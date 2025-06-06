package internal

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
func NewMerkleTree(txs []*types.Transaction) (*MerkleTree, error) {
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

	// Строим дерево снизу вверх
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

	// Если только один лист, возвращаем его как корень
	if len(leaves) == 1 {
		return leaves[0], nil
	}

	// Создаем родительские узлы
	var parents []*Node
	for i := 0; i < len(leaves); i += 2 {
		parent := &Node{}
		if i+1 < len(leaves) {
			// Есть пара узлов
			parent.left = leaves[i]
			parent.right = leaves[i+1]
			leaves[i].parent = parent
			leaves[i+1].parent = parent
			// Вычисляем хэш родителя: H(left || right)
			hash := sha512.New()
			hash.Write(append(leaves[i].Hash, leaves[i+1].Hash...))
			parent.Hash = hash.Sum(nil)
		} else {
			// Непарное количество листьев, дублируем последний
			parent.left = leaves[i]
			leaves[i].parent = parent
			parent.Hash = leaves[i].Hash
		}
		parents = append(parents, parent)
	}

	// Рекурсивно строим следующий уровень
	root, err := buildTree(parents)
	if err != nil {
		return nil, err
	}
	return root, nil
}

// VerifyTransaction verifies if a transaction is in the Merkle Tree.
func (m *MerkleTree) VerifyTransaction(tx *types.Transaction) (bool, error) {
	// Ищем индекс транзакции в листьях
	var leafIndex int
	var found bool
	for i, leaf := range m.Leaves {
		if leaf.tx == tx {
			leafIndex = i
			found = true
			break
		}
	}
	if !found {
		return false, errors.New("transaction not found in tree")
	}

	// Собираем Merkle Path
	path, err := m.getMerklePath(leafIndex)
	if err != nil {
		return false, err
	}

	// Проверяем путь
	return m.verifyPath(tx.Hash[:], leafIndex, path), nil
}

// getMerklePath returns the Merkle Path for a leaf at the given index.
func (m *MerkleTree) getMerklePath(index int) ([][]byte, error) {
	if index < 0 || index >= len(m.Leaves) {
		return nil, errors.New("invalid leaf index")
	}

	var path [][]byte
	current := m.Leaves[index]

	// Идем вверх по дереву
	for current.parent != nil {
		parent := current.parent
		if parent.left == current && parent.right != nil {
			// Добавляем хэш правого узла в путь
			path = append(path, parent.right.Hash)
		} else if parent.right == current {
			// Добавляем хэш левого узла в путь
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
			// Текущий узел — левый, добавляем правый хэш
			hash.Write(append(currentHash, siblingHash...))
		} else {
			// Текущий узел — правый, добавляем левый хэш
			hash.Write(append(siblingHash, currentHash...))
		}
		currentHash = hash.Sum(nil)
		index /= 2
	}

	// Сравниваем вычисленный хэш с корневым
	return string(currentHash) == string(m.Root.Hash)
}
