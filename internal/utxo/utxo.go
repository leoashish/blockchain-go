// Package utxo implements the unspent transaction output (UTXO) model — the sharpest part of a cryptocurrency.
package utxo

import "errors"

var (
	// ErrBadAmount: the amount is zero or negative. GIVEN.
	ErrBadAmount = errors.New("utxo: amount must be positive")
	// ErrOutputNotFound: the output id does not exist, or it was already spent. GIVEN.
	ErrOutputNotFound = errors.New("utxo: output not found (or already spent)")
	// ErrInsufficientFunds: the selected outputs do not cover the amount. Nothing was consumed. GIVEN.
	ErrInsufficientFunds = errors.New("utxo: insufficient funds")
	// ErrDuplicateOutput: an output with this id already exists. GIVEN.
	ErrDuplicateOutput = errors.New("utxo: output id already exists")
)

// Output is one unspent transaction output: a slice of coins owned by Address.
type Output struct {
	ID     string `json:"id"`
	Owner  string `json:"owner"`
	Amount int64  `json:"amount"`
}

// Ledger tracks every unspent output ever created by the chain. Spending = removing outputs; funding = adding them.
type Ledger struct {
	unspent map[string]Output
}

// NewLedger returns an empty ledger. GIVEN.
func NewLedger() *Ledger {
	return &Ledger{unspent: make(map[string]Output)}
}
