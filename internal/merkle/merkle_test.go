package merkle_test

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"testing"

	"blockchain/internal/merkle"
)

// Step 08 -- Merkle trees.   Run:  go test ./internal/merkle -run Step08 -v

func leaves(data ...string) [][]byte {
	out := make([][]byte, len(data))
	for i, s := range data {
		out[i] = []byte(s)
	}
	return out
}

func sha256Concat(a, b []byte) []byte {
	h := sha256.Sum256(append(append([]byte{}, a...), b...))
	return h[:]
}

func TestStep08_Root_OfEmptyLeaves_IsSixtyFourZeroHexChars(t *testing.T) {
	if got := merkle.Root(nil); got != strings.Repeat("00", 32) {
		t.Fatalf("Root(nil) = %q, want 64 zeros (the zero hash)", got)
	}
}

func TestStep08_Root_SingleLeaf_IsItsOwnHash(t *testing.T) {
	leaf := []byte("only transaction")
	want := hex.EncodeToString(mustSha(leaf))

	if got := merkle.Root([][]byte{leaf}); got != want {
		t.Fatalf("Root = %q, want %q (sha256 of the leaf)", got, want)
	}
}

func mustSha(b []byte) []byte {
	h := sha256.Sum256(b)
	return h[:]
}

func TestStep08_Root_OfTwoLeaves_EqualsShaOfConcat(t *testing.T) {
	a, b := []byte("one"), []byte("two")
	want := hex.EncodeToString(sha256Concat(mustSha(a), mustSha(b)))

	if got := merkle.Root([][]byte{a, b}); got != want {
		t.Fatalf("Root = %q, want %q (sha256(sha256(a) ‖ sha256(b)))", got, want)
	}
}

func TestStep08_Root_Changes_WhenAnyLeafByteChanges(t *testing.T) {
	base := [][]byte{[]byte("alpha"), []byte("beta"), []byte("gamma"), []byte("delta")}
	baseRoot := merkle.Root(base)

	for i := range base {
		tampered := leaves("alpha", "beta", "gamma", "delta")
		tampered[i] = []byte(append(append([]byte{}, base[i]...), 'X'))
		if got := merkle.Root(tampered); got == baseRoot {
			t.Errorf("changing leaf %d did not change the root", i)
		}
	}
}

func TestStep08_Root_DuplicatesTheOddNode_AtEveryLevel(t *testing.T) {
	l := leaves("a", "b", "c") // one odd leaf at level 0

	// manual: sha(a), sha(b), sha(sha(a)‖sha(b)); sha(c‖c); root = sha(prev ‖ sha(c‖c))
	level1 := [][]byte{mustSha(l[0]), mustSha(l[1]), mustSha(l[2])}
	left := sha256Concat(level1[0], level1[1])
	right := sha256Concat(level1[2], level1[2]) // odd duplicated
	want := hex.EncodeToString(sha256Concat(left, right))

	if got := merkle.Root(l); got != want {
		t.Fatalf("Root = %q, want %q (odd node duplicated at every level)", got, want)
	}
}

func TestStep08_Prove_And_Verify_RoundTrip_AtEverySizeFromOneToSeventeen(t *testing.T) {
	for n := 1; n <= 17; n++ {
		ls := make([][]byte, n)
		for i := range ls {
			ls[i] = []byte{byte('a' + i)}
		}
		root := merkle.Root(ls)
		for i := 0; i < n; i++ {
			proof, err := merkle.Prove(ls, i)
			if err != nil {
				t.Fatalf("n=%d i=%d Prove: %v", n, i, err)
			}
			if !merkle.Verify(ls[i], i, proof, root) {
				t.Fatalf("n=%d i=%d: Verify(leaf, i, proof, root) = false, want true", n, i)
			}
		}
	}
}

func TestStep08_Verify_RejectsATamperedProof_OrWrongIndexOrWrongRoot(t *testing.T) {
	ls := leaves("aa", "bb", "cc", "dd", "ee")
	root := merkle.Root(ls)

	proof, err := merkle.Prove(ls, 2)
	if err != nil {
		t.Fatalf("Prove: %v", err)
	}

	// flip the last character of the first proof hash
	tamper := []string{proof[0]}
	b := tamper[0]
	if b[len(b)-1] == 'a' {
		tamper[0] = b[:len(b)-1] + "b"
	} else {
		tamper[0] = b[:len(b)-1] + "a"
	}
	if merkle.Verify(ls[2], 2, tamper, root) {
		t.Fatal("a tampered proof element must not verify")
	}
	if merkle.Verify(ls[2], 3, proof, root) {
		t.Fatal("the wrong index must not verify")
	}
	badRoot := root
	if badRoot[len(badRoot)-1] == 'a' {
		badRoot = badRoot[:len(badRoot)-1] + "b"
	} else {
		badRoot = badRoot[:len(badRoot)-1] + "a"
	}
	if merkle.Verify(ls[2], 2, proof, badRoot) {
		t.Fatal("the wrong root must not verify")
	}
}

func TestStep08_Prove_RejectsOutOfRangeIndex(t *testing.T) {
	ls := leaves("x", "y", "z")

	if _, err := merkle.Prove(ls, 5); !errors.Is(err, merkle.ErrIndexOutOfRange) {
		t.Fatalf("Prove(3 leaves, index 5) = %v, want merkle.ErrIndexOutOfRange", err)
	}
	if _, err := merkle.Prove(ls, -1); !errors.Is(err, merkle.ErrIndexOutOfRange) {
		t.Fatalf("Prove(3 leaves, index -1) = %v, want merkle.ErrIndexOutOfRange", err)
	}
}

func TestStep08_Prove_And_Verify_HoldUpOnAThousandLeafTree(t *testing.T) {
	const n = 1000
	ls := make([][]byte, n)
	for i := range ls {
		ls[i] = []byte(hex.EncodeToString([]byte{byte(i), byte(i >> 8)}))
	}
	root := merkle.Root(ls)

	for _, i := range []int{0, 1, n / 2, n - 1} {
		proof, err := merkle.Prove(ls, i)
		if err != nil {
			t.Fatalf("Prove(%d): %v", i, err)
		}
		if !merkle.Verify(ls[i], i, proof, root) {
			t.Fatalf("Verify(index %d) = false, want true", i)
		}
	}
}

func TestStep08_Verify_RejectsAProofWithMissingOrReorderedElements(t *testing.T) {
	ls := leaves("aa", "bb", "cc", "dd")
	root := merkle.Root(ls)

	proof, err := merkle.Prove(ls, 0)
	if err != nil {
		t.Fatalf("Prove: %v", err)
	}

	short := proof[:len(proof)-1]
	if merkle.Verify(ls[0], 0, short, root) {
		t.Fatal("a proof with the last element missing must not verify")
	}

	reordered := append([]string{}, proof...)
	reordered[0], reordered[1] = reordered[1], reordered[0]
	if merkle.Verify(ls[0], 0, reordered, root) {
		t.Fatal("a reordered proof must not verify")
	}
}
