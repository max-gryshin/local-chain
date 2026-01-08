package merkle

import (
	"crypto/sha512"
	"errors"

	"local-chain/internal/types"
)

// Node represents a node in the Merkle Tree.
type Node struct {
	left   *Node
	right  *Node
	parent *Node
	Hash   []byte
	tx     *types.Transaction
}

// MerkleTree represents a Merkle Tree.
type MerkleTree struct {
	Root   *Node
	Leaves []*Node
}

// NewMerkleTree creates a new Merkle Tree from a list of transactions.
func NewMerkleTree(txs ...*types.Transaction) (*MerkleTree, error) {
	if len(txs) == 0 {
		return nil, errors.New("no transactions provided")
	}

	leaves := make([]*Node, len(txs))
	for i, tx := range txs {
		leaves[i] = &Node{
			tx:   tx,
			Hash: tx.Hash,
		}
	}

	// build tree from the bottom to the top
	root, err := buildTree(leaves)
	if err != nil {
		return nil, err
	}

	return &MerkleTree{
		Root:   root,
		Leaves: leaves,
	}, nil
}

// buildTree builds the Merkle Tree from a list of leaf nodes.
func buildTree(leaves []*Node) (*Node, error) {
	if len(leaves) == 0 {
		return nil, errors.New("no leaves provided")
	}

	// If just 1 leaf then return a root
	if len(leaves) == 1 {
		return leaves[0], nil
	}

	// create parent nodes
	var parents []*Node
	for i := 0; i < len(leaves); i += 2 {
		parent := &Node{}
		if i+1 < len(leaves) {
			// there is pair of nodes
			parent.left = leaves[i]
			parent.right = leaves[i+1]
			leaves[i].parent = parent
			leaves[i+1].parent = parent
			// computing parent's hash: H(left || right)
			hash := sha512.New()
			hash.Write(append(leaves[i].Hash, leaves[i+1].Hash...))
			parent.Hash = hash.Sum(nil)
		} else {
			// An unpaired number of leaves, duplicating the last one
			parent.left = leaves[i]
			leaves[i].parent = parent
			parent.Hash = leaves[i].Hash
		}
		parents = append(parents, parent)
	}

	// Recursively build the next level
	root, err := buildTree(parents)
	if err != nil {
		return nil, err
	}
	return root, nil
}

// VerifyTransaction verifies if a transaction is in the Merkle Tree.
func (m *MerkleTree) VerifyTransaction(tx *types.Transaction) (bool, error) {
	// looking tx index in leafs
	var leafIndex int
	var found bool
	for i, leaf := range m.Leaves {
		if leaf.tx.ID == tx.ID {
			leafIndex = i
			found = true
			break
		}
	}
	if !found {
		return false, errors.New("transaction not found in a tree")
	}

	// build Merkle Path
	path, err := m.getMerklePath(leafIndex)
	if err != nil {
		return false, err
	}

	return m.verifyPath(tx.Hash[:], leafIndex, path), nil
}

// getMerklePath returns the Merkle Path for a leaf at the given index.
func (m *MerkleTree) getMerklePath(index int) ([][]byte, error) {
	if index < 0 || index >= len(m.Leaves) {
		return nil, errors.New("invalid leaf index")
	}

	var path [][]byte
	current := m.Leaves[index]

	// walk up the tree
	for current.parent != nil {
		parent := current.parent
		if parent.left == current && parent.right != nil {
			// add hash right node in the path
			path = append(path, parent.right.Hash)
		} else if parent.right == current {
			// add hash left node in the path
			path = append(path, parent.left.Hash)
		}
		current = parent
	}

	return path, nil
}

// verifyPath verifies the Merkle Path for a transaction hash.
func (m *MerkleTree) verifyPath(txHash []byte, index int, path [][]byte) bool {
	currentHash := txHash
	for _, siblingHash := range path {
		hash := sha512.New()
		if index%2 == 0 {
			// current node left - add right hash
			hash.Write(append(currentHash, siblingHash...))
		} else {
			// current node right - add left hash
			hash.Write(append(siblingHash, currentHash...))
		}
		currentHash = hash.Sum(nil)
		index /= 2
	}

	// compare computed hash with the root
	return string(currentHash) == string(m.Root.Hash)
}
