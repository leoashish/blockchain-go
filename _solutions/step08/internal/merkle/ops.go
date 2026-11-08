package merkle

import (
	"crypto/sha256"
	"encoding/hex"
)

// REFERENCE SOLUTION for Step 08. Try it yourself first!

func Root(leaves [][]byte) string {
	if len(leaves) == 0 {
		return hex.EncodeToString(make([]byte, 32))
	}
	cur := hashLeaves(leaves)
	for len(cur) > 1 {
		cur = pairUp(cur)
	}
	return hex.EncodeToString(cur[0])
}

func Prove(leaves [][]byte, index int) ([]string, error) {
	if index < 0 || index >= len(leaves) {
		return nil, ErrIndexOutOfRange
	}
	proof := []string{}
	cur := hashLeaves(leaves)
	idx := index
	for len(cur) > 1 {
		if len(cur)%2 == 1 {
			cur = append(append([][]byte{}, cur...), cur[len(cur)-1]) // duplicate odd node (Bitcoin's rule)
		}
		proof = append(proof, hex.EncodeToString(cur[idx^1]))
		cur = pairUp(cur)
		idx /= 2
	}
	return proof, nil
}

func Verify(leaf []byte, index int, proof []string, root string) bool {
	if index < 0 {
		return false
	}
	h := hashLeaf(leaf)
	idx := index
	for _, s := range proof {
		sib, err := hex.DecodeString(s)
		if err != nil || len(sib) != 32 {
			return false
		}
		if idx%2 == 0 {
			h = hashPair(h, sib)
		} else {
			h = hashPair(sib, h)
		}
		idx /= 2
	}
	return hex.EncodeToString(h) == root
}

func hashLeaf(b []byte) []byte {
	h := sha256.Sum256(b)
	return h[:]
}

func hashPair(a, b []byte) []byte {
	sum := sha256.Sum256(append(append([]byte{}, a...), b...))
	return sum[:]
}

func hashLeaves(leaves [][]byte) [][]byte {
	out := make([][]byte, len(leaves))
	for i, l := range leaves {
		out[i] = hashLeaf(l)
	}
	return out
}

func pairUp(cur [][]byte) [][]byte {
	if len(cur)%2 == 1 {
		cur = append(append([][]byte{}, cur...), cur[len(cur)-1]) // duplicate odd node (Bitcoin's rule)
	}
	out := make([][]byte, 0, len(cur)/2)
	for i := 0; i < len(cur); i += 2 {
		out = append(out, hashPair(cur[i], cur[i+1]))
	}
	return out
}
