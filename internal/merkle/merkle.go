// Package merkle implements Merkle trees: one small hash that stands for an entire list of hashes.
package merkle

import "errors"

// ErrIndexOutOfRange: the leaf index is outside the tree. GIVEN.
var ErrIndexOutOfRange = errors.New("merkle: index out of range")
