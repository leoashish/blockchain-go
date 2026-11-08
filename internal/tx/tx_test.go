package tx_test

import (
	"errors"
	"strings"
	"testing"

	"blockchain/internal/tx"
)

// Step 04 -- transactions and signatures.   Run:  go test ./internal/tx -run Step04 -v

func keyPair(t *testing.T) tx.KeyPair {
	t.Helper()
	kp, err := tx.GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair: %v", err)
	}
	return kp
}

func signedTx(t *testing.T, sender, recipient tx.KeyPair, amount int64, timestamp int64) tx.Transaction {
	t.Helper()
	transaction := tx.Transaction{
		From:      tx.Address(sender.Public),
		To:        tx.Address(recipient.Public),
		Amount:    amount,
		Timestamp: timestamp,
	}
	signed, err := tx.Sign(transaction, sender.Private)
	if err != nil {
		t.Fatalf("Sign: %v", err)
	}
	return signed
}

func TestStep04_GenerateKeyPair_MakesAWorkingPair(t *testing.T) {
	kp, err := tx.GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair: %v", err)
	}
	if kp.Private == nil || kp.Public == nil {
		t.Fatal("the pair must contain both keys")
	}
	if !kp.Public.Equal(&kp.Private.PublicKey) {
		t.Fatal("the public half of the pair does not match the private key")
	}
}

func TestStep04_Address_IsAHexEncodedDecodableKey(t *testing.T) {
	kp := keyPair(t)
	addr := tx.Address(kp.Public)

	if addr == "" {
		t.Fatal("Address returned nothing")
	}
	for _, r := range addr {
		if !strings.ContainsRune("0123456789abcdef", r) {
			t.Fatalf("Address contains %q, want only lowercase hex", r)
		}
	}
	if tx.Address(kp.Public) != addr {
		t.Fatal("Address is not deterministic")
	}
}

func TestStep04_SigningBytes_IsDeterministicAndCoversEverythingButTheSignature(t *testing.T) {
	sender := keyPair(t)
	recipient := keyPair(t)
	base := tx.Transaction{From: tx.Address(sender.Public), To: tx.Address(recipient.Public), Amount: 100, Timestamp: 1700000000}

	first := tx.SigningBytes(base)
	second := tx.SigningBytes(base)
	if len(first) == 0 {
		t.Fatal("SigningBytes returned nothing")
	}
	if string(first) != string(second) {
		t.Fatal("SigningBytes is not deterministic")
	}

	variants := map[string]tx.Transaction{
		"from":      {From: tx.Address(recipient.Public), To: base.To, Amount: base.Amount, Timestamp: base.Timestamp},
		"to":        {From: base.From, To: base.From, Amount: base.Amount, Timestamp: base.Timestamp},
		"amount":    {From: base.From, To: base.To, Amount: base.Amount + 1, Timestamp: base.Timestamp},
		"timestamp": {From: base.From, To: base.To, Amount: base.Amount, Timestamp: base.Timestamp + 1},
	}
	for field, v := range variants {
		if string(tx.SigningBytes(v)) == string(first) {
			t.Errorf("changing %s did not change SigningBytes", field)
		}
	}

	signed := base
	signed.Signature = []byte("anything")
	if string(tx.SigningBytes(signed)) != string(first) {
		t.Fatal("SigningBytes must ignore the Signature field")
	}
}

func TestStep04_Sign_Then_Verify_RoundTrips(t *testing.T) {
	sender := keyPair(t)
	recipient := keyPair(t)

	signed := signedTx(t, sender, recipient, 50, 1700000000)
	if len(signed.Signature) != 64 {
		t.Fatalf("Signature = %d bytes, want 64 (r ‖ s)", len(signed.Signature))
	}
	if !tx.Verify(signed) {
		t.Fatal("Verify(signed transaction) = false, want true")
	}
}

