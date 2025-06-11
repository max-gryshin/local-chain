package service

import (
	"crypto/sha512"
	"strconv"
	"time"

	"local-chain/internal"

	"local-chain/internal/types"
)

type BlockHeader struct {
	Timestamp  int64
	PrevHash   []byte
	Hash       []byte
	MerkleRoot []byte
}
type Block struct {
	Header       BlockHeader
	Transactions []*types.Transaction
}

// Blockchain represents a private blockchain.
type Blockchain struct {
	Blocks []*Block
}

// NewBlockchain creates a new blockchain with a genesis block.
func NewBlockchain() *Blockchain {
	genesisBlock := &Block{
		Header: BlockHeader{
			Timestamp: time.Now().UnixNano(),
			PrevHash:  []byte{},
			Hash:      []byte{},
		},

		Transactions: []*types.Transaction{},
	}
	return &Blockchain{
		Blocks: []*Block{genesisBlock},
	}
}

// Создание нового блока - это добавление транзакций в блокчейн тобиш способ их фиксации так как блок содержит
// ссылку на предыдущий блок и хеш корня дерева Меркла, который позволяет проверить целостность транзакций в блоке.
// Тоесть чтобы проверить присутствует ли транзакция в блоке надо построить дерево меркла блока потом верифицировать транакцию

// когда пользователь создает транзакцию он отправляет ее в пул транзакций (пул по сути блок который еще незаписан в блокчейн)
// по истечении какого то времени блокчейн берет пул ничаниает процесс записи блока в цепочку (консенсус голосование вся фигня пока опускаем)

// возможно стоит выделить отдельный тип как пул транзакций, который будет хранить те транзакции которые еще не включены в блокчейн
// пул будет гулять по сети и ее участники будуть его обновлять, добавляя туда свои транзакции
// и удаляя те которые уже включены в блокчейн, а также те которые не прошли (верификацию???) (и не прошли консенсус?)
// Таким образом для добавления блока в AddBlock нужно передовать тип пул транзакций на основе которого будет строиться блок.
// и таким образом addBlock отвечает уже за создание и добавление нового блока в блокчейн (и соответственно лог транзакцй)

// AddBlock adds a new block to the blockchain.
func (bc *Blockchain) AddBlock(pool *types.Pool) error {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]

	txs := make([]*types.Transaction, 0, len(pool.Transactions))
	for _, txPool := range pool.Transactions {
		txs = append(txs, txPool.Tx)
	}
	merkleTree, err := internal.NewMerkleTree(txs)
	if err != nil {
		return err
	}

	newBlock := &Block{
		Header: BlockHeader{
			Timestamp:  time.Now().UnixNano(),
			PrevHash:   prevBlock.computeHash(),
			MerkleRoot: merkleTree.Root.Hash,
		},
		Transactions: txs,
	}
	bc.Blocks = append(bc.Blocks, newBlock)
	return nil
}

// computeHash computes the hash of a block.
func (b *Block) computeHash() []byte {
	hash := sha512.New()
	hash.Write(append([]byte(strconv.FormatInt(b.Header.Timestamp, 10)), b.Header.PrevHash...))
	return hash.Sum(nil)
}
