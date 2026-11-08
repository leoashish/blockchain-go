package chain_test

import (
	"errors"
	"strings"
	"testing"

	"blockchain/internal/chain"
	"blockchain/internal/pow"
)

// Step 03 -- the chain: validation.   Run:  go test ./internal/chain -run Step03 -v

func newChain() *chain.Blockchain {
	return chain.NewBlockchain("genesis", 1_000_000_000, 2)
}

func requireBlocks(t *testing.T, bc *chain.Blockchain, want int) {
	t.Helper()
	if len(bc.Blocks) != want {
		t.Fatalf("the chain has %d blocks, want %d (the test cannot index into an empty chain)", len(bc.Blocks), want)
	}
}

func TestStep03_NewBlockchain_CreatesASingleMinedGenesisBlock(t *testing.T) {
	bc := newChain()

	requireBlocks(t, bc, 1)
	if bc.Difficulty != 2 {
		t.Fatalf("Difficulty = %d, want 2", bc.Difficulty)
	}
	g := bc.Blocks[0]
	if g.Index != 0 {
		t.Fatalf("genesis Index = %d, want 0", g.Index)
	}
	if g.PrevHash != chain.GenesisPrevHash {
		t.Fatalf("genesis PrevHash = %q, want %q (sixty-four zeros)", g.PrevHash, chain.GenesisPrevHash)
	}
	if g.Data != "genesis" || g.Timestamp != 1_000_000_000 {
		t.Fatalf("genesis = %+v, want the genesisData and genesisTimestamp you passed", g)
	}
}

func TestStep03_GenesisBlock_MeetsTheDifficulty(t *testing.T) {
	bc := newChain()

	requireBlocks(t, bc, 1)
	if !pow.Verify(bc.Blocks[0], bc.Difficulty) {
		t.Fatal("the genesis block does not meet the chain's difficulty: it must be mined, not just built")
	}
}

func TestStep03_AddBlock_AppendsMinedBlocksLinkedToTheChain(t *testing.T) {
	bc := newChain()

	requireBlocks(t, bc, 1)
	g := bc.Blocks[0]

	b1 := bc.AddBlock("transactions of block 1", 1_000_000_060)
	b2 := bc.AddBlock("transactions of block 2", 1_000_000_120)

	requireBlocks(t, bc, 3)
	if b1.Index != 1 || b2.Index != 2 {
		t.Fatalf("indices are b1=%d, b2=%d, want 1 and 2", b1.Index, b2.Index)
	}
	if b1.PrevHash != g.Hash {
		t.Fatalf("block 1 PrevHash = %q, want the genesis hash %q", b1.PrevHash, g.Hash)
	}
	if b2.PrevHash != b1.Hash {
		t.Fatalf("block 2 PrevHash = %q, want block 1's hash %q", b2.PrevHash, b1.Hash)
	}
}

func TestStep03_AddBlock_EveryBlockMeetsTheDifficulty(t *testing.T) {
	bc := newChain()
	bc.AddBlock("one", 1_000_000_060)
	bc.AddBlock("two", 1_000_000_120)

	requireBlocks(t, bc, 3)
	for i, b := range bc.Blocks {
		if !pow.Verify(b, bc.Difficulty) {
			t.Fatalf("block %d does not meet difficulty %d", i, bc.Difficulty)
		}
	}
}

func TestStep03_Validate_AcceptsAFreshlyBuiltChain(t *testing.T) {
	bc := newChain()
	bc.AddBlock("one", 1_000_000_060)

	if err := bc.Validate(); err != nil {
		t.Fatalf("Validate() = %v, want nil", err)
	}
}

func TestStep03_Validate_DetectsATamperedBlock(t *testing.T) {
	bc := newChain()
	bc.AddBlock("one", 1_000_000_060)
	requireBlocks(t, bc, 2)

	// Flip a byte in block 1's data WITHOUT recomputing its hash.
	bc.Blocks[1].Data = "oneeeee"

	if err := bc.Validate(); err == nil {
		t.Fatal("Validate() = nil, want an error")
	} else if !errors.Is(err, chain.ErrInvalidHash) {
		t.Fatalf("Validate() = %v, want an error wrapping chain.ErrInvalidHash", err)
	}
}

func TestStep03_Validate_DetectsAWrongIndex_AfterRehashing(t *testing.T) {
	bc := newChain()
	bc.AddBlock("one", 1_000_000_060)
	requireBlocks(t, bc, 2)

	tainted := bc.Blocks[1]
	tainted.Index = 7 // reassign the index, and re-hash so the Hash still checks out
	tainted.Hash = chain.CalculateHash(tainted)
	bc.Blocks[1] = tainted

	if err := bc.Validate(); err == nil {
		t.Fatal("Validate() = nil, want an error")
	} else if !errors.Is(err, chain.ErrWrongIndex) {
		t.Fatalf("Validate() = %v, want an error wrapping chain.ErrWrongIndex", err)
	}
}

func TestStep03_Validate_DetectsABrokenLink_AfterRehashing(t *testing.T) {
	bc := newChain()
	bc.AddBlock("one", 1_000_000_060)
	requireBlocks(t, bc, 2)

	tampered := bc.Blocks[1]
	tampered.PrevHash = strings.Repeat("deadbeef", 8) // point block 1 at a fake parent, re-hash
	tampered.Hash = chain.CalculateHash(tampered)
	bc.Blocks[1] = tampered

	if err := bc.Validate(); err == nil {
		t.Fatal("Validate() = nil, want an error")
	} else if !errors.Is(err, chain.ErrBrokenLink) {
		t.Fatalf("Validate() = %v, want an error wrapping chain.ErrBrokenLink", err)
	}
}

func TestStep03_Validate_DetectsABrokenGenesisLink(t *testing.T) {
	bc := newChain()
	requireBlocks(t, bc, 1)

	tangled := bc.Blocks[0]
	tangled.PrevHash = strings.Repeat("abc", 22)[:64] // not the genesis link
	tangled.Hash = chain.CalculateHash(tangled)
	bc.Blocks[0] = tangled

	if err := bc.Validate(); err == nil {
		t.Fatal("Validate() = nil, want an error")
	} else if !errors.Is(err, chain.ErrBrokenLink) {
		t.Fatalf("Validate() = %v, want an error wrapping chain.ErrBrokenLink", err)
	}
}

func TestStep03_Validate_AcceptsAnEmptyChain(t *testing.T) {
	bc := &chain.Blockchain{Difficulty: 2}

	if err := bc.Validate(); err != nil {
		t.Fatalf("Validate() = %v, want nil (nothing to invalidate)", err)
	}
}
