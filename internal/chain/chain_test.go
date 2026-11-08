package chain_test

import (
	"strings"
	"testing"

	"blockchain/internal/chain"
)

// Step 01 -- the block: SHA-256 hashing and linking.   Run:  go test ./internal/chain -run Step01 -v

func TestStep01_CalculateHash_IsDeterministic_AndValidHex(t *testing.T) {
	b := chain.Block{Index: 0, Timestamp: 1700000000, Data: "genesis", PrevHash: chain.GenesisPrevHash, Nonce: 0}

	first := chain.CalculateHash(b)
	second := chain.CalculateHash(b)

	if first != second {
		t.Fatalf("CalculateHash is not deterministic: %q then %q", first, second)
	}
	if len(first) != 64 {
		t.Fatalf("hash is %d hex characters, want 64 (SHA-256). An empty hash means CalculateHash was not implemented.", len(first))
	}
	for _, r := range first {
		if !strings.ContainsRune("0123456789abcdef", r) {
			t.Fatalf("hash %q contains %q, want only lowercase hex characters", first, r)
		}
	}
}

func TestStep01_CalculateHash_Changes_WhenAnyIdentityFieldChanges(t *testing.T) {
	base := chain.Block{Index: 7, Timestamp: 123456789, Data: "base data", PrevHash: strings.Repeat("5e", 32), Nonce: 42}
	baseHash := chain.CalculateHash(base)

	variants := map[string]chain.Block{
		"index":     {Index: base.Index + 1, Timestamp: base.Timestamp, Data: base.Data, PrevHash: base.PrevHash, Nonce: base.Nonce},
		"timestamp": {Index: base.Index, Timestamp: base.Timestamp + 1, Data: base.Data, PrevHash: base.PrevHash, Nonce: base.Nonce},
		"data":      {Index: base.Index, Timestamp: base.Timestamp, Data: base.Data + "!", PrevHash: base.PrevHash, Nonce: base.Nonce},
		"prevHash":  {Index: base.Index, Timestamp: base.Timestamp, Data: base.Data, PrevHash: strings.Repeat("a5", 32), Nonce: base.Nonce},
		"nonce":     {Index: base.Index, Timestamp: base.Timestamp, Data: base.Data, PrevHash: base.PrevHash, Nonce: base.Nonce + 1},
	}

	for field, variant := range variants {
		if got := chain.CalculateHash(variant); got == baseHash {
			t.Errorf("changing %s did not change the hash", field)
		}
	}

	// The Hash field itself must not influence the hash (it is the output).
	withHash := base
	withHash.Hash = baseHash
	if got := chain.CalculateHash(withHash); got != baseHash {
		t.Fatalf("CalculateHash must ignore the Hash field: got %q", got)
	}
}

func TestStep01_NewBlock_ComputesTheHashFromTheIdentityFields(t *testing.T) {
	b := chain.NewBlock(3, "some transactions", strings.Repeat("f0", 32), 1700000123)

	if b.Index != 3 {
		t.Fatalf("NewBlock: Index = %d, want 3", b.Index)
	}
	if b.Hash == "" {
		t.Fatal("NewBlock: Hash is empty -- the hash was not computed")
	}
	if b.Hash != chain.CalculateHash(b) {
		t.Fatalf("NewBlock: Hash = %q, but CalculateHash(block) = %q", b.Hash, chain.CalculateHash(b))
	}
}

func TestStep01_NewGenesisBlock_IsTheBeginningOfTheChain(t *testing.T) {
	g := chain.NewGenesisBlock("the first block ever", 1_000_000_000)

	if g.Index != 0 {
		t.Fatalf("genesis Index = %d, want 0", g.Index)
	}
	if g.PrevHash != chain.GenesisPrevHash {
		t.Fatalf("genesis PrevHash = %q, want %q (sixty-four zeros)", g.PrevHash, chain.GenesisPrevHash)
	}
	if len(chain.GenesisPrevHash) != 64 || strings.Count(chain.GenesisPrevHash, "0") != 64 {
		t.Fatalf("GenesisPrevHash = %q, want exactly sixty-four zeros", chain.GenesisPrevHash)
	}
	if g.Hash == "" || g.Hash != chain.CalculateHash(g) {
		t.Fatalf("genesis Hash = %q, want CalculateHash(genesis) = %q", g.Hash, chain.CalculateHash(g))
	}
}

func TestStep01_NewBlock_LinksToThePreviousBlock_InAChainOfThree(t *testing.T) {
	g := chain.NewGenesisBlock("genesis", 1_000_000_000)
	b1 := chain.NewBlock(1, "first block's data", g.Hash, 1_000_000_060)
	b2 := chain.NewBlock(2, "second block's data", b1.Hash, 1_000_000_120)

	if b1.PrevHash != g.Hash {
		t.Fatalf("block 1 PrevHash = %q, want the genesis hash %q", b1.PrevHash, g.Hash)
	}
	if b2.PrevHash != b1.Hash {
		t.Fatalf("block 2 PrevHash = %q, want block 1's hash %q", b2.PrevHash, b1.Hash)
	}
	if b1.Hash == b2.Hash {
		t.Fatal("two different blocks got the same hash")
	}
}

func TestStep01_VerifyHash_AcceptsIntactAndRejectsTampered(t *testing.T) {
	b := chain.NewBlock(9, "blocks in a blockchain do not lie", strings.Repeat("1b", 32), 1_700_000_000)

	if !chain.VerifyHash(b) {
		t.Fatal("VerifyHash(untampered block) = false, want true")
	}

	tampered := b
	tampered.Data = b.Data + "." // one character appended
	if chain.VerifyHash(tampered) {
		t.Fatal("VerifyHash(tampered block) = true, want false: the hash covers Data")
	}

	posed := b
	posed.Hash = strings.Repeat("deadbeef", 8)
	if chain.VerifyHash(posed) {
		t.Fatal("VerifyHash(block with a mismatched Hash) = true, want false")
	}
}
