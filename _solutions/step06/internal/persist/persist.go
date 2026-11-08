package persist

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"blockchain/internal/chain"
	"blockchain/internal/pow"
)

// REFERENCE SOLUTION for Step 06. Try it yourself first!

type fileChain struct {
	Difficulty int           `json:"difficulty"`
	Blocks     []chain.Block `json:"blocks"`
}

func SaveChain(path string, bc *chain.Blockchain) error {
	data, err := json.Marshal(fileChain{Difficulty: bc.Difficulty, Blocks: bc.Blocks})
	if err != nil {
		return fmt.Errorf("persist: encode %s: %w", path, err)
	}
	tmp := path + ".tmp"
	f, err := os.OpenFile(tmp, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return fmt.Errorf("persist: open temp file: %w", err)
	}
	if _, err := f.Write(data); err != nil {
		_ = f.Close()
		return fmt.Errorf("persist: write temp file: %w", err)
	}
	if err := f.Sync(); err != nil {
		_ = f.Close()
		return fmt.Errorf("persist: sync temp file: %w", err)
	}
	if err := f.Close(); err != nil {
		return fmt.Errorf("persist: close temp file: %w", err)
	}
	if err := os.Rename(tmp, path); err != nil {
		return fmt.Errorf("persist: rename temp file: %w", err)
	}
	return nil
}

func LoadChain(path string) (*chain.Blockchain, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("persist: read %s: %w", path, err)
	}
	var fc fileChain
	if err := json.Unmarshal(data, &fc); err != nil {
		return nil, fmt.Errorf("persist: parse %s: %w", path, err)
	}
	bc := &chain.Blockchain{Difficulty: fc.Difficulty, Blocks: fc.Blocks}

	if err := bc.Validate(); err != nil {
		return nil, fmt.Errorf("persist: chain validation failed: %w", err)
	}
	for i, b := range bc.Blocks {
		if !pow.Verify(b, bc.Difficulty) {
			return nil, fmt.Errorf("persist: block %d does not meet difficulty %d", i, bc.Difficulty)
		}
	}
	if len(bc.Blocks) == 0 {
		return nil, errors.New("persist: the saved chain contains no genesis block")
	}
	return bc, nil
}
