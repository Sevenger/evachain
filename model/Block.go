package model

import (
	"github.com/sevenger/evachain"
	"strconv"
	"time"
)

type Data interface{}

type Block struct {
	Index       int    `comment:"下标"`
	TimeStamp   string `comment:"时间戳"`
	Data        Data   `comment:"数据"`
	Hash        string `comment:"哈希"`
	PreviewHash string `comment:"前一个哈希"`
}

func GenerateNextBlock(lastBlock Block, blockData Data) *Block {
	index := lastBlock.Index + 1
	timeStamp := strconv.Itoa(int(time.Now().Unix()))
	hash := evachain.CalculateHash(index, lastBlock.Hash, blockData)
	return NewBlock(index, timeStamp, blockData, hash, lastBlock.Hash)
}

func NewBlock(index int, timeStamp string, data Data, hash string, previewHash string) *Block {
	return &Block{
		Index:       index,
		TimeStamp:   timeStamp,
		Data:        data,
		Hash:        hash,
		PreviewHash: previewHash,
	}
}
