// Package chain implements ChainCoin's blocks and the hashing that links them.
// Step 01: the block itself -- SHA-256 hashing and how blocks link to their predecessors.
package chain

// GenesisPrevHash is the "previous hash" of the genesis block: sixty-four zeros. GIVEN.
const GenesisPrevHash = "0000000000000000000000000000000000000000000000000000000000000000"

// Block is one block in the chain.
//
// STEP 01: a block is just memory right now -- index, timestamp, data, the hash of the previous block, a nonce (unused until
// Step 02's proof of work) and the hash of the identity fields themselves. The interesting question this step answers:
// *what goes into the hash, and what does not?* The Hash field itself cannot, obviously -- so the interesting bit is that
// the hash covers the block's IDENTITY (index, timestamp, data, prevHash, nonce) and excludes nothing else.
type Block struct {
	Index     int    `json:"index"`
	Timestamp int64  `json:"timestamp"`
	Data      string `json:"data"`
	PrevHash  string `json:"prev_hash"`
	Nonce     int64  `json:"nonce"`
	Hash      string `json:"hash"`
}
