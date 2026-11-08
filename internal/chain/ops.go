package chain

import "blockchain/internal/todo"

// STEP 03: implement the chain. Read docs/steps/step-03-the-chain.md first.   Run:  go test ./internal/chain -run Step03

// NewBlockchain creates a Blockchain with one block: the genesis block (index 0, PrevHash = GenesisPrevHash), mined with
// the given difficulty. The genesis block carries genesisData and genesisTimestamp.
func NewBlockchain(genesisData string, genesisTimestamp int64, difficulty int) *Blockchain {
	return &Blockchain{} // TODO Step 03
}

// AddBlock appends a new block carrying data (mined with the chain's difficulty) after the current last block, and returns
// it. The new block's Index is len(Blocks) so far, and its PrevHash is the current last block's Hash.
func (bc *Blockchain) AddBlock(data string, timestamp int64) Block {
	return Block{} // TODO Step 03
}

// Validate checks every block: Hash matches contents, Index is contiguous from 0, PrevHash links to the previous block
// (the genesis block must link to GenesisPrevHash). Any failure returns an error wrapping the matching sentinel.
func (bc *Blockchain) Validate() error {
	return todo.NotImplemented("Step 03: Blockchain.Validate")
}
