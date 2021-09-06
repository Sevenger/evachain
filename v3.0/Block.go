package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type BlockChain struct {
	Blocks []*Block `json:"blocks"`
}

func NewBlockChain() *BlockChain {
	return &BlockChain{
		Blocks: []*Block{{
			High:      0,
			TimeStamp: 1630480460,
			PreHash:   "",
			Hash:      "280074b880e45acfced3772b80d83a08fd90e0affd69e15dcc583c66e40863d4",
			Data:      "Hello, BlockChain!",
		}},
	}
}

func (bc *BlockChain) GetGenesisBlock() IBlock {
	return bc.Blocks[0]
}

func (bc *BlockChain) IsValidChain() bool {
	if bc.Blocks[0] != bc.GetGenesisBlock() {
		return false
	}

	for i := 0; i < len(bc.Blocks)-1; i++ {
		last, next := bc.GetBlock(i), bc.GetBlock(i+1)
		if next.IsValidBlock(last) {
			return false
		}
	}

	return true
}

func (bc *BlockChain) AddBlock(block IBlock) {
	bc.Blocks = append(bc.Blocks, block.(*Block))
}

func (bc *BlockChain) GetBlock(high int) IBlock {
	if high > len(bc.Blocks) {
		return nil
	}
	return bc.Blocks[high]
}

func (bc *BlockChain) GetLatestBlock() IBlock {
	return bc.Blocks[len(bc.Blocks)-1]
}

func (bc *BlockChain) GetBlockHigh() int {
	return len(bc.Blocks)
}

func (bc *BlockChain) GenerateNextBlock(data string) IBlock {
	lastBlock := bc.GetLatestBlock()
	newBlock := NewBlock(lastBlock.GetHigh()+1, time.Now().Unix(), lastBlock.GetHash(), "", data)
	newBlock.Hash = newBlock.CalculateHash()
	return newBlock
}

func (bc *BlockChain) ReplaceChain(chain IBlockchain) {
}

type Block struct {
	High      int64  `json:"high"`
	TimeStamp int64  `json:"time_stamp"`
	PreHash   string `json:"pre_hash"`
	Hash      string `json:"hash"`
	Data      string `json:"data"`
}

func NewBlock(high, timeStamp int64, PreviousHash, hash, data string) *Block {
	return &Block{
		High:      high,
		TimeStamp: timeStamp,
		PreHash:   PreviousHash,
		Hash:      hash,
		Data:      data,
	}
}

func GenerateNextBlock(lastBlock *Block, data string) *Block {
	block := NewBlock(lastBlock.High+1, time.Now().Unix(), lastBlock.Hash, "", data)
	block.Hash = block.CalculateHash()
	return block
}

func (b *Block) GetHigh() int64 {
	return b.High
}

func (b *Block) GetTimeStamp() int64 {
	return b.TimeStamp
}

func (b *Block) GetPreviousHash() string {
	return b.PreHash
}

func (b *Block) GetHash() string {
	return b.Hash
}

func (b *Block) GetData() string {
	return b.Data
}

func (b *Block) GetString() string {
	return fmt.Sprintf("Block[High: %d, TimeStamp: %d, PreHash: %s, Hash: %s, Data: %s]", b.High, b.TimeStamp, b.PreHash, b.Hash, b.Data)
}

func (b *Block) CalculateHash() string {
	code := strconv.Itoa(int(b.High)) + b.PreHash + b.Data
	hash := sha256.Sum256([]byte(code))
	return hex.EncodeToString(hash[:])
}

func (b *Block) IsValidBlock(lastBlock IBlock) bool {
	if lastBlock.GetHigh()+1 != b.GetHigh() {
		return false
	} else if lastBlock.GetHash() != b.GetPreviousHash() {
		return false
	} else if b.CalculateHash() != b.GetHash() {
		return false
	}

	return true
}

func CalculateHash(index int64, previousHash string, data string) string {
	sb := strings.Builder{}
	sb.WriteString(strconv.Itoa(int(index)))
	sb.WriteString(previousHash)
	sb.WriteString(data)
	hash := sha256.Sum256([]byte(sb.String()))
	return hex.EncodeToString(hash[:])
}
