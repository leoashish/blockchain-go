package tx

import (
	"crypto/ecdsa"

	"blockchain/internal/todo"
)

// STEP 04: implement transactions. Read docs/steps/step-04-transactions.md first.   Run:  go test ./internal/tx -run Step04

// GenerateKeyPair creates a fresh ECDSA P-256 key pair using crypto/rand.
func GenerateKeyPair() (KeyPair, error) {
	return KeyPair{}, todo.NotImplemented("Step 04: GenerateKeyPair")
}

// Address returns the hex-encoded PKIX encoding of the public key.
func Address(pub *ecdsa.PublicKey) string {
	return "" // TODO Step 04
}

// SigningBytes returns the canonical byte string to sign / verify. Everything except the signature goes in:
// length-prefixed From, To, Amount and Timestamp, hashed with SHA-256.
func SigningBytes(t Transaction) []byte {
	return nil // TODO Step 04
}

// Sign signs t with the private key (ECDSA P-256 + SHA-256) and returns the transaction with its Signature set (64 bytes: r ‖ s).
func Sign(t Transaction, priv *ecdsa.PrivateKey) (Transaction, error) {
	return Transaction{}, todo.NotImplemented("Step 04: Sign")
}

// Verify reports whether t's signature is valid for t's claimed From address. No panics on garbage:
// malformed addresses or signatures simply fail to verify.
func Verify(t Transaction) bool {
	return false // TODO Step 04
}

// Add validates a transaction (positive amount; signature verifies) and queues it for mining.
func (m *Mempool) Add(t Transaction) error {
	return todo.NotImplemented("Step 04: Mempool.Add")
}

// Pending returns a copy of the queued transactions, oldest first.
func (m *Mempool) Pending() []Transaction {
	return nil // TODO Step 04
}

// Drain returns all queued transactions and empties the mempool (as a miner would once they are in a block).
func (m *Mempool) Drain() []Transaction {
	return nil // TODO Step 04
}
