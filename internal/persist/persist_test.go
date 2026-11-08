package persist_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"blockchain/internal/chain"
	"blockchain/internal/persist"
)

// Step 06 -- persistence.   Run:  go test ./internal/persist -run Step06 -v

func buildChain() *chain.Blockchain {
	bc := chain.NewBlockchain("genesis", 1_000_000_000, 2)
	bc.AddBlock("block one", 1_000_000_060)
	bc.AddBlock("block two", 1_000_000_120)
	return bc
}

func TestStep06_SaveThenLoad_RoundTripsTheChain(t *testing.T) {
	bc := buildChain()
	path := filepath.Join(t.TempDir(), "chain.json")

	if err := persist.SaveChain(path, bc); err != nil {
		t.Fatalf("SaveChain: %v", err)
	}

	loaded, err := persist.LoadChain(path)
	if err != nil {
		t.Fatalf("LoadChain: %v", err)
	}

	if loaded.Difficulty != bc.Difficulty {
		t.Fatalf("Difficulty = %d, want %d", loaded.Difficulty, bc.Difficulty)
	}
	if len(loaded.Blocks) != len(bc.Blocks) {
		t.Fatalf("loaded %d blocks, want %d", len(loaded.Blocks), len(bc.Blocks))
	}
	for i := range bc.Blocks {
		if loaded.Blocks[i] != bc.Blocks[i] {
			t.Fatalf("block %d = %+v, want %+v", i, loaded.Blocks[i], bc.Blocks[i])
		}
	}
}

func TestStep06_Load_MissingFile_IsAnError(t *testing.T) {
	_, err := persist.LoadChain(filepath.Join(t.TempDir(), "does-not-exist.json"))
	if err == nil {
		t.Fatal("LoadChain(missing file) = nil, want an error")
	}
	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("LoadChain(missing file) = %v, want an error wrapping os.ErrNotExist", err)
	}
}

func TestStep06_Load_RejectsATamperedFile(t *testing.T) {
	bc := buildChain()
	path := filepath.Join(t.TempDir(), "chain.json")
	if err := persist.SaveChain(path, bc); err != nil {
		t.Fatalf("SaveChain: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	cut := len(data) / 2 // flip a byte somewhere in the middle of the JSON
	data[cut] ^= 0xFF
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}

	if _, err := persist.LoadChain(path); err == nil {
		t.Fatal("LoadChain(tampered file) = nil, want an error")
	}
}

func TestStep06_Load_RejectsAChainWhoseDifficultyIsImpossible(t *testing.T) {
	bc := buildChain()
	path := filepath.Join(t.TempDir(), "chain.json")
	if err := persist.SaveChain(path, bc); err != nil {
		t.Fatalf("SaveChain: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	replaced := bytes.Replace(data, []byte(`"difficulty":2`), []byte(`"difficulty":64`), 1)
	if bytes.Equal(replaced, data) {
		t.Fatal("test setup: could not locate the difficulty field in the saved JSON (format changed?)")
	}
	if err := os.WriteFile(path, replaced, 0o600); err != nil {
		t.Fatal(err)
	}

	if _, err := persist.LoadChain(path); err == nil {
		t.Fatal("LoadChain with an impossible difficulty = nil, want an error")
	}
}

func TestStep06_Save_WritesAWholeChainOrNothing(t *testing.T) {
	bc := buildChain()
	dir := t.TempDir()
	path := filepath.Join(dir, "chain.json")

	if err := persist.SaveChain(path, bc); err != nil {
		t.Fatalf("SaveChain: %v", err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("the chain was not saved: %v", err)
	}
	if _, err := os.Stat(path + ".tmp"); !os.IsNotExist(err) {
		t.Fatal("a temp file was left behind after SaveChain: the rename must have consumed it")
	}

	var decoded struct {
		Difficulty int           `json:"difficulty"`
		Blocks     []chain.Block `json:"blocks"`
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("the saved file is not valid JSON: %v", err)
	}
	if decoded.Difficulty != 2 || len(decoded.Blocks) != 3 {
		t.Fatalf("saved chain = difficulty %d, %d blocks, want 2 and 3", decoded.Difficulty, len(decoded.Blocks))
	}
}

func TestStep06_Load_RoundTrip_IsDeterministic(t *testing.T) {
	bc := buildChain()
	path := filepath.Join(t.TempDir(), "chain.json")
	if err := persist.SaveChain(path, bc); err != nil {
		t.Fatalf("SaveChain: %v", err)
	}

	a, err := persist.LoadChain(path)
	if err != nil {
		t.Fatal(err)
	}
	b, err := persist.LoadChain(path)
	if err != nil {
		t.Fatal(err)
	}

	if len(a.Blocks) != len(b.Blocks) {
		t.Fatalf("two loads produced %d and %d blocks", len(a.Blocks), len(b.Blocks))
	}
	for i := range a.Blocks {
		if a.Blocks[i] != b.Blocks[i] {
			t.Fatalf("block %d differs between loads", i)
		}
	}
}

func TestStep06_SaveToAnUnwritablePath_IsAnError(t *testing.T) {
	bc := buildChain()
	missingDir := filepath.Join(t.TempDir(), "nope", "chain.json")

	err := persist.SaveChain(missingDir, bc)
	if err == nil {
		t.Fatal("SaveChain to a missing directory = nil, want an error")
	}
	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("SaveChain to a missing directory = %v, want an error wrapping os.ErrNotExist", err)
	}
}
