// Package persist saves the chain to disk and proves the chain on the way back in.
package persist

import (
	"blockchain/internal/chain"

	"blockchain/internal/todo"
)

// STEP 06: implement persistence. Read docs/steps/step-06-persistence.md first.   Run:  go test ./internal/persist -run Step06

// SaveChain writes the chain to path ATOMICALLY: write to path+".tmp", fsync, then os.Rename (the rename is atomic on
// POSIX). After SaveChain returns, the file either contains the whole chain or does not exist at all.
func SaveChain(path string, bc *chain.Blockchain) error {
	return todo.NotImplemented("Step 06: SaveChain")
}

// LoadChain reads a chain file and, BEFORE returning it, fully verifies the chain:
//   - every block's Hash matches its contents,
//   - indices are contiguous from 0,
//   - the genesis block links to GenesisPrevHash, every other block links to its predecessor,
//   - every block meets the chain's stored difficulty.
//
// A tampered or corrupted file must never load. Errors are wrapped so callers can inspect them.
func LoadChain(path string) (*chain.Blockchain, error) {
	return nil, todo.NotImplemented("Step 06: LoadChain")
}
