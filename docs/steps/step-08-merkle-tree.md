# Step 08 — Merkle trees: one hash to rule a thousand

| | |
|---|---|
| **Build** | `internal/merkle/ops.go` — `Root`, `Prove`, `Verify` |
| **Given** | `ErrIndexOutOfRange` (`internal/merkle/merkle.go`); `SHA-256` via `crypto/sha256` |
| **Tests** | `internal/merkle/merkle_test.go` — 10 tests |
| **Run** | `go test ./internal/merkle -run Step08 -v` |
| **You'll learn** | Merkle trees, inclusion proofs, the odd-node duplication rule, why block headers carry only a "root" and light clients prove membership without the block's full data |
| **Time** | ~90 min |

## 1. The story

Block headers on real chains carry *one* number that commits to every transaction in the block: the **Merkle root**. A block may contain ten thousand transactions; the header is one hash. If any transaction is altered, the root no longer matches the header — and anyone holding the header can prove which transaction was forged.

That is what this step builds: a tree of hashes from leaves to a single root, and a **proof** that lets a light client convince itself that a transaction is in the block — by sending only the path, not the whole tree.

## 2. Concept corner

- **Hash pairs upward**: every leaf is hashed; neighbours are hashed together over the **concatenation of their 32-byte hashes** — never of the raw content (hash concat is ambiguous, hash-of-hashes is the whole point of a tree).
- **Odd nodes are duplicated** — Bitcoin's rule: an odd leaf at any level is paired with *itself*.
- **Proof = the siblings on the path**: ~log2(n) hashes are enough for a light client to verify membership without holding the block.
- **Root of an empty tree is the all-zero 32 bytes** — Bitcoin's convention, also ours.
- **Sideways is the verification order**: `idx % 2 == 0` → you are the left child (concatenate h‖sib); otherwise swap. At every level, the sibling is `index ^ 1`. Follow it bottom-up and you arrive at the root — or not, which is exactly what makes a fake proof fail.
- **Hex everywhere in our API**: hashes are plain strings — the tests teach it; real systems use raw bytes.

## 3. Requirements

| Function | Behaviour |
|---|---|
| `Root(leaves [][]byte) string` | Hash leaves → pair-hash level by level until one node remains → hex of the root. Empty input → sixty-four zeros. |
| `Prove(leaves [][]byte, index int) ([]string, error)` | Sibling hashes bottom-up (hex). Index out of range → `ErrIndexOutOfRange`. |
| `Verify(leaf []byte, index int, proof []string, root string) bool` | Recompute the root from `leaf`, `index` and the sibling hashes (bottom-up, hex). Malformed hex (or wrong length) → false. Never panics; false on any mismatch. |

## 4. What the tests pin down

- Zero leaves → `strings.Repeat("00", 32)`.
- One leaf → `sha256(leaf)`.
- Two leaves → `sha256(sha256(a) ‖ sha256(b))`.
- Any single-byte change anywhere flips the root.
- An odd node at any level is duplicated for pairing.
- `Prove`+`Verify` round-trip for tree sizes 1…17 for **every** index; a proof verifies only for the right `(leaf, index)` pair.
- A tampered leaf, a flipped nibble in any proof element, the wrong index, or the wrong root — all rejected.
- A 1000-leaf tree verifies its first, middle and last member.
- A proof with a missing (or reordered) element fails verification.
- Out-of-range indices → `ErrIndexOutOfRange` (both negative and too large).

## 5. Hints (open one at a time)

<details><summary>Hint 1 — one hashing helper is enough</summary>

`hashPair(a, b []byte) []byte` — `sha256.Sum256(append(append([]byte{}, a...), b...))`. Leaf hashing uses the same pattern without concatenation. Copy inputs before append; never mutate child hashes in place.

</details>

<details><summary>Hint 2 — the tree loop</summary>

`for len(cur) > 1 { if odd { duplicate last }; next := pairs hashed }` — same shape for Root and Prove, so factor `hashLeaves` and `pairUp`.

</details>

<details><summary>Hint 3 — Prove</summary>

Walk down just like Root, but collect `cur[idx^1]` at each level before collapsing it; `idx /= 2` each level. Index bounds checked up front.

</details>

<details><summary>Hint 4 — Verify</summary>

`h := sha256(leaf)`; for each proof element: decode hex (reject on error), `h = idx%2==0 ? sha(h‖sib) : sha(sib‖h)`; `idx/=2`. Return `h == root`.

</details>

## 6. Stretch challenges

1. **Sparse merkle tree**: handle `n` fixed at 212 leaves with empty leaves pre-hashed to 32 zeros (what real light clients do).
2. **Second-preimage test**: prove that a hash-concatenation tree is weaker than a hash-of-hash tree for 2-to-2 substitutions. Put it in Reflect.
3. **Tie it back to Step 04**: make a transaction tree where `Data = Root(tx-ids)`; verification via `Verify(tx.ID(), i, proof, blockData)`.
4. **Streaming Prove**: call Prove for every index on a 10 000-leaf tree — how long does a real node spend per request? (That's why real nodes cache.)

## 7. Reflect

- Why does the header hash represent the whole block, and why can't a forged transaction pass verification with a valid proof?
- Why hash the concat of hashes rather than the concat of raw contents?
- Why duplicate the odd node (rather than carry odd nodes up unchanged)?
- Why do light clients only need `Prove`/`Verify` and never hold the full block body?

## 8. Done when

`go test ./internal/merkle -run Step08 -v` → **ok**. Then `go vet ./...`, `gofmt -l .`, and compare with `_solutions/step08/…`.

## 9. The end — and beyond

ChainCoin now stores, signs, organizes, persists and proves itself. Remaining stretch wants if you're still hungry: **difficulty retargeting** (closed-loop control); **HD wallets**; **fee market (mempool eviction under pressure)**; **json-rpc ?mini-node**.
