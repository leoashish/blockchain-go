package merkle

import "blockchain/internal/todo"

// STEP 08: implement Merkle trees. Read docs/steps/step-08-merkle-tree.md first.   Run:  go test ./internal/merkle -run Step08

// Root computes the Merkle root of leaves:
//   - every leaf is hashed with SHA-256 to become a node,
//   - nodes pair up, each pair is hashed with SHA-256 over the CONCATENATION of the two children,
//   - an odd node at any level is duplicated to make a pair (Bitcoin's rule),
//   - zero leaves -> the all-zero 32 bytes, "0000000000000000000000000000000000000000000000000000000000000000".
//
// Returns the root as lowercase hex.
func Root(leaves [][]byte) string {
	return "" // TODO Step 08
}

// Prove returns the sibling hashes (hex) on the path from `index` to the root, ordered bottom-up.
// An index outside the leaves -> ErrIndexOutOfRange.
func Prove(leaves [][]byte, index int) ([]string, error) {
	return nil, todo.NotImplemented("Step 08: Prove")
}

// Verify reports whether `leaf` is the `index`-th leaf of a tree whose root is `root`, given `proof` (the sibling hashes
// bottom-up, hex). A proof that does not reproduce `root` — or a malformed hex element — returns false. It must never panic.
func Verify(leaf []byte, index int, proof []string, root string) bool {
	return false // TODO Step 08
}
