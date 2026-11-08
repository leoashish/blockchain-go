package utxo_test

import (
	"errors"
	"testing"

	"blockchain/internal/utxo"
)

// Step 05 -- the UTXO set.   Run:  go test ./internal/utxo -run Step05 -v

func TestStep05_Coinbase_MintsAndTracksAnOutput(t *testing.T) {
	l := utxo.NewLedger()

	if err := l.Coinbase("reward-1", "alice", 1000); err != nil {
		t.Fatalf("Coinbase: %v", err)
	}
	if got := l.Balance("alice"); got != 1000 {
		t.Fatalf("Balance(alice) = %d, want 1000", got)
	}
	outputs := l.Outputs("alice")
	if len(outputs) != 1 || outputs[0].ID != "reward-1" || outputs[0].Amount != 1000 {
		t.Fatalf("Outputs(alice) = %+v, want the reward-1 output of 1000", outputs)
	}
}

func TestStep05_Coinbase_RejectsBadAmountsAndDuplicateIDs(t *testing.T) {
	l := utxo.NewLedger()

	if err := l.Coinbase("free-1", "alice", 0); !errors.Is(err, utxo.ErrBadAmount) {
		t.Fatalf("Coinbase(0) = %v, want utxo.ErrBadAmount", err)
	}
	if err := l.Coinbase("free-2", "alice", -5); !errors.Is(err, utxo.ErrBadAmount) {
		t.Fatalf("Coinbase(-5) = %v, want utxo.ErrBadAmount", err)
	}

	if err := l.Coinbase("dup", "alice", 100); err != nil {
		t.Fatalf("first Coinbase(dup): %v", err)
	}
	if err := l.Coinbase("dup", "bob", 200); !errors.Is(err, utxo.ErrDuplicateOutput) {
		t.Fatalf("second Coinbase(dup) = %v, want utxo.ErrDuplicateOutput", err)
	}
	if got := l.Balance("alice"); got != 100 {
		t.Fatalf("Balance(alice) = %d, want 100 (the first output only)", got)
	}
}

func TestStep05_Outputs_AreIsolatedPerOwner_AndSorted(t *testing.T) {
	l := utxo.NewLedger()
	_ = l.Coinbase("c", "alice", 300)
	_ = l.Coinbase("a", "alice", 100)
	_ = l.Coinbase("b", "bob", 200)

	outputs := l.Outputs("alice")
	if len(outputs) != 2 {
		t.Fatalf("Outputs(alice) = %+v, want 2 (bob's output excluded)", outputs)
	}
	if outputs[0].ID != "a" || outputs[1].ID != "c" {
		t.Fatalf("Outputs(alice) ids = %v, %v, want a < c (sorted)", outputs[0].ID, outputs[1].ID)
	}
	if got := l.Balance("bob"); got != 200 {
		t.Fatalf("Balance(bob) = %d, want 200", got)
	}
}

func TestStep05_Transfer_ConsumesInputsAndCreatesPaymentAndChange(t *testing.T) {
	l := utxo.NewLedger()
	_ = l.Coinbase("coinbase-a", "alice", 1000)

	if err := l.Transfer("tx1", "alice", "bob", 400); err != nil {
		t.Fatalf("Transfer: %v", err)
	}

	if got := l.Balance("alice"); got != 600 {
		t.Fatalf("Balance(alice) = %d, want 600 (1000 - 400)", got)
	}
	if got := l.Balance("bob"); got != 400 {
		t.Fatalf("Balance(bob) = %d, want 400", got)
	}

	bobOutputs := l.Outputs("bob")
	if len(bobOutputs) != 1 || bobOutputs[0].ID != "tx1:0" || bobOutputs[0].Amount != 400 {
		t.Fatalf("bob's outputs = %+v, want tx1:0 of 400", bobOutputs)
	}

	aliceOutputs := l.Outputs("alice")
	if len(aliceOutputs) != 1 || aliceOutputs[0].ID != "tx1:1" || aliceOutputs[0].Amount != 600 {
		t.Fatalf("alice's outputs = %+v, want tx1:1 of 600 (change)", aliceOutputs)
	}
	for _, o := range aliceOutputs {
		if o.ID == "coinbase-a" {
			t.Fatal("the consumed coinbase-a output is still in the ledger")
		}
	}
}

func TestStep05_Transfer_ExactAmountLeavesNoChange(t *testing.T) {
	l := utxo.NewLedger()
	_ = l.Coinbase("coinbase-b", "alice", 500)

	if err := l.Transfer("tx2", "alice", "bob", 500); err != nil {
		t.Fatalf("Transfer: %v", err)
	}
	if got := l.Balance("alice"); got != 0 {
		t.Fatalf("Balance(alice) = %d, want 0 (everything moved)", got)
	}
	outputs := l.Outputs("alice")
	if len(outputs) != 0 {
		t.Fatalf("alice has %d outputs left, want 0 (no change output)", len(outputs))
	}
}

