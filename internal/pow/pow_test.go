package pow_test

import (
	"strings"
	"testing"

	"blockchain/internal/chain"
	"blockchain/internal/pow"
)

// Step 02 -- proof of work.   Run:  go test ./internal/pow -run Step02 -v

func TestStep02_MeetsDifficulty_AcceptsOnlyLeadingZeros(t *testing.T) {
	cases := map[string]struct {
		hash       string
		difficulty int
		want       bool
	}{
		"all leading zeros, exactly enough": {hash: "0000abcd", difficulty: 4, want: true},
		"zeros then a digit":                {hash: "0001abcd", difficulty: 3, want: true},
		"more zeros than required":          {hash: "000000abcd", difficulty: 2, want: true},
		"zeros in the middle don't count":   {hash: "abc00000", difficulty: 3, want: false},
		"not enough zeros":                  {hash: "0000abcd", difficulty: 5, want: false},
		"zero difficulty passes":            {hash: "abcd", difficulty: 0, want: true},
		"empty hash":                        {hash: "", difficulty: 1, want: false},
	}
	for name, tc := range cases {
		if got := pow.MeetsDifficulty(tc.hash, tc.difficulty); got != tc.want {
			t.Errorf("%s: MeetsDifficulty(%q, %d) = %v, want %v", name, tc.hash, tc.difficulty, got, tc.want)
		}
	}
}

func TestStep02_Mine_ProducesAHashWithTheRequiredLeadingZeros(t *testing.T) {
	b := chain.NewBlock(7, "some data", strings.Repeat("ab", 32), 1700000000)

	mined := pow.Mine(b, 2)

	if !strings.HasPrefix(mined.Hash, "00") {
		t.Fatalf("after mining with difficulty 2, hash = %q, want at least two leading zeros", mined.Hash)
	}
}

func TestStep02_Mine_PreservesEverythingExceptTheNonceAndHash(t *testing.T) {
	b := chain.NewBlock(7, "the data stays the same", strings.Repeat("ab", 32), 1700000000)

	mined := pow.Mine(b, 1)

	if mined.Index != b.Index || mined.Timestamp != b.Timestamp || mined.Data != b.Data || mined.PrevHash != b.PrevHash {
		t.Fatalf("mining changed the identity fields: %+v -> %+v", b, mined)
	}
	if mined.Nonce < 0 {
		t.Fatalf("mined Nonce = %d, want a non-negative nonce", mined.Nonce)
	}
	if mined.Hash == "" {
		t.Fatal("the mined block has an empty Hash: mining must set it")
	}
}

func TestStep02_Mine_IsDeterministic_ForTheSameInputs(t *testing.T) {
	b := chain.NewBlock(3, "deterministic mining", strings.Repeat("cd", 32), 1700000000)

	first := pow.Mine(b, 2)
	second := pow.Mine(b, 2)

	if !strings.HasPrefix(first.Hash, "00") {
		t.Fatal("mining must set a hash that meets the difficulty")
	}
	if first != second {
		t.Fatalf("mining the same block twice gave different results: %+v vs %+v", first, second)
	}
}

func TestStep02_Verify_AcceptsAProperlyMinedBlock(t *testing.T) {
	b := chain.NewBlock(7, "verify me", strings.Repeat("ef", 32), 1700000000)
	mined := pow.Mine(b, 2)

	if !pow.Verify(mined, 2) {
		t.Fatalf("Verify(mined, 2) = false, want true")
	}
}

func TestStep02_Verify_RejectsAnIntactBlockWithoutEnoughLeadingZeros(t *testing.T) {
	b := chain.Block{Index: 7, Timestamp: 1700000000, Data: "proof is cheaper to check than to make", PrevHash: strings.Repeat("ef", 32), Nonce: 0}

	// Find a nonce whose hash does NOT start with '0' at all (exists easily: ~16 tries on average).
	for {
		b.Hash = chain.CalculateHash(b)
		if !strings.HasPrefix(b.Hash, "0") {
			break
		}
		b.Nonce++
	}

	if !chain.VerifyHash(b) {
		t.Fatal("the block must be intact (its Hash matches its contents) for this test to mean anything")
	}
	if pow.Verify(b, 1) {
		t.Fatal("Verify(b, 1) = true, want false: an intact block whose hash has no leading zero does not meet difficulty 1")
	}
}

func TestStep02_Verify_RejectsATamperedBlock(t *testing.T) {
	b := chain.NewBlock(7, "tamper with me", strings.Repeat("ef", 32), 1700000000)
	mined := pow.Mine(b, 2)

	if !pow.Verify(mined, 2) {
		t.Fatal("baseline: the freshly mined block must verify (mining is not done yet)")
	}

	tampered := mined
	tampered.Data = tampered.Data + "."

	if pow.Verify(tampered, 2) {
		t.Fatal("Verify(tampered) = true, want false: changing one character breaks the hash")
	}
}
