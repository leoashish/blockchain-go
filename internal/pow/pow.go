// Package pow implements proof of work: the honest way to make a block expensive to produce and cheap to check.
package pow

import (
	"blockchain/internal/chain"
)

// MeetsDifficulty reports whether hash begins with at least `difficulty` zero hex characters.
// A difficulty of zero (or less) always passes.
func MeetsDifficulty(hash string, difficulty int) bool {
	return false // TODO Step 02
}

// Mine increments the block's Nonce until its hash meets the required difficulty, sets Hash, and returns the block.
// Iteration starts from the block's current Nonce, so mining the same block twice is deterministic.
func Mine(b chain.Block, difficulty int) chain.Block {
	return chain.Block{} // TODO Step 02
}

// Verify reports whether b is intact (chain.VerifyHash) AND its hash meets the given difficulty.
func Verify(b chain.Block, difficulty int) bool {
	return false // TODO Step 02
}
