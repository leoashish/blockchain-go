package tx

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/binary"
	"encoding/hex"
	"math/big"
	"strconv"
)

// REFERENCE SOLUTION for Step 04. Try it yourself first!

func GenerateKeyPair() (KeyPair, error) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return KeyPair{}, err
	}
	return KeyPair{Private: priv, Public: &priv.PublicKey}, nil
}

func Address(pub *ecdsa.PublicKey) string {
	der, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return "" // unreachable: a valid P-256 public key always marshals
	}
	return hex.EncodeToString(der)
}

func parseAddress(address string) (*ecdsa.PublicKey, error) {
	der, err := hex.DecodeString(address)
	if err != nil {
		return nil, err
	}
	pubIface, err := x509.ParsePKIXPublicKey(der)
	if err != nil {
		return nil, err
	}
	pub, ok := pubIface.(*ecdsa.PublicKey)
	if !ok {
		return nil, ErrBadSignature
	}
	return pub, nil
}

func SigningBytes(t Transaction) []byte {
	h := sha256.New()
	fields := []string{
		t.From,
		t.To,
		strconv.FormatInt(t.Amount, 10),
		strconv.FormatInt(t.Timestamp, 10),
	}
	var lenbuf [binary.MaxVarintLen64]byte
	for _, field := range fields {
		n := binary.PutUvarint(lenbuf[:], uint64(len(field)))
		h.Write(lenbuf[:n])
		h.Write([]byte(field))
	}
	return h.Sum(nil)
}

func Sign(t Transaction, priv *ecdsa.PrivateKey) (Transaction, error) {
	digest := SigningBytes(t)
	r, s, err := ecdsa.Sign(rand.Reader, priv, digest)
	if err != nil {
		return Transaction{}, err
	}
	rBytes := r.Bytes()
	sBytes := s.Bytes()
	sig := make([]byte, 64)
	copy(sig[32-len(rBytes):32], rBytes)
	copy(sig[64-len(sBytes):64], sBytes)
	t.Signature = sig
	return t, nil
}

func Verify(t Transaction) bool {
	if len(t.Signature) != 64 {
		return false
	}
	pub, err := parseAddress(t.From)
	if err != nil {
		return false
	}
	r := new(big.Int).SetBytes(t.Signature[:32])
	s := new(big.Int).SetBytes(t.Signature[32:])
	return ecdsa.Verify(pub, SigningBytes(t), r, s)
}

func (m *Mempool) Add(t Transaction) error {
	if t.Amount <= 0 {
		return ErrBadAmount
	}
	if !Verify(t) {
		return ErrBadSignature
	}
	m.txs = append(m.txs, t)
	return nil
}

func (m *Mempool) Pending() []Transaction {
	return append([]Transaction(nil), m.txs...)
}

func (m *Mempool) Drain() []Transaction {
	drained := m.txs
	m.txs = nil
	return drained
}
