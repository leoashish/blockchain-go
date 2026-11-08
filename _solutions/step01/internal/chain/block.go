package chain

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"strconv"
)

// REFERENCE SOLUTION for Step 01. Try it yourself first!

func CalculateHash(b Block) string {
	h := sha256.New()
	fields := []string{
		strconv.Itoa(b.Index),
		strconv.FormatInt(b.Timestamp, 10),
		b.Data,
		b.PrevHash,
		strconv.FormatInt(b.Nonce, 10),
	}
	var lenbuf [binary.MaxVarintLen64]byte
	for _, field := range fields {
		n := binary.PutUvarint(lenbuf[:], uint64(len(field)))
		h.Write(lenbuf[:n])
		h.Write([]byte(field))
	}
	return hex.EncodeToString(h.Sum(nil))
}

func NewBlock(index int, data string, prevHash string, timestamp int64) Block {
	b := Block{Index: index, Timestamp: timestamp, Data: data, PrevHash: prevHash}
	b.Hash = CalculateHash(b)
	return b
}

func NewGenesisBlock(data string, timestamp int64) Block {
	return NewBlock(0, data, GenesisPrevHash, timestamp)
}

func VerifyHash(b Block) bool {
	return b.Hash == CalculateHash(b)
}
