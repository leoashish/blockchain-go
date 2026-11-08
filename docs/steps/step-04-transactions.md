# Step 04 — Transactions: ECDSA signatures and the mempool

| | |
|---|---|
| **Build** | `internal/tx/ops.go` — `GenerateKeyPair`, `Address`, `SigningBytes`, `Sign`, `Verify`, `Mempool.{Add, Pending, Drain}` |
| **Given** | `Transaction`, `KeyPair`, `Mempool`, `ErrBadAmount`, `ErrBadSignature`, `NewMempool` (`internal/tx/tx.go`) |
| **Tests** | `internal/tx/tx_test.go` — 10 tests |
| **Run** | `go test ./internal/tx -run Step04 -race -v` |
| **You'll learn** | ECDSA (P-256), PKIX public-key serialisation, canonical signing bytes, signature verification, mempools |
| **Time** | ~2 h |

## 1. The story

ChainCoin has no users yet — blocks are data, but data needs **owners**. A transaction says "I, Alice, move 50 coins to Bob" — and proves it with a signature only Alice's key can produce. The mempool is where those signed promises wait to be mined into a block.

This is the step where ChainCoin becomes a real currency: balances become enforceable, and forgery stops being an assertion and becomes a hard mathematical failure.

## 2. Concept corner

- **ECDSA P-256** (NIST curve): key pairs are cheap, signatures are compact, and verification is fast. `crypto/ecdsa` in the standard library does all of the math.
- **PKIX serialisation**: `x509.MarshalPKIXPublicKey(pub)` gives a stable DER encoding of the public key; our "address" is its hex encoding — the universal, no-ambiguity way to print a key. (Real coins hash it; we skip that for clarity — the lesson is unchanged.)
- **Sign everything except the signature**: `SigningBytes` covers `From`, `To`, `Amount`, `Timestamp` — never `Signature`. Otherwise verifying your own signature would require the signature already being known… an infinite regress.
- **Length-prefix your sign bytes**: same discipline as Step 01 — `From|To|…` concatenation is ambiguous; `uvarint(len) ‖ field` is not.
- **r and s, 64 bytes**: ECDSA gives you two 256-bit integers; we pack them into a fixed 64-byte signature. (DER would also work; raw is simpler.)

## 3. Requirements

| Function | Behaviour |
|---|---|
| `GenerateKeyPair() (KeyPair, error)` | Fresh P-256 pair via `crypto/rand`. |
| `Address(pub *ecdsa.PublicKey) string` | Hex-encoded PKIX encoding of the public key. Deterministic. |
| `SigningBytes(t Transaction) []byte` | SHA-256 of length-prefixed `From‖To‖Amount‖Timestamp`. Ignores `Signature`. |
| `Sign(t Transaction, priv *ecdsa.PrivateKey) (Transaction, error)` | ECDSA-signs `SigningBytes(t)`; sets `Signature` to 64 bytes (r ‖ s). |
| `Verify(t Transaction) bool` | False on any garbage: malformed hex in `From`, malformed PKIX, wrong length signature, or bad signature. No panics. |
| `Mempool.Add(t Transaction) error` | Queues `t` iff `Amount > 0` (`ErrBadAmount`) and `Verify(t)` (`ErrBadSignature`). |
| `Mempool.Pending() []Transaction` | A **copy**, oldest first. |
| `Mempool.Drain() []Transaction` | All of them, and the pool is emptied (as a miner does once they are in a block). |

## 4. What the tests pin down

- Generated key pairs are usable (public matches private).
- Addresses are deterministic and lowercase hex.
- `SigningBytes` is deterministic, covers every non-signature field, ignores `Signature`.
- Sign → Verify round-trips; tampering with any field breaks verification; garbage (bad hex, short signature, empty) never verifies and never panics.
- The mempool accepts valid transactions, rejects zero-amount and forged ones (with the right sentinels), drains correctly, and hands out copies — not its backing slice.

## 5. Hints (open one at a time)

<details><summary>Hint 1 — GenerateKeyPair</summary>

`priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)`

</details>

<details><summary>Hint 2 — Address</summary>

`der, _ := x509.MarshalPKIXPublicKey(pub); return hex.EncodeToString(der)` (unreachable error — a valid key always marshals).

</details>

<details><summary>Hint 3 — SigningBytes</summary>

Iterate `[]string{t.From, t.To, strconv.FormatInt(t.Amount, 10), strconv.FormatInt(t.Timestamp, 10)}`, length-prefix each into `sha256.New()` via `binary.PutUvarint`, return `h.Sum(nil)`.

</details>

<details><summary>Hint 4 — Sign</summary>

`r, s, err := ecdsa.Sign(rand.Reader, priv, digest)`, then write `r` and `s` big-endian into a zeroed 64-byte slice with sign-byte offset copies.

</details>

<details><summary>Hint 5 — Verify</summary>

`hex.DecodeString(t.From)` → `x509.ParsePKIXPublicKey` → assert `(*ecdsa.PublicKey)` → `ecdsa.Verify(pub, SigningBytes(t), r, s)`. Every failure mode returns false.

</details>

## 6. Stretch challenges

1. **Compressed addresses:** hash the PKIX bytes (SHA-256) and take the hex of the first 16 bytes — 32 chars instead of ~180. Now addresses look like real addresses.
2. **Replay protection:** coin indexes are nonces too — for this toy, `Timestamp` serves as a per-tx nonce. Explain why replay is impossible *within* this design, and where it breaks.
3. **Multi-input transactions:** allow `Inputs []Input` each with its own signature. (Step 05's UTXO model is exactly where this pays off.)
4. **Benchmarks:** time GenerateKeyPair / Sign / Verify. How many tx can this node verify per second?

## 7. Reflect

- Why must the signature cover `Timestamp`?
- Why is `From` not sufficient identity for a transaction's signer without also binding `To`?
- What can an attacker do with a stolen private key, and why can't they do the same with a stolen address?
- Why does the mempool reject zero-amount transactions at ingest?

## 8. Done when

`go test ./internal/tx -run Step04 -race -v` → **ok**. Then `go vet ./...` and `gofmt -l .`, and compare with `_solutions/step04/…`.

## 9. Next

Step 05 — **the UTXO set**: balances from the ledger alone, coinbase rewards, change outputs, and why spending is deletion.
