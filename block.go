package main

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"strings"
	"time"
)

var GenesisBlock = &Block{
	Index:        0,
	TimeStamp:    1630480460,
	PreviousHash: "",
	Hash:         "280074b880e45acfced3772b80d83a08fd90e0affd69e15dcc583c66e40863d4",
	Data:         "Hello, EvaChain!",
}

var EvaChain = Chain{GenesisBlock}

type Chain []*Block

type Block struct {
	Index        int64  `json:"index"`
	TimeStamp    int64  `json:"time_stamp"`
	PreviousHash string `json:"previous_hash"`
	Hash         string `json:"hash"`
	Data         string `json:"data"`
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

func GenerateBlock(lastBlock *Block, data string) *Block {
	newBlock := NewBlock(lastBlock.Index+1, time.Now().Unix(), lastBlock.PreviousHash, "", data)
	for i := 0; ; i++ {

	}
	return newBlock
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

func IsValidHash(hash string, difficulty int) bool {
	prefix := strings.Repeat("0", difficulty)
	return strings.HasPrefix(hash, prefix)
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

func IsValidChain(chain Chain) bool {
	if *chain[0] != *GenesisBlock {
		logMsg("is not the same chain")
		return false
	}

	for i := 0; i < len(chain)-1; i++ {
		if !IsValidBlock(chain[i+1], chain[i]) {
			logMsg("contains invalid block:", chain[i])
			return false
		}
	}
	return true
}

func GetLatestBlock() *Block {
	return EvaChain[len(EvaChain)-1]
}

func AddBlock(block *Block) {
	if IsValidBlock(block, GetLatestBlock()) {
		EvaChain = append(EvaChain, block)
	} else {
		logMsg("invalid block:", block)
	}
}

func ReplaceChain(chain Chain) {
	if IsValidChain(chain) && len(chain) > len(EvaChain) {
		logMsg("replace the chain")
		EvaChain = chain
		BoardCast(ResponseLatestMsg())
	}
}
