# Step 06 — Persistence: save it, reload it, trust it

| | |
|---|---|
| **Build** | `internal/persist/persist.go` — `SaveChain`, `LoadChain` |
| **Given** | `chain.Blockchain`, `chain.Validate`, `pow.Verify`, `chain.GenesisPrevHash` (Steps 01–03) |
| **Tests** | `internal/persist/persist_test.go` — 7 tests |
| **Run** | `go test ./internal/persist -run Step06 -race -v` |
| **You'll learn** | atomic writes (temp+rename), `fsync` discipline, "parse ≠ trust", verify-on-load |
| **Time** | ~60 min |

## 1. The story

A blockchain you cannot reload is a blockchain you cannot run. Step 06 makes ChainCoin durable — and, more importantly, **trustworthy on the way back in**: a JSON file that claims to be "a chain" is indistinguishable from a malicious peer until it passes every check the chain itself uses.

Two rules make this step feel like engineering rather than plumbing:

1. **Atomic writes** — the file is either complete or absent, never half-written.
2. **No trust without verification** — every load re-checks everything.

## 2. Concept corner

- **Write-temp then rename**: write to `path.tmp` → `fsync` → `os.Rename`. Either the whole chain appears, or nothing does; a crash mid-save never leaves a torn half-file. (Renames are atomic on POSIX.)
- **`json.Marshal` of the chain**: the entire `Blockchain{Difficulty, Blocks}` round-trips textually. Human-inspectable is a feature: you can `cat chain.json` and read the ledger.
- **Parse ≠ trust**: `json.Unmarshal` only tells you the file is *syntactically* valid. Semantically it could be a chain of zeros that "works" as JSON but is nonsense — so `LoadChain` re-runs `Validate` (hash, index, links, genesis) and `pow.Verify` (proof of work, difficulty) before returning.
- **Fail closed**: any verification failure returns an error and NOTHING is loaded. A node that loads a bad chain is strictly worse than a node that refuses to start.

## 3. Requirements

| Function | Behaviour |
|---|---|
| `SaveChain(path string, bc *chain.Blockchain) error` | Serialise the chain to JSON; write to `path + ".tmp"`; `fsync`; `os.Rename` to `path`. Errors wrapped with `%w`. |
| `LoadChain(path string) (*chain.Blockchain, error)` | Read + parse the file (missing file → error wrapping `os.ErrNotExist`). Then **fully verify**: every block hash matches; indices contiguous; genesis links to sixty-four zeros; every other block links to its predecessor (via `chain.Validate()`); every block meets the stored difficulty (via `pow.Verify`). Any failure → an error, nothing returned. |

## 4. What the tests pin down

- Save → load round-trips: difficulty, block count, every block field (exact equality).
- A missing file is an error wrapping `os.ErrNotExist`.
- A middle-byte-flipped file never loads.
- A chain whose stored difficulty is impossible (e.g. 64) never loads.
- Save leaves no `.tmp` behind; the file parses as JSON with the right difficulty/block count.
- Two consecutive loads produce identical chains.
- Saving to a missing directory errors with `os.ErrNotExist`.

## 5. Hints (open one at a time)

<details><summary>Hint 1 — SaveChain</summary>

`json.Marshal` a small private struct with `difficulty` and `blocks`; open `path+".tmp"` with `O_CREATE|O_WRONLY|O_TRUNC`; write; `Sync()`; close; `os.Rename(tmp, path)`. Wrap each step with `%w`.

</details>

<details><summary>Hint 2 — LoadChain: the order</summary>

Read file (missing → `os.ErrNotExist` wrapped) → `json.Unmarshal` → build `Blockchain` → `bc.Validate()` → each block `pow.Verify(b, bc.Difficulty)` → return `bc`.

</details>

<details><summary>Hint 3 — Validate covers the genesis rule</summary>

`chain.Validate()` already demands hash-matches, contiguity, and the genesis link — reuse it; don't re-implement it.

</details>

<details><summary>Hint 4 — The impossible difficulty test</summary>

Editing `"difficulty":2` → `"difficulty":64` makes every block fail `pow.Verify` — deterministic, no luck involved.

</details>

## 6. Stretch challenges

1. **fsync the directory** after `Rename` — what is lost without it?
2. **Checkpoint format:** serialise only `{difficulty, tip}` + per-block records; reload in O(blocks) with O(1) allocation. Measure.
3. **Compression:** `gzip` the JSON; measure the size win on a 1000-block chain.
4. **Versioning:** add `"version": 1` to the file; LoadChain rejects unknown versions. When would you bump it?

## 7. Reflect

- What exactly can be lost if you `Write` without `fsync`?
- Why does `Rename` (not `Write`) make the save atomic?
- Why is `json.Unmarshal` never enough? Construct a concrete example of a syntactically valid, nonsense-but-expensive "chain".
- What does fail-closed cost you, and what does it save you?

## 8. Done when

`go test ./internal/persist -run Step06 -race -v` → **ok**. Then `go vet ./...`, `gofmt -l .`, and compare with `_solutions/step06/…`.

## 9. The end — and beyond

Congratulations: you have built ChainCoin — hashing, proof of work, a validating chain, signatures, a mempool, a UTXO ledger, and atomic, verified persistence.

Stretch roadmap if you're still hungry:

- **Forks and fork choice:** hand two nodes different genesis data, let them both mine, resolve to the longer chain.
- **Peer sync:** serve `chain.json` over HTTP; verify a received chain before adopting it (you already wrote the verifier).
- **Merkle roots over transactions:** block `Data` = merkle tree of tx IDs; `Validate` rejects forged tx lists.
- **Difficulty adjustment:** every 2016 blocks, retarget so blocks average 10 minutes.
- **The README moment:** a `demo.sh` that creates Alice and Bob, coins Alice via coinbase, has her pay Bob, persists the chain, and re-verifies it.
