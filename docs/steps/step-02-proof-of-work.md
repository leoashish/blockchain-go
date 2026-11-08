# Step 02 — Proof of Work: difficulty, mining, verifying

| | |
|---|---|
| **Build** | `internal/pow/pow.go` — `MeetsDifficulty`, `Mine`, `Verify` |
| **Given** | `chain.CalculateHash`, `chain.Block`, `chain.VerifyHash` (from Step 01) |
| **Tests** | `internal/pow/pow_test.go` — 7 tests |
| **Run** | `go test ./internal/pow -run Step02 -v` |
| **You'll learn** | proof of work, difficulty targets, why mining is expensive but verifying is free |
| **Time** | ~45 min |

## 1. The story

Blocks (Step 01) prove nothing yet: anyone can compute a valid hash. Proof of work makes blocks **expensive** and verifying **free** — that asymmetry is the whole consensus mechanism:

- **Mining**: keep bumping the nonce until the hash starts with `N` zeros. Pure brute-force search; expected tries = 16^N hex digits.
- **Verifying**: recompute one hash and check the zeros. Milliseconds.

Bitcoin uses 20 leading zero *bits* (~2^20 ≈ 1,000,000 expected tries). We use leading zero *hex characters*; `N=2` means ~256 tries — perfect for a learner, and `N=64` would be the heat-death-of-universe tier.

## 2. Concept corner

- **The honest asymmetry**: mining is O(16^N) hashes; verifying is O(1). This is why "the longest valid chain" is also the most-expensive chain to build — and why re-writing history costs more than writing it the first time.
- **Hex-prefix difficulty (not bits)**: `strings.HasPrefix(hash, strings.Repeat("0", difficulty))` — simple, readable, and perfectly adequate for teaching. (Real networks use integer comparisons against a big target; same idea, better encoding.)
- **Determinism by accident**: we iterate the nonce *ascending from its current value*, so mining the same block twice gives the same result. Real miners randomise nonce selection instead.
- **Verify means everything**: `chain.VerifyHash(b) && MeetsDifficulty(b.Hash, difficulty)` — a tampered block fails the hash; an intact-but-unlucky block fails the difficulty.

## 3. Requirements

| Function | Behaviour |
|---|---|
| `MeetsDifficulty(hash string, difficulty int) bool` | True iff `hash` starts with at least `difficulty` zero hex characters. `difficulty <= 0` → always true. |
| `Mine(b chain.Block, difficulty int) chain.Block` | Increments `b.Nonce` (from its current value) until `chain.CalculateHash(b)` meets the difficulty; sets `b.Hash`; returns the mined block. Identity fields unchanged. |
| `Verify(b chain.Block, difficulty int) bool` | True iff `chain.VerifyHash(b)` AND `MeetsDifficulty(b.Hash, difficulty)`. |

## 4. What the tests pin down

- `MeetsDifficulty`: leading zeros only (middle zeros don't count), exact-length matches, zero/negative difficulty passes anything, empty hash fails.
- `Mine` with difficulty 2 produces a hash starting with `"00"`; preserves Index/Timestamp/Data/PrevHash; sets Hash; is deterministic for the same inputs.
- `Verify`: accepts a mined block; rejects an intact block without enough leading zeros; rejects a tampered block (hash no longer matches).

## 5. Hints (open one at a time)

<details><summary>Hint 1 — MeetsDifficulty</summary>

Guard `difficulty <= 0` → true. Guard `len(hash) < difficulty` → false. Then `strings.HasPrefix(hash, strings.Repeat("0", difficulty))`.

</details>

<details><summary>Hint 2 — Mine</summary>

`for { b.Hash = chain.CalculateHash(b); if MeetsDifficulty(b.Hash, difficulty) { return b }; b.Nonce++ }`

</details>

<details><summary>Hint 3 — Verify</summary>

`return chain.VerifyHash(b) && MeetsDifficulty(b.Hash, difficulty)` — the elegance of proof of work is exactly this expressibility.

</details>

## 6. Stretch challenges

1. Time the miner: `N` digits for N=1..4. Does your O(16^N) match theory (~16^N expected tries)?
2. **Difficulty adjustment:** add a function that picks difficulty so that blocks take ~10 seconds (measure hashes/sec, solve for N).
3. **Bits vs hex:** implement `MeetsDifficultyBits(hashHex string, bits int) bool` operating on bits and compare against the hex version.
4. Change mining to start from a random offset (`crypto/rand`) — mining becomes non-deterministic. Why do real miners do this?

## 7. Reflect

- Why does an attacker need to redo *all* the work to rewrite history, but an honest verifier never has to?
- What changes (and what stays) if difficulty is 1 hex zero vs 8?
- Why is proof of work called "work" and not just "ordering"?
- Why can a difficulty-2 block *coincidentally* satisfy difficulty 4? (And why is that fine — even good?)

## 8. Done when

`go test ./internal/pow -run Step02 -race -v` → **ok**. Then `go vet ./...` and `gofmt -l .`, and compare with `_solutions/step02/…`.

## 9. Next

Step 03 — **the blockchain**: genesis, contiguous indices, link validation, and catching every way to tamper.
