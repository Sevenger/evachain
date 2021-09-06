package main

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"strings"
	"time"
)

type IBlockChain interface {
	GetGenesisBlock() *Block
	IsValidaChain() bool
	AddBlock(block *Block)
	GetLatestBlock() *Block
	GetBlockHigh() int
	Get() *IBlock
}

type BlockChain struct {
	Blocks []*Block
}

func (bc BlockChain) GetGenesisBlock() *Block {
	return &Block{
		Index:        0,
		TimeStamp:    1630480460,
		PreviousHash: "",
		Hash:         "280074b880e45acfced3772b80d83a08fd90e0affd69e15dcc583c66e40863d4",
		Data:         "Hello, EvaChain!",
	}
}

func (bc BlockChain) IsValidaChain() bool {
	if bc.Blocks[0] != bc.GetGenesisBlock() {
		logMsg("is not the same chain")
		return false
	}

	for i := 0; i < len(bc.Blocks)-1; i++ {
		if !IsValidBlock(bc.Blocks[i+1], bc.Blocks[i]) {
			logMsg("contains invalid block:", bc.Blocks[i])
			return false
		}
	}
	return true
}

func (bc BlockChain) AddBlock(block *Block) {
	panic("implement me")
}

func (bc BlockChain) GetLatestBlock() *Block {
	return bc.Blocks[len(bc.Blocks)-1]
}

func (bc BlockChain) GetBlockHigh() int {
	return len(bc.Blocks)
}

func (bc BlockChain) Get() *IBlock {
	var b interface{} = bc.Blocks[0]
	block := (b).(*IBlock)
	return block
}

type IBlock interface {
	Post()
}

type BlockImpl struct {
	Index        int64  `json:"index"`
	TimeStamp    int64  `json:"time_stamp"`
	PreviousHash string `json:"previous_hash"`
	Hash         string `json:"hash"`
	Data         string `json:"data"`
}

type Block struct {
	Index        int64  `json:"index"`
	TimeStamp    int64  `json:"time_stamp"`
	PreviousHash string `json:"previous_hash"`
	Hash         string `json:"hash"`
	Data         string `json:"data"`
}

func (b Block) Post() {
	panic("implement me")
}

func NewBlock(index, timeStamp int64, previousHash, hash, data string) *Block {
	return &Block{
		Index:        index,
		TimeStamp:    timeStamp,
		PreviousHash: previousHash,
		Hash:         hash,
		Data:         data,
	}
}

func GenerateNextBlock(data string) *Block {
	previousHash := GetLatestBlock().Hash
	index := GetLatestBlock().Index + 1
	timeStamp := time.Now().Unix()
	hash := CalculateHash(index, previousHash, data)
	return NewBlock(index, timeStamp, previousHash, hash, data)
}

func CalculateHash(index int64, previousHash string, data string) string {
	sb := strings.Builder{}
	sb.WriteString(strconv.Itoa(int(index)))
	sb.WriteString(previousHash)
	sb.WriteString(data)
	hash := sha256.Sum256([]byte(sb.String()))
	return hex.EncodeToString(hash[:])
}

func CalculateHashForBlock(block *Block) string {
	return CalculateHash(block.Index, block.PreviousHash, block.Data)
}

func IsValidBlock(newBlock, lastBlock *Block) bool {
	if newBlock.Index != lastBlock.Index+1 {
		logMsg("invalid index")
		return false
	} else if newBlock.PreviousHash != lastBlock.Hash {
		logMsg("invalid previous hash")
		return false
	} else if newBlock.Hash != CalculateHashForBlock(newBlock) {
		logMsg("invalid hash")
		return false
	}
	return true
}

func AddBlock(block *Block) {
	if IsValidBlock(block, GetLatestBlock()) {
		EvaChain = append(EvaChain, block)
	} else {
		logMsg("invalid block:", block)
	}
}

func ReplaceChain(chain IBlockChain) {
	if chain.IsValidaChain() && len(chain) > len(EvaChain) {
		logMsg("replace the chain")
		EvaChain = chain
		BoardCast(ResponseLatestMsg())
	}
}
