package fork

import (
	"blockchain/internal/chain"
	"blockchain/internal/pow"
)

// REFERENCE SOLUTION for Step 07. Try it yourself first!

func ForkChoice(chains []*chain.Blockchain) (*chain.Blockchain, error) {
	if len(chains) == 0 {
		return nil, ErrNoCandidates
	}
	var best *chain.Blockchain
	var bestWork uint64
	for _, bc := range chains {
		if !validChain(bc) {
			continue
		}
		work := workOf(bc)
		if best == nil || work > bestWork || (work == bestWork && tipLess(bc, best)) {
			best = bc
			bestWork = work
		}
	}
	if best == nil {
		return nil, ErrNoValidChain
	}
	return best, nil
}

func validChain(bc *chain.Blockchain) bool {
	if bc == nil {
		return false
	}
	if len(bc.Blocks) == 0 {
		return false
	}
	if err := bc.Validate(); err != nil {
		return false
	}
	for _, b := range bc.Blocks {
		if !pow.Verify(b, bc.Difficulty) {
			return false
		}
	}
	return true
}

func workOf(bc *chain.Blockchain) uint64 {
	var per uint64 = 1
	if bits := uint(bc.Difficulty) * 4; bits >= 64 {
		per = ^uint64(0)
	} else {
		per = uint64(1) << bits
	}
	var sum uint64
	for range bc.Blocks {
		if ^uint64(0)-sum < per {
			return ^uint64(0)
		}
		sum += per
	}
	return sum
}

func tipLess(a, b *chain.Blockchain) bool {
	return a.Blocks[len(a.Blocks)-1].Hash < b.Blocks[len(b.Blocks)-1].Hash
}
