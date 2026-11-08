# Step 01 — The block: SHA-256 hashing and linking

| | |
|---|---|
| **Build** | `internal/chain/block.go` — `CalculateHash`, `NewBlock`, `NewGenesisBlock`, `VerifyHash` |
| **Given** | `Block` (the struct: Index, Timestamp, Data, PrevHash, Nonce, Hash), `GenesisPrevHash` (`internal/chain/chain.go`) |
| **Tests** | `internal/chain/chain_test.go` — 6 tests |
| **Run** | `go test ./internal/chain -run Step01 -v` |
| **You'll learn** | `crypto/sha256`, `encoding/hex`, why hash input needs unambiguous encoding, how blocks link by hash, tamper detection |
| **Time** | ~45–60 min |

## 1. The story

**ChainCoin** is a small cryptocurrency built from first principles. It will eventually have a chain of blocks, proof of work, transactions signed with real keys, balances computed from a UTXO set, and peer sync. None of that matters until the **atom** exists: a **block**.

A block is nothing but data + two numbers + a pointer:

- `Index` — position in the chain (0 = genesis).
- `Timestamp` — when it was made (Unix seconds).
- `Data` — whatever data the block carries (later: a list of transactions).
- `PrevHash` — the hash of the previous block. **This is the link that makes it a chain.**
- `Nonce` — unused until Step 02's proof of work; built into the hash already.
- `Hash` — the hash of the block's identity fields (everything except `Hash` itself).

Small, flat, boring — and exactly why the format works: **tamper with one bit anywhere and `Hash` stops matching.**

## 2. Concept corner

- **SHA-256 is a fingerprinting machine.** Feed in any byte string; get an exact 256-bit digest back. Change one bit anywhere → a completely different digest, every time. That is the whole foundation: blocks prove identity with their hash.
- **Hex, not raw bytes.** `[]byte("…")` is unreadable and contains ambiguous bytes; hex (`0-9a-f`) is the canonical, human-inspectable way to print a hash — exactly how you see them on explorers (and in the tests).
- **The hash must exclude `Hash` itself.** `Hash = SHA256(…Hash…)` has no fixed point you can write down. So the hash covers the **identity** fields: `Index`, `Timestamp`, `Data`, `PrevHash`, `Nonce` — everything except `Hash`.
- **Hash input needs unambiguous encoding.** `"0" + "1"` and `"01"` are the same string — but `Index=01, Timestamp=2` must not collide with `Index=0, Timestamp=12`. Concatenating numbers without separators is broken. The clean fix (used by the reference solution) is to **length-prefix** every field before writing it into the hasher: `uvarint(len(field)) || field`, over and over. No two different field lists ever hash to the same bytes.
- **Why link by hash, not by index?** An index can be re-numbered; a hash cannot be faked without producing a collision. Tampering is only visible if the next block's `PrevHash` literally *contains* the tampered hash.
- **Step 01 invents no failure modes.** Hashing and constructing a block cannot fail, so none of these functions returns an error. Validation errors arrive with real validation (Step 03) — the honest way to design signatures.

## 3. Requirements

| Function | Behaviour |
|---|---|
| `CalculateHash(b Block) string` | SHA-256 of the identity fields, returned as 64 lowercase hex characters. Deterministic, `Hash` field ignored. Length-prefix the fields so the input bytes are unambiguous. |
| `NewBlock(index, data, prevHash, timestamp) Block` | Create the block with the given parts (nonce 0) and set its `Hash` to `CalculateHash` of the parts. |
| `NewGenesisBlock(data, timestamp) Block` | `NewBlock(0, data, GenesisPrevHash, timestamp)` — index 0, `PrevHash` = sixty-four zeros. |
| `VerifyHash(b Block) bool` | `b.Hash == CalculateHash(b)` — true iff the block is untampered. |

## 4. What the tests pin down

- `CalculateHash` is deterministic (same block ⇒ same hash) and a valid 64-char lowercase hex string.
- Changing **any** identity field (Index, Timestamp, Data, PrevHash, Nonce) changes the hash — and the `Hash` field itself has no effect.
- `NewBlock` sets `Hash` to `CalculateHash(block)`.
- `NewGenesisBlock` is block 0 with `PrevHash` = sixty-four zeros.
- Three chained blocks link correctly (`b1.PrevHash == g.Hash`, `b2.PrevHash == b1.Hash`).
- `VerifyHash` accepts an intact block and rejects one with a single-bit Data change or a mismatched `Hash`.

## 5. Hints (open one at a time)

<details><summary>Hint 1 — the pieces</summary>

`crypto/sha256`: `h := sha256.New()`, write fields, `h.Sum(nil)`. `encoding/hex`: `hex.EncodeToString(sum)`.

</details>

<details><summary>Hint 2 — what goes into the hash</summary>

Iterate `[]string{ strconv.Itoa(b.Index), strconv.FormatInt(b.Timestamp, 10), b.Data, b.PrevHash, strconv.FormatInt(b.Nonce, 10) }`.

</details>

<details><summary>Hint 3 — unambiguous encoding</summary>

Length-prefix each field: `n := binary.PutUvarint(lenbuf[:], uint64(len(field))); h.Write(lenbuf[:n]); h.Write([]byte(field))`. Now no two different field lists produce the same bytes.

</details>

<details><summary>Hint 4 — VerifyHash is tiny</summary>

`return b.Hash == CalculateHash(b)` — the elegance of "recompute and compare".

</details>

## 6. Stretch challenges

1. **JSON round trip:** add a test that `json.Marshal(block)` → `json.Unmarshal` preserves the block. (Blocks will cross the wire in later steps.)
2. **Second-round question:** why can't two different Data strings produce the same length-prefixed bytes? Prove it to yourself.
3. **Benchmark:** `go test -bench=.` — how many hashes/sec does this implementation do? (Remember it for Step 02's mining targets…)
4. **Zero-heavy hash:** can you, today, make `Hash` start with a `0`? (Answer: luck. Step 02 will weaponise this.)

## 7. Reflect

- Why can't `Hash` cover itself, and who else in a blockchain has the same non-fixed-point problem? (Hint: transactions reference themselves too… be careful.)
- What breaks if `CalculateHash` delimits fields with `|` instead of length-prefixing? Try to find an actual collision on paper.
- Why does block `PrevHash` point at the *previous* block's `Hash` instead of its `Index`?
- Why did this step have **no** errors to return? When will validation errors arrive?

## 8. Done when

`go test ./internal/chain -run Step01 -race -v` → **ok**. Then `go vet ./...` and `gofmt -l .`, and compare with `_solutions/step01/…`.

## 9. Next

Step 02 — **Proof of work**: nonces, difficulty, mining until the hash starts with zeros, verifying.
