package fork_test

import (
	"errors"
	"strings"
	"testing"

	"blockchain/internal/chain"
	"blockchain/internal/fork"
)

// Step 07 -- fork choice.   Run:  go test ./internal/fork -run Step07 -v

func buildFork(t *testing.T, difficulty int, blocks ...string) *chain.Blockchain {
	t.Helper()
	bc := chain.NewBlockchain("genesis", 1_000_000_000, difficulty)
	for i, data := range blocks {
		bc.AddBlock(data, 1_000_000_060+int64(i)*60)
	}
	if w := len(blocks) + 1; len(bc.Blocks) != w {
		t.Fatalf("could not build the chain (got %d blocks, want %d): is step 03 solved?", len(bc.Blocks), w)
	}
	return bc
}

func TestStep07_ForkChoice_LongestValidChainWins(t *testing.T) {
	short := buildFork(t, 2) // genesis only
	long := buildFork(t, 2, "one", "two")

	picked, err := fork.ForkChoice([]*chain.Blockchain{short, long})
	if err != nil {
		t.Fatalf("ForkChoice: %v", err)
	}
	if picked != long {
		t.Fatalf("picked the %d-block chain, want the %d-block chain (greatest work)", len(picked.Blocks), len(long.Blocks))
	}
}

func TestStep07_ForkChoice_InvalidLongerChainLoses_ToAShorterValidOne(t *testing.T) {
	short := buildFork(t, 2, "mine")
	forged := buildFork(t, 2, "one", "two")
	forged.Blocks[1].Data = "forged in the dark" // hash no longer matches its contents

	picked, err := fork.ForkChoice([]*chain.Blockchain{short, forged})
	if err != nil {
		t.Fatalf("ForkChoice: %v", err)
	}
	if picked != short {
		t.Fatal("the forged longer chain must not beat the honest shorter one")
	}
}

func TestStep07_ForkChoice_GreaterWorkWins_EvenWithFewerBlocks(t *testing.T) {
	heavy := buildFork(t, 3)           // 1 block at difficulty 3 = 4096 units of work
	light := buildFork(t, 2, "a", "b") // 3 blocks at difficulty 2 = 768 units of work

	picked, err := fork.ForkChoice([]*chain.Blockchain{heavy, light})
	if err != nil {
		t.Fatalf("ForkChoice: %v", err)
	}
	if picked != heavy {
		t.Fatal("the harder-to-mint chain wins: difficulty is cumulative work, not block count")
	}
}

func TestStep07_ForkChoice_EqualWork_PicksLowerTipHash(t *testing.T) {
	x := buildFork(t, 2, "alpha")
	y := buildFork(t, 2, "beta")

	picked, err := fork.ForkChoice([]*chain.Blockchain{x, y})
	if err != nil {
		t.Fatalf("ForkChoice: %v", err)
	}

	want := x
	if y.Blocks[1].Hash < x.Blocks[1].Hash {
		want = y
	}
	if picked != want {
		t.Fatalf("equal work must resolve deterministically: picked the wrong chain (want lower tip hash)")
	}
}

func TestStep07_ForkChoice_EmptyInput_ReturnsErrNoCandidates(t *testing.T) {
	_, err := fork.ForkChoice(nil)
	if !errors.Is(err, fork.ErrNoCandidates) {
		t.Fatalf("ForkChoice(nil) = %v, want fork.ErrNoCandidates", err)
	}
}

func TestStep07_ForkChoice_AllInvalid_ReturnsErrNoValidChain(t *testing.T) {
	bad1 := buildFork(t, 2, "one")
	bad1.Blocks[1].Data = "x"

	bad2 := buildFork(t, 2, "one", "two")
	bad2.Blocks[2].Hash = strings.Repeat("00", 32) // a hash that matches no content

	_, err := fork.ForkChoice([]*chain.Blockchain{bad1, bad2})
	if !errors.Is(err, fork.ErrNoValidChain) {
		t.Fatalf("ForkChoice(all invalid) = %v, want fork.ErrNoValidChain", err)
	}
}

func TestStep07_ForkChoice_NilChainsAreSkipped(t *testing.T) {
	good := buildFork(t, 2)

	picked, err := fork.ForkChoice([]*chain.Blockchain{nil, good})
	if err != nil {
		t.Fatalf("ForkChoice: %v", err)
	}
	if picked != good {
		t.Fatal("a nil entry must not win, merely be skipped")
	}
}

func TestStep07_ForkChoice_ImpossibleDifficultyChainIsDiscarded(t *testing.T) {
	valid := buildFork(t, 2, "honest")
	parsedAsMagic := buildFork(t, 2, "impossible")
	parsedAsMagic.Difficulty = 65 // by construction nothing can ever meet a 65-hex-character proof of work

	picked, err := fork.ForkChoice([]*chain.Blockchain{valid, parsedAsMagic})
	if err != nil {
		t.Fatalf("ForkChoice: %v", err)
	}
	if picked != valid {
		t.Fatal("a chain that could never have been mined must be discarded")
	}
}
