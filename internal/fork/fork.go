// Package fork chooses between competing chains — the moment consensus gets its vote.
package fork

import "errors"

var (
	// ErrNoCandidates: ForkChoice was offered zero chains (or every offered chain was nil). GIVEN.
	ErrNoCandidates = errors.New("fork: no candidate chains offered")
	// ErrNoValidChain: every offered chain failed validation or proof-of-work verification. GIVEN.
	ErrNoValidChain = errors.New("fork: none of the offered chains is valid")
)
