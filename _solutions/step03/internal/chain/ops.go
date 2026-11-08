package chain

import (
	"fmt"
	"strings"
)

// REFERENCE SOLUTION for Step 03. Try it yourself first!

func NewBlockchain(genesisData string, genesisTimestamp int64, difficulty int) *Blockchain {
	g := Block{Index: 0, Timestamp: genesisTimestamp, Data: genesisData, PrevHash: GenesisPrevHash}
	g = mine(g, difficulty)
	return &Blockchain{Difficulty: difficulty, Blocks: []Block{g}}
}

func (bc *Blockchain) AddBlock(data string, timestamp int64) Block {
	prev := Block{Index: -1, PrevHash: GenesisPrevHash}
	if len(bc.Blocks) > 0 {
		prev = bc.Blocks[len(bc.Blocks)-1]
	}
	b := Block{Index: prev.Index + 1, Timestamp: timestamp, Data: data, PrevHash: prev.Hash}
	b = mine(b, bc.Difficulty)
	bc.Blocks = append(bc.Blocks, b)
	return b
}

func (bc *Blockchain) Validate() error {
	for i, b := range bc.Blocks {
		if b.Hash != CalculateHash(b) {
			return fmt.Errorf("block %d: %w", i, ErrInvalidHash)
		}
		if b.Index != i {
			return fmt.Errorf("block %d: %w", i, ErrWrongIndex)
		}
		if i == 0 {
			if b.PrevHash != GenesisPrevHash {
				return fmt.Errorf("genesis: %w", ErrBrokenLink)
			}
			continue
		}
		if b.PrevHash != bc.Blocks[i-1].Hash {
			return fmt.Errorf("block %d: %w", i, ErrBrokenLink)
		}
	}
	return nil
}

// mine keeps bumping the nonce until the block's hash meets the difficulty.
func mine(b Block, difficulty int) Block {
	for {
		b.Hash = CalculateHash(b)
		if meetsZeros(b.Hash, difficulty) {
			return b
		}
		b.Nonce++
	}
}

func meetsZeros(hash string, difficulty int) bool {
	if difficulty <= 0 {
		return true
	}
	if len(hash) < difficulty {
		return false
	}
	return strings.HasPrefix(hash, strings.Repeat("0", difficulty))
}
