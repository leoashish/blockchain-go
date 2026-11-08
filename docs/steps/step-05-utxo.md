# Step 05 — The UTXO set: balances from the ledger alone

| | |
|---|---|
| **Build** | `internal/utxo/ops.go` — `Ledger.{Coinbase, Outputs, Balance, Transfer, Spend}` |
| **Given** | `Output`, `Ledger` (`unspent` map), `ErrBadAmount`, `ErrOutputNotFound`, `ErrInsufficientFunds`, `ErrDuplicateOutput`, `NewLedger` (`internal/utxo/utxo.go`) |
| **Tests** | `internal/utxo/utxo_test.go` — 10 tests |
| **Run** | `go test ./internal/utxo -run Step05 -race -v` |
| **You'll learn** | the UTXO model, fee/change outputs, deterministic selection, atomic ledger transitions |
| **Time** | ~90 min |

## 1. The story

ChainCoin now has blocks, proof of work, transactions and key pairs. What it does NOT yet have: *who owns what*.

Blockchains do not maintain account balances — they maintain **unspent outputs**. Every transaction consumes outputs and creates new ones. A balance is just the sum of the unspent outputs an address owns. There is no such thing as a balance row in Bitcoin — and that is a beautiful, unusual design.

The UTXO set ("unspent transaction output set") is what a node keeps in memory so it does not re-walk the entire chain to answer "what does Alice own?".

## 2. Concept corner

- **Outputs, not balances**: coins exist *as* outputs. To pay someone, you **consume** outputs you own and **create** new outputs. There is never a negative UTXO or a balance column.
- **Consume whole outputs**: a transaction that wants to spend 30 coins from a 100-coin output consumes the entire 100, and sends **70 back to the payer as change**. (Omitting change = donating the difference to the miner — a fun bug to relive on mainnet.)
- **Coinbase**: the one transaction allowed to create coins from nothing — the miner's reward. We model it as a special case (`Coinbase`); on a real network, a coinbase is simply the first transaction of each block with exactly one input ("the promise of new coins").
- **Atomicity**: selection → consumption → creation must be all-or-nothing. Our `Transfer` checks everything up front and consumes only at the end — any failure (`ErrInsufficientFunds`, duplicate id…) leaves the ledger untouched.
- **Deterministic selection**: outputs are selected in ascending `(amount, id)` order — minimal change, maximal reproducibility, zero map-iteration nondeterminism in tests.

## 3. Requirements

| Function | Behaviour |
|---|---|
| `Coinbase(id, owner string, amount int64) error` | Mint a new output. `amount <= 0` or empty `id` → `ErrBadAmount`. `id` already exists → `ErrDuplicateOutput`. |
| `Outputs(owner string) []Output` | Every unspent output owned by `owner`, **sorted by ID**. Always non-nil (empty slice when none). |
| `Balance(owner string) int64` | Sum of the owner's unspent outputs' amounts. |
| `Transfer(id, from, to string, amount int64) error` | `amount <= 0` → `ErrBadAmount`. Pre-create ids taken (`id+":0"`, `id+":1"`) → `ErrDuplicateOutput`. Select `from`'s unspent outputs ascending `(amount, id)` until covered; consume them; create payment `id+":0"` for `amount` → `to`; create change `id+":1"` for the remainder → `from`. Insufficient → `ErrInsufficientFunds`, nothing consumed. On any error, nothing is consumed and nothing is created. |
| `Spend(id string) (Output, error)` | Remove and return one output. Unknown/spent id → `ErrOutputNotFound`. |

## 4. What the tests pin down

- Minting coins, tracking balances, per-owner isolation, sorted `Outputs`.
- Bad amounts (zero, negative), duplicate ids.
- Exact-consumption transfers create payment + change; exact-amount leaves no change.
- Insufficient funds: nothing is consumed, nothing created — balances unchanged before and after.
- Deterministic ascending selection (o1, o3 consumed; o2 untouched).
- Pre-taken ids: the whole transaction is refused before inputs are touched.
- Spend: single output removed, `ErrOutputNotFound` on re-spend/missing.

## 5. Hints (open one at a time)

<details><summary>Hint 1 — Outputs / Balance</summary>

Range the `unspent` map, filter by `Owner`, collect into a slice, sort by ID. Balance is the same walk with a running sum.

</details>

<details><summary>Hint 2 — Transfer: preflight</summary>

Reject `amount <= 0`; check `id+":0"` and `id+":1"` are free; select candidates sorted ascending by `(amount, id)`; if `sum(candidates) < amount` → `ErrInsufficientFunds` (nothing touched!).

</details>

<details><summary>Hint 3 — Transfer: consume and create</summary>

Delete the selected ids; create `id+":0"` → `to` for `amount`; if `change := covered - amount; change > 0` create `id+":1"` → `from`.

</details>

<details><summary>Hint 4 — Spend</summary>

Look up, delete, return. One line each — the UTXO model turns "deleting a row" into a citizen's whole financial life.

</details>

## 6. Stretch challenges

1. **Fees:** add `fee int64` to Transfer; change = covered − amount − fee. Where does the fee go?
2. **Dust:** refuse payments below a configurable dust threshold (UTXO bloat).
3. **Collision-safety proof:** show (on paper) that two of your tests can never produce the same output id by accident.
4. **Snapshotting:** serialise the whole ledger as JSON and rebuild it on load — the "state sync" half of a real full node.

## 7. Reflect

- Why does the UTXO model make "where did this coin come from?" answerable, but makes "who is Bob?" not?
- What is the relation between the UTXO set and a proof-of-work chain's difficulty?
- Why is change sent back to the *same* address (usually to a fresh one in reality)?
- Why must output creation happen after consumption is validated?

## 8. Done when

`go test ./internal/utxo -run Step05 -race -v` → **ok**. Then `go vet ./...` and `gofmt -l .`, and compare with `_solutions/step05/…`.

## 9. Next

Step 06 — **persistence**: save the chain to disk atomically, and load it back *verified* — a file that claims to be a chain is not a chain until it passes `pow.Verify` and `Validate`.
