# Step 03 — The chain: genesis and validation

| | |
|---|---|
| **Build** | `internal/chain/blockchain.go` (given: sentinels + `Blockchain`); `internal/chain/ops.go` — `NewBlockchain`, `AddBlock`, `Validate` |
| **Given** | `Blockchain{Difficulty int; Blocks []Block}` — `Blocks` is exported so tests can tamper |
| **Tests** | `internal/chain/blockchain_test.go` — 10 tests |
| **Run** | `go test ./internal/chain -run Step03 -v` |
| **You'll learn** | contiguous indexing, link validation, the genesis rule, whole-chain verification |
| **Time** | ~60–90 min |

## 1. The story

A block is nothing without its place in a chain. Step 03 gives ChainCoin its **spine**: the genesis block and a `Validate()` that re-checks every block — hashes, indices, links — unconditionally.

This is the function a node runs (a) on startup before trusting anything, and (b) on every block it receives from the network. Get it wrong, and a malicious peer can feed you garbage. Get it right, and "trust the math, not the peer" becomes literal.

## 2. Concept corner

- **The genesis rule**: block 0's `PrevHash` is sixty-four zeros — a magic constant that means "there is nothing before this". Any other value → `ErrBrokenLink`.
- **Contiguous indices**: `Index` must equal the position from 0, in order. Gaps and duplicates → `ErrWrongIndex`.
- **Every block re-checks its own hash**: hash-vs-contents → `ErrInvalidHash`.
- **Validation order matters**: check hash integrity first, then index, then links — otherwise errors get attributed to the wrong field.
- **Difficulty is not checked here**: `Validate` checks structure. Proof-of-work validity is Step 02's job (`pow.Verify`) — used when importing a chain (Step 06). The separation keeps `chain` free of an import cycle with `pow`.

## 3. Requirements

| Function | Behaviour |
|---|---|
| `NewBlockchain(genesisData string, genesisTimestamp int64, difficulty int) *Blockchain` | Creates the chain with exactly one block: index 0, `PrevHash = GenesisPrevHash`, mined to the difficulty, carrying `genesisData`/`genesisTimestamp`. |
| `(bc *Blockchain) AddBlock(data string, timestamp int64) Block` | Appends a block mined with `bc.Difficulty` to the current last block; Index = len so far; `PrevHash` = the current last block's Hash. Returns the new block. |
| `(bc *Blockchain) Validate() error` | For every block: Hash matches contents (`ErrInvalidHash`); Index equals position (`ErrWrongIndex`); genesis links to sixty-four zeros and every other block links to its predecessor (`ErrBrokenLink`). Wrap the sentinel with `%w`. Empty chain → nil. |

## 4. What the tests pin down

- New chains start with one mined genesis block (Index 0, PrevHash = 64 zeros, Difficulty carried over).
- `AddBlock` appends correctly-linked, mined blocks with contiguous indices.
- `Validate` accepts a freshly built chain, an empty chain.
- Tampering with `Data` → `ErrInvalidHash`; reassigning `Index` (and recomputing Hash) → `ErrWrongIndex`; re-pointing `PrevHash` → `ErrBrokenLink`; breaking the genesis link → `ErrBrokenLink` — in each case the error satisfies `errors.Is(err, sentinel)`.

## 5. Hints (open one at a time)

<details><summary>Hint 1 — NewBlockchain</summary>

Build the genesis block, mine it, return `&Blockchain{Difficulty: difficulty, Blocks: []Block{g}}`.

</details>

<details><summary>Hint 2 — AddBlock</summary>

`prev` may be absent (empty chain): use `GenesisPrevHash` and Index -1 in that case. Otherwise take `Blocks[len-1]`.

</details>

<details><summary>Hint 3 — mine locally</summary>

You need the mining loop from Step 02 here; either duplicate the 6 lines (a `meetsZeros` helper) or call `pow.Mine` — but note `pow` imports `chain`, so **chain cannot import pow** (Go forbids import cycles).

</details>

<details><summary>Hint 4 — Validate</summary>

Loop, check hash → check index → check link (genesis special-case), wrap any failure with `fmt.Errorf("block %d: %w", i, sentinel)`.

</details>

## 6. Stretch challenges

1. Export the error detail: return exactly which block index failed, not just the sentinel.
2. **Cross-check PoW:** in `Validate`, also ensure each block was actually mined (with Step 02). Note: this check must come last and fails the same sentinel family.
3. **Merkle root:** set `Data` = sha256 of concatenated transaction IDs (Step 04+). Now Validate proves *transaction integrity*, not just block shape.
4. **Fuzz:** write a test that corrupts one random byte of one random block and demands Validate() errors. (Every single-bit corruption must be caught.)

## 7. Reflect

- Why is ordering (hash → index → link) important in Validate?
- Why can a hash-mismatch be detected *without* the network, but a broken link also requires the previous block?
- What happens if two blocks are swapped wholesale (both recomputed with correct hashes)? Does Validate catch it? (No — that's what proof of work and fork choice are for; Step 02 / future work.)
- Why is `Blocks` exported as a field here? What would seal it?

## 8. Done when

`go test ./internal/chain -run Step03 -race -v` → **ok**. Then `go vet ./...` and `gofmt -l .`, and compare with `_solutions/step03/…`.

## 9. Next

Step 04 — **transactions**: ECDSA keys, PKIX addresses, real signatures, and a mempool of pending transactions.
