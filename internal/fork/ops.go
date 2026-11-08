package fork

import (
	"blockchain/internal/chain"

	"blockchain/internal/todo"
)

// STEP 07: implement fork choice. Read docs/steps/step-07-fork-choice.md first.   Run:  go test ./internal/fork -run Step07

// ForkChoice resolves competing chains that share the same genesis.
//
//   - A chain is ELIGIBLE only if: it has at least one block, chain.Validate() passes, and every block meets its
//     declared difficulty (pow.Verify). Ineligible chains are discarded.
//   - The chain with the MOST cumulative proof-of-work wins (each block contributes 16^difficulty of work) —
//     the "greatest work wins" rule of Nakamoto consensus.
//   - If two eligible chains have equal work, the one whose tip hash is lower (byte order of the hex string) wins.
//     A deterministic stand-in for "first seen wins" — every honest node picks the same chain.
//
// Returns ErrNoCandidates for zero offers; ErrNoValidChain if nothing offered is eligible.
func ForkChoice(chains []*chain.Blockchain) (*chain.Blockchain, error) {
	return nil, todo.NotImplemented("Step 07: ForkChoice")
}