func TestStep04_Verify_RejectsATamperedTransaction(t *testing.T) {
	sender := keyPair(t)
	recipient := keyPair(t)

	signed := signedTx(t, sender, recipient, 50, 1700000000)

	variants := map[string]tx.Transaction{
		"amount":    {From: signed.From, To: signed.To, Amount: signed.Amount + 1, Timestamp: signed.Timestamp, Signature: signed.Signature},
		"to":        {From: signed.From, To: signed.From, Amount: signed.Amount, Timestamp: signed.Timestamp, Signature: signed.Signature},
		"timestamp": {From: signed.From, To: signed.To, Amount: signed.Amount, Timestamp: signed.Timestamp + 1, Signature: signed.Signature},
		"from":      {From: tx.Address(recipient.Public), To: signed.To, Amount: signed.Amount, Timestamp: signed.Timestamp, Signature: signed.Signature},
	}
	for field, v := range variants {
		if tx.Verify(v) {
			t.Errorf("Verify(%s tampered) = true, want false", field)
		}
	}
}

func TestStep04_Verify_DoesNotPanicOnGarbage(t *testing.T) {
	// Starts from real keys, so it is red until GenerateKeyPair and Sign work.
	sender := keyPair(t)
	recipient := keyPair(t)

	signed := signedTx(t, sender, recipient, 10, 1700000000)

	garbage := signed
	garbage.From = "!!!not hex!!!"
	if tx.Verify(garbage) {
		t.Fatal("Verify(address that is not even hex) = true, want false")
	}

	shortSig := signed
	shortSig.Signature = signed.Signature[:32]
	if tx.Verify(shortSig) {
		t.Fatal("Verify(32-byte signature) = true, want false")
	}

	emptySig := signed
	emptySig.Signature = nil
	if tx.Verify(emptySig) {
		t.Fatal("Verify(transaction with no signature) = true, want false")
	}
}

func TestStep04_Mempool_AddAcceptsAValidAndRejectsInvalidTransactions(t *testing.T) {
	sender := keyPair(t)
	recipient := keyPair(t)

	mempool := tx.NewMempool()

	signed := signedTx(t, sender, recipient, 75, 1700000000)
	if err := mempool.Add(signed); err != nil {
		t.Fatalf("Add(valid) = %v, want nil", err)
	}
	if pending := mempool.Pending(); len(pending) != 1 {
		t.Fatalf("Pending = %d transactions, want 1", len(pending))
	}

	zero := signed
	zero.Amount = 0
	if err := mempool.Add(zero); !errors.Is(err, tx.ErrBadAmount) {
		t.Fatalf("Add(amount 0) = %v, want tx.ErrBadAmount", err)
	}

	forged := signed
	forged.Signature = make([]byte, 64)
	if err := mempool.Add(forged); !errors.Is(err, tx.ErrBadSignature) {
		t.Fatalf("Add(forged) = %v, want tx.ErrBadSignature", err)
	}

	if pending := mempool.Pending(); len(pending) != 1 {
		t.Fatalf("Pending = %d after rejects, want 1 (the valid one)", len(pending))
	}
}

func TestStep04_Mempool_DrainReturnsAllAndEmptiesThePool(t *testing.T) {
	sender := keyPair(t)
	recipient := keyPair(t)

	mempool := tx.NewMempool()
	for i := 0; i < 3; i++ {
		signed := signedTx(t, sender, recipient, int64(i+1), 1700000000+int64(i))
		if err := mempool.Add(signed); err != nil {
			t.Fatalf("Add: %v", err)
		}
	}

	drained := mempool.Drain()
	if len(drained) != 3 {
		t.Fatalf("Drain = %d transactions, want 3", len(drained))
	}
	if pending := mempool.Pending(); len(pending) != 0 {
		t.Fatalf("after Drain, Pending = %d, want 0", len(pending))
	}

	drainedAgain := mempool.Drain()
	if len(drainedAgain) != 0 {
		t.Fatalf("second Drain = %d transactions, want 0", len(drainedAgain))
	}
}

func TestStep04_Mempool_PendingIsACopyNotTheBackingSlice(t *testing.T) {
	sender := keyPair(t)
	recipient := keyPair(t)

	mempool := tx.NewMempool()
	signed := signedTx(t, sender, recipient, 10, 1700000000)
	if err := mempool.Add(signed); err != nil {
		t.Fatalf("Add: %v", err)
	}

	pending := mempool.Pending()
	pending[0].Amount = 999999

	if got := mempool.Pending(); got[0].Amount == 999999 {
		t.Fatal("Pending returned the backing slice: mutating it corrupted the pool")
	}
}
