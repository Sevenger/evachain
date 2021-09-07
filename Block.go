package main

import (
	"time"
)

type BlockChain struct {
	Blocks []*Block `json:"blocks"`
}

func NewBlockChain() *BlockChain {
	return &BlockChain{
		Blocks: []*Block{
			NewBlock(
				0,
				1630480460,
				"",
				"280074b880e45acfced3772b80d83a08fd90e0affd69e15dcc583c66e40863d4",
				"Hello, BlockChain!",
			),
		},
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
	switch block.(type) {
	case *Block:
		bc.Blocks = append(bc.Blocks, block.(*Block))
	case BlockImpl:
		bc.Blocks = append(bc.Blocks, &Block{block.(BlockImpl)})
	}
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

func (bc *BlockChain) GetAllBlocks() []IBlock {
	blocks := make([]IBlock, len(bc.Blocks))
	for i, v := range bc.Blocks {
		blocks[i] = v
	}
	return blocks
}

func (bc *BlockChain) GetBlockHigh() int64 {
	return int64(len(bc.Blocks))
}

func (bc *BlockChain) GenerateNextBlock(data string) IBlock {
	lastBlock := bc.GetLatestBlock()
	newBlock := NewBlock(lastBlock.GetHigh()+1, time.Now().Unix(), lastBlock.GetHash(), "", data)
	newBlock.Hash = newBlock.CalculateHash()
	return newBlock
}

func (bc *BlockChain) ReplaceBlocks(blocks []IBlock) {
	tmpBlocks := make([]*Block, len(blocks))
	for i, v := range blocks {
		tmpBlocks[i] = &Block{v.(BlockImpl)}
	}
	bc.Blocks = tmpBlocks
}

type Block struct {
	BlockImpl
}

func NewBlock(high, timeStamp int64, PreviousHash, hash, data string) *Block {
	return &Block{
		BlockImpl{
			High:      high,
			TimeStamp: timeStamp,
			PreHash:   PreviousHash,
			Hash:      hash,
			Data:      data,
		},
	}
}
