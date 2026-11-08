package chain

import "errors"

var (
	// ErrInvalidHash: the block's Hash does not match its contents (someone tampered with it, or it was never computed). GIVEN.
	ErrInvalidHash = errors.New("chain: block hash does not match its contents")
	// ErrWrongIndex: the block's Index is not the next contiguous index from the previous block. GIVEN.
	ErrWrongIndex = errors.New("chain: block index is not the next contiguous index")
	// ErrBrokenLink: the block does not reference the previous block's hash. GIVEN.
	ErrBrokenLink = errors.New("chain: block does not link to the previous block")
)

// Blockchain is the whole chain, with the genesis block at index 0.
//
// STEP 03: the chain proves its own integrity. Every block, when validated, must satisfy:
//   - its Hash matches its contents (ErrInvalidHash),
//   - its Index is contiguous from the previous block (ErrWrongIndex),
//   - it references the previous block's Hash (ErrBrokenLink), with the genesis block pointing at GenesisPrevHash (sixty-four zeros).
//
// Blocks is exported so you (and the tests) can corrupt entries and watch Validate catch the damage.
// A production store would seal it behind a read API; here the fun is seeing why tampering cannot hide.
type Blockchain struct {
	Difficulty int
	Blocks     []Block
}