func TestStep05_Transfer_InsufficientFundsChangesNothing(t *testing.T) {
	l := utxo.NewLedger()
	_ = l.Coinbase("fund-1", "alice", 100)

	err := l.Transfer("tx3", "alice", "bob", 150)
	if !errors.Is(err, utxo.ErrInsufficientFunds) {
		t.Fatalf("Transfer(insufficient) = %v, want utxo.ErrInsufficientFunds", err)
	}

	if got := l.Balance("alice"); got != 100 {
		t.Fatalf("Balance(alice) = %d, want 100 (nothing consumed)", got)
	}
	if got := len(l.Outputs("alice")); got != 1 {
		t.Fatalf("alice has %d outputs, want 1 (the original, untouched)", got)
	}
	if got := l.Balance("bob"); got != 0 {
		t.Fatalf("Balance(bob) = %d, want 0 (payment never happened)", got)
	}
}

func TestStep05_Transfer_ZeroAmountMovesNothing(t *testing.T) {
	l := utxo.NewLedger()
	_ = l.Coinbase("fund-2", "alice", 100)

	if err := l.Transfer("tx4", "alice", "bob", 0); !errors.Is(err, utxo.ErrBadAmount) {
		t.Fatalf("Transfer(0) = %v, want utxo.ErrBadAmount", err)
	}

	if got := l.Balance("alice"); got != 100 {
		t.Fatalf("Balance(alice) = %d, want 100 (untouched)", got)
	}
}

func TestStep05_Transfer_SelectsAscendingOutputsDeterministically(t *testing.T) {
	l := utxo.NewLedger()
	_ = l.Coinbase("o1", "alice", 100)
	_ = l.Coinbase("o2", "alice", 250)
	_ = l.Coinbase("o3", "alice", 100)

	if err := l.Transfer("tx5", "alice", "bob", 150); err != nil {
		t.Fatalf("Transfer: %v", err)
	}

	// Ascending (amount, id): o1(100), o3(100), o2(250). A 150 transfer consumes o1 and o3 (200).
	// Alice keeps o2 (250) and the change (50): 300 in total.
	if got := l.Balance("alice"); got != 300 {
		t.Fatalf("Balance(alice) = %d, want 300 (o2 = 250 and the change tx5:1 = 50)", got)
	}
	if got := l.Outputs("alice"); len(got) != 2 || got[0].ID != "o2" {
		t.Fatalf("alice's outputs = %+v, want o2 (unspent) and tx5:1 (change)", got)
	}
	for _, id := range []string{"o1", "o3"} {
		_, err := l.Spend(id)
		if !errors.Is(err, utxo.ErrOutputNotFound) {
			t.Fatalf("Spend(%s) after consuming it = %v, want utxo.ErrOutputNotFound", id, err)
		}
	}
}

func TestStep05_Transfer_RefusesWhenOutputIDsAreTaken_BeforeConsumingAnything(t *testing.T) {
	l := utxo.NewLedger()
	_ = l.Coinbase("fund-3", "alice", 1000)
	_ = l.Coinbase("tx6:0", "alice", 5) // pre-claim the payment output id

	err := l.Transfer("tx6", "alice", "bob", 400)
	if !errors.Is(err, utxo.ErrDuplicateOutput) {
		t.Fatalf("Transfer with a taken id = %v, want utxo.ErrDuplicateOutput", err)
	}

	if got := l.Balance("alice"); got != 1005 {
		t.Fatalf("Balance(alice) = %d, want 1005 (nothing consumed)", got)
	}
}

func TestStep05_Spend_RemovesAndReturns_OneOutputOnly(t *testing.T) {
	l := utxo.NewLedger()
	_ = l.Coinbase("s1", "alice", 70)

	spent, err := l.Spend("s1")
	if err != nil {
		t.Fatalf("Spend: %v", err)
	}
	if spent.ID != "s1" || spent.Owner != "alice" || spent.Amount != 70 {
		t.Fatalf("Spend returned %+v, want s1/alice/70", spent)
	}
	if got := l.Balance("alice"); got != 0 {
		t.Fatalf("Balance(alice) = %d, want 0 after spending", got)
	}

	if _, err := l.Spend("s1"); !errors.Is(err, utxo.ErrOutputNotFound) {
		t.Fatalf("second Spend(s1) = %v, want utxo.ErrOutputNotFound", err)
	}
	if _, err := l.Spend("never-existed"); !errors.Is(err, utxo.ErrOutputNotFound) {
		t.Fatalf("Spend(missing) = %v, want utxo.ErrOutputNotFound", err)
	}
}
