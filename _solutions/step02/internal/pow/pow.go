package pow

import (
	"strings"

	"blockchain/internal/chain"
)

// REFERENCE SOLUTION for Step 02. Try it yourself first!

func MeetsDifficulty(hash string, difficulty int) bool {
	if difficulty <= 0 {
		return true
	}
	if len(hash) < difficulty {
		return false
	}
	return strings.HasPrefix(hash, strings.Repeat("0", difficulty))
}

func Mine(b chain.Block, difficulty int) chain.Block {
	for {
		b.Hash = chain.CalculateHash(b)
		if MeetsDifficulty(b.Hash, difficulty) {
			return b
		}
		b.Nonce++
	}
}

func Verify(b chain.Block, difficulty int) bool {
	if !chain.VerifyHash(b) {
		return false
	}
	return MeetsDifficulty(b.Hash, difficulty)
}
