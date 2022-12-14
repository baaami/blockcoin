package blockchain

import (
	"errors"
	"strings"
	"time"

	"github.com/baaami/blockcoin/db"
	"github.com/baaami/blockcoin/utils"
)

type Block struct {
	Hash         string `json:"hash"`
	PrevHash     string `json:"prevHash,omitempty"`
	Height       int    `json:"height"`
	Difficulty   int    `json:"difficulty"`
	Nonce        int    `json:"nonce"`
	Timestamp    int    `json:"timestamp"`
	Transactions []*Tx  `json:"transactions"`
}

func persistBlock(b *Block) {
	db.SaveBlock(b.Hash, utils.ToBytes(b))
}

func (b *Block) restore(data []byte) {
	utils.FromBytes(b, data)
}

var ErrNotFound = errors.New("Block not found ")

func FindBlock(hash string) (*Block, error) {
	blockBytes := db.Block(hash)
	if blockBytes == nil {
		return nil, ErrNotFound
	}

	block := &Block{}
	block.restore(blockBytes)
	return block, nil
}

func (b *Block) mine() {
	// insert "00"
	target := strings.Repeat("0", b.Difficulty)
	for {
		hash := utils.Hash(b)
		if strings.HasPrefix(hash, target) {
			b.Timestamp = int(time.Now().Unix())
			b.Hash = hash
			break
		} else {
			b.Nonce++
		}
	}
}

func createBlock(prevHash string, height, diff int) *Block {
	block := Block{
		Hash:       "",
		PrevHash:   prevHash,
		Height:     height,
		Difficulty: diff,
		Nonce:      0,
	}
	// 작업 증명
	block.mine()
	// 채굴을 끝내고 해시를 찾고 전부 끝낸 다음에 transaction들을 block에 넣어줌
	block.Transactions = Mempool().TxToConfirm()
	// block을 db에 저장
	persistBlock(&block)
	return &block
}
