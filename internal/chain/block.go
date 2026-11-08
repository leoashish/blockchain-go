package chain

// STEP 01: implement the block. Read docs/steps/step-01-the-block.md first.   Run:  go test ./internal/chain -run Step01

// CalculateHash computes the SHA-256 hash of the block's identity fields -- Index, Timestamp, Data, PrevHash, Nonce -- and
// returns it as 64 lowercase hex characters. The Hash field itself must NOT be part of the input: it cannot hash itself.
// Hashing implies the fields must be encoded unambiguously: no two different blocks may hash to the same input bytes.
func CalculateHash(b Block) string {
	return "" // TODO Step 01
}

// NewBlock creates a block with the given parts and sets its Hash to CalculateHash of the parts. Nonce starts at 0.
func NewBlock(index int, data string, prevHash string, timestamp int64) Block {
	return Block{} // TODO Step 01
}

// NewGenesisBlock creates block 0 (the genesis block): index 0, PrevHash = GenesisPrevHash, hash computed from the parts.
func NewGenesisBlock(data string, timestamp int64) Block {
	return Block{} // TODO Step 01
}

// VerifyHash reports whether b.Hash equals CalculateHash(b) -- i.e. the block has not been tampered with.
func VerifyHash(b Block) bool {
	return false // TODO Step 01
}
