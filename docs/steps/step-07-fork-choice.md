# Step 07 — Fork choice: most work wins

| | |
|---|---|
| **Build** | `internal/fork/ops.go` — `ForkChoice` |
| **Given** | `ErrNoCandidates`, `ErrNoValidChain` (`internal/fork/fork.go`); `chain.Blockchain`, `chain.Validate`, `pow.Verify` (steps 01–03) |
| **Tests** | `internal/fork/fork_test.go` — 8 tests |
| **Run** | `go test ./internal/fork -run Step07 -race -v` |
| **You'll learn** | fork resolution, cumulative work, deterministic tie-breaking, rejection of impossible difficulty declarations |
| **Time** | ~60–90 min |

## 1. The story

Two miners found the same block height at the same time. Both chains are self-contained, both pass validation — which chain is THE chain?

This is the last piece of honest Nakamoto consensus: **fork choice**. ChainCoin nodes all see both branches, evaluate both identically, and pick the same one — no coordination needed. The answer: **the chain with the most cumulative proof-of-work, period.** Difficulty × block count, summed.

## 2. Concept corner

- **Most work, not most blocks**: a difficulty-3 block is 16× the work of a difficulty-2 block (each hex zero multiplies work by 16). Work sums: `Σ 16^difficulty` per block. That is why nobody can out-fork a long honest chain with cheap low-difficulty forks.
- **Eligibility first, best second**: a chain is eligible only if (a) it has blocks, (b) `chain.Validate()` passes (hash/index/link/genesis), (c) every block meets its **declared difficulty**. Declared difficulty that nothing could produce (e.g. 65 leading zeros — geometrically impossible for a SHA-256 hex string) is a liar and is discarded.
- **Determinism or it doesn't work**: ties are broken by comparing tip hashes — every honest node picks the same chain, because the comparison is a pure byte-order decision (every honest node sees the same bytes).
- **Why hash order?** Real networks use "first seen wins" — order of observation. A deterministic tie-break here plays the same role without a clock: the comparison is a pure byte-order decision (every honest node sees the same bytes).
- **Cumulative work is monotone**: extending a chain only ever increases its score. Re-opening history means doing *more* work than the total spent on the honest chain — which is what "finality with probability" means in Raft-less land too.

## 3. Requirements

| Function | Behaviour |
|---|---|
| `ForkChoice(chains []*chain.Blockchain) (*chain.Blockchain, error)` | `len == 0` → `ErrNoCandidates`. Discard every chain that: is nil; has no blocks; fails `chain.Validate()`; or has a block whose `Hash` fails `pow.Verify(b, bc.Difficulty)`. Among the **eligible**, maximise `Σ 16^difficulty`. Ties → lower tip hash wins (hex string comparison). None eligible → `ErrNoValidChain`. |

## 4. What the tests pin down

- Longest valid chain wins between competitors.
- A longer *forged* chain (corrupt data) loses to a shorter honest one.
- Greater work wins even with fewer blocks (difficulty 3 beats difficulty 2 with more blocks).
- Equal work ties resolve deterministically to the lower tip hash.
- Zero offers → `ErrNoCandidates`; every offer invalid → `ErrNoValidChain`; nil candidates are skipped, not fatal.
- A chain declaring an impossible difficulty (65) is discarded even though its structure validates.

## 5. Hints (open one at a time)

<details><summary>Hint 1 — eligibility is a function</summary>

`validChain(bc)`: nil → false; `len == 0` → false; `chain.Validate()`; every block `pow.Verify(b, bc.Difficulty)`.

</details>

<details><summary>Hint 2 — work is a sum</summary>

`per := 1 << (difficulty*4)` (saturate at max uint64 if ≥ 64 bits); loop blocks, add, saturate on overflow.

</details>

<details><summary>Hint 3 — the tie-break</summary>

`tipLess(a, b)`: compare `a.Blocks[len(a.Blocks)-1].Hash` vs `b.Blocks[len(b.Blocks)-1].Hash` — hex bytes compare exactly like big-endian numbers.

</details>

<details><summary>Hint 4 — accumulate carefully</summary>

Track `best` (pointer) and `bestWork` (uint64). Only initialize on the first eligible chain; on equal work, compare tips.

</details>

## 6. Stretch challenges

1. **Most-work vs longest:** build a fork pair where the longest chain has *less* cumulative work — watch the tie-break fail and learn why Bitcoin ≥ 2017 picked "greatest work" precisely because of this.
2. **Difficulty retargeting:** adjust difficulty every 10 blocks so blocks average 6 seconds — a closed control loop (measure, compare, retarget).
3. **Orphan handling:** when the loser is discarded, what becomes of transactions in its orphaned blocks? Implement mempool re-queue.
4. **Fuzz forks:** generate k random valid forks, pick one, and verify every node picks the same one.

## 7. Reflect

- Why must every honest node reach the same fork independently?
- Why is hash order a workable tie-break but block height alone isn't enough?
- Why must "impossible" difficulties be filtered before work is even compared?
- Why does most-work beat longest-valid only as long as work is monotone?

## 8. Done when

`go test ./internal/fork -run Step07 -race -v` → **ok** (apply steps 01–07 in order: `for s in 01 02 03 04 05 06 07; do bash scripts/apply-solution.sh $s; done`). Then `go vet ./...`, `gofmt -l .`.

## 9. The end of ChainCoin — and beyond

You now have: hashing, proof of work, a validating chain, signatures, a mempool, a UTXO ledger, atomic verified persistence, and fork choice. That is — genuinely — the architectural skeleton of a cryptocurrency.

Stretch ideas if you're still hungry: difficulty retargeting; orphan/re-queue; HD wallets; a `json-rpc`? next node (`localhost:8545`-ish); a 64-block demo network with three simulated peers.
