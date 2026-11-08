package utxo

import "blockchain/internal/todo"

// STEP 05: implement the UTXO set. Read docs/steps/step-05-utxo.md first.   Run:  go test ./internal/utxo -run Step05

// Coinbase mints a brand-new output out of thin air (as mining does): id must be free, amount must be positive.
func (l *Ledger) Coinbase(id, owner string, amount int64) error {
	return todo.NotImplemented("Step 05: Ledger.Coinbase")
}

// Outputs returns every unspent output owned by the address, sorted by ID for determinism.
func (l *Ledger) Outputs(owner string) []Output {
	return nil // TODO Step 05
}

// Balance returns the total unspent amount owned by the address.
func (l *Ledger) Balance(owner string) int64 {
	return 0 // TODO Step 05
}

// Transfer moves amount from `from` to `to`. It selects `from`'s unspent outputs (ascending amount, then ID) until the
// amount is covered, consumes them, and creates:
//   - output id+":0" owned by `to` for exactly `amount`
//   - output id+":1" owned by `from` for the change, if any
//
// Errors: ErrBadAmount (amount <= 0), ErrDuplicateOutput (id+":0" or id+":1" already exists), ErrInsufficientFunds
// (outputs do not cover the amount). On any error NOTHING is consumed and NOTHING is created.
func (l *Ledger) Transfer(id, from, to string, amount int64) error {
	return todo.NotImplemented("Step 05: Ledger.Transfer")
}

// Spend removes and returns one output by id. A spent output is gone from the ledger forever.
func (l *Ledger) Spend(id string) (Output, error) {
	return Output{}, todo.NotImplemented("Step 05: Ledger.Spend")
}
