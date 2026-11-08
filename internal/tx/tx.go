// Package tx implements ChainCoin's transactions: who is paying whom, signed with real ECDSA keys.
package tx

import (
	"crypto/ecdsa"
	"errors"
)

var (
	// ErrBadAmount: the amount is zero or negative. A transaction must move a positive number of coins. GIVEN.
	ErrBadAmount = errors.New("tx: amount must be positive")
	// ErrBadSignature: the signature does not verify. The transaction is not from the claimed sender, or it was tampered with. GIVEN.
	ErrBadSignature = errors.New("tx: signature does not verify")
)

// Transaction is a coin transfer.
//
// From and To are hex-encoded uncompressed SEC P-256 public keys (as produced by Address).
// The signature covers everything except itself (From, To, Amount, Timestamp) -- sign THIS, never the signature.
type Transaction struct {
	From      string `json:"from"`
	To        string `json:"to"`
	Amount    int64  `json:"amount"`
	Timestamp int64  `json:"timestamp"`
	Signature []byte `json:"signature"`
}

// KeyPair is an ECDSA P-256 key pair. The public key is what everyone else sees (the address); the private key never leaves the sender's machine.
type KeyPair struct {
	Private *ecdsa.PrivateKey
	Public  *ecdsa.PublicKey
}

// Mempool holds transactions waiting to be mined into a block.
type Mempool struct {
	txs []Transaction
}

// NewMempool returns a fresh, empty mempool. GIVEN.
func NewMempool() *Mempool {
	return &Mempool{}
}
