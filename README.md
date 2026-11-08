# Blockchain-go — learn blockchain by building a cryptocurrency from first principles, step by step

A hands-on course in **blockchain**, built the same way as [shopflow-go](../shopflow-go), [dataflow-go](../dataflow-go)
and [distkv-go](../distkv-go): **one step at a time, each with a PRD, failing tests, stubs, collapsible hints,
stretch challenges and a reference solution**.

You are building **ChainCoin** — a small proof-of-work cryptocurrency — out of first principles, in Go. Step 01
gives you the atom (the block: hashing and linking). Steps 02–08 add **proof of work**, the **chain**,
**signed transactions**, a **UTXO set**, **verified persistence**, **fork choice** and **Merkle trees**.

The goal: every day you add one small, testable piece of the system, and you *understand* why it is there —
because you wrote it, tested it red-first, and debugged it — instead of only reading about it.

## What you are building, and why this shape

A cryptocurrency is the smallest possible "real" blockchain system, and it contains every lesson in order:

1. **Hashing** — one unforgeable fingerprint per block. (Step 01.)
2. **Proof of work** — why "zeroes after n tries" is consensus. (Step 02.)
3. **A chain** — genesis, validation, forks. (Step 03.)
4. **Transactions** — transfers, signed with ECDSA keys (real signatures, not toy ones). (Step 04.)
5. **UTXO** — balances, coinbase rewards, "who owns what" from the ledger alone. (Step 05.)
6. **Persistence** — save the chain, reload it, verify it. (Step 06.)
7. **Fork choice** — most work wins; every honest node picks the same chain. (Step 07.)
8. **Merkle trees** — one root commits to every transaction. (Step 08.)

By the end you have something every cryptocurrency has in common, at a size that fits in a README —
**so you know where the real ones get correct and where they just make trade-offs**.

## Start here

```bash
cd blockchain-go
go version                      # 1.22 or newer
go test ./...                   # everything fails CLEANLY at first: no panics, no hangs. That is the starting state.
go test ./internal/chain -run Step01 -v
```

You should see **67 failing tests** across the 8 steps. That is the "red" in red-green-refactor. 🎯

Work through a step:

1. Read the step's PRD (`docs/steps/step-01-the-block.md` — the current one; siblings are in `docs/steps/`) — especially
   the **Concept corner**.
2. Open the tests (`internal/<package>/<step>_test.go`) and read them **before** coding. They are the spec.
3. Open the stub (the doc comment on each stub repeats the requirements).
4. Run only that step while you work: `go test ./internal/pow -run Step02 -v`.
5. When green, read **Reflect**, try a **Stretch**, run `go vet ./...` and `gofmt -l .`, then compare with
   `_solutions/stepNN/…`. The differences are where the learning is.
6. `git commit`, next step.

**Stuck?** Hints in each PRD go from gentle to explicit. As a last resort:

```bash
bash scripts/apply-solution.sh 01            # macOS / Linux / Git-Bash
.\scripts\apply-solution.ps1 01              # Windows PowerShell
```

> Note: higher steps build on earlier ones. To test step 07, apply 01–03 first (`for s in 01 02 03; do bash scripts/apply-solution.sh $s; done`). Step 06 needs 01–03 and 05.

## Roadmap

| # | Package | Feature | Go ideas | Patterns / principles | Tests |
|---|---------|---------|----------|-----------------------|:------:|
| 01 | `chain` | [The block: SHA-256 hashing and linking](docs/steps/step-01-the-block.md) | `crypto/sha256`, `encoding/hex`, canonical encoding | hashing, integrity | 6 |
| 02 | `pow` | [Proof of Work: difficulty, mining, verifying](docs/steps/step-02-proof-of-work.md) | `strings`, `time`, tight loops | mining, consensus | 7 |
| 03 | `chain` | [The blockchain: genesis, validation, forgery detection](docs/steps/step-03-the-chain.md) | sentinels, `%w`, traversal | invariants, validation | 10 |
| 04 | `tx` | [Transactions: ECDSA keys, PKIX addresses, signatures, mempool](docs/steps/step-04-transactions.md) | `crypto/ecdsa`, `crypto/x509`, `crypto/rand` | signatures, mempool | 9 |
| 05 | `utxo` | [The UTXO set: balances, coinbase rewards, atomic transfers](docs/steps/step-05-utxo.md) | map models, deterministic ordering | UTXO, ledger | 10 |
| 06 | `persist` | [Persistence: atomic save, verified load](docs/steps/step-06-persistence.md) | `encoding/json`, `fsync`, rename | atomic writes, fail-closed | 7 |
| 07 | `fork` | [Fork choice: most work wins, deterministic tie-break](docs/steps/step-07-fork-choice.md) | cumulative work, determinism, sentinel errors | Nakamoto consensus | 8 |
| 08 | `merkle` | [Merkle trees: one hash for every leaf](docs/steps/step-08-merkle-tree.md) | `crypto/sha256`, trees, proofs | Merkle trees | 10 |

**67 tests in total.** Steps 1–2 stand alone. Step 3 builds on 1–2; step 4 is standalone-ish (uses Go's ECDSA only); step 5 is independent (UTXO); step 6 builds on 1–3 and 5; step 7 builds on 1–3; step 8 stands alone.

## Project layout

```
blockchain-go/
├── README.md
├── go.mod
├── docs/
│   └── steps/             # ← the PRDs, one per step
├── internal/
│   ├── todo/              # ErrNotImplemented for stubs
│   ├── chain/             # step 01 (block), step 03 (chain validation)
│   ├── pow/               # step 02: proof of work
│   ├── tx/                # step 04: transactions, signatures, mempool
│   ├── utxo/              # step 05: unspent outputs, balances
│   ├── persist/           # step 06: save, reload, verify
│   ├── fork/              # step 07: fork choice (most work wins)
│   └── merkle/            # step 08: Merkle trees (roots, proofs)
├── _solutions/            # reference solutions (ignored by the Go tool)
│   └── step01/ … step08/
└── scripts/               # apply-solution.sh / apply-solution.ps1
```

**Why `_solutions`?** The Go tool ignores directories whose names start with `_`, so `go test ./...` and `go vet ./...` never try to compile the reference solutions as if they were packages.

## Ground rules that make this effective

- Don't edit the tests to make them pass. If a test looks wrong, that is a great question — read the PRD again, then ask.
- Resist peeking at `_solutions/` until your tests are green.
- Run `go vet ./...` and `gofmt -l .` before every commit.
- Keep the **Reflect** answers in a notes file. Explaining *why* in your own words is what makes it stick.

## How this project was built

- Starting state: all **67 step tests fail cleanly** — no panics, no hangs, no compile errors in the stub packages.
- Solved state (verified with `_solutions/step01`…`_solutions/step08` applied cumulatively): **all 67 pass under `-race`** with `go vet` and `gofmt` clean.
- The reference solutions are in `_solutions/step01/` through `_solutions/step08/`.
