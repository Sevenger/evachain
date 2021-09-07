package main

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
)

type IBlockchain interface {
	AddBlock(block IBlock)
	GetGenesisBlock() IBlock
	GetBlock(high int) IBlock
	GetLatestBlock() IBlock
	GetAllBlocks() []IBlock
	GenerateNextBlock(data string) IBlock
	GetBlockHigh() int64
	IsValidChain() bool
	ReplaceBlocks(blocks []IBlock)
}

type IBlock interface {
	GetHigh() int64
	GetTimeStamp() int64
	GetPreHash() string
	GetHash() string
	GetData() string
	CalculateHash() string
	IsValidBlock(lastBlock IBlock) bool
}

type BlockImpl struct {
	High      int64
	TimeStamp int64
	PreHash   string
	Hash      string
	Data      string
}

func (impl BlockImpl) GetHigh() int64 {
	return impl.High
}

func (impl BlockImpl) GetTimeStamp() int64 {
	return impl.TimeStamp
}

func (impl BlockImpl) GetPreHash() string {
	return impl.PreHash
}

func (impl BlockImpl) GetHash() string {
	return impl.Hash
}

func (impl BlockImpl) GetData() string {
	return impl.Data
}

func (impl BlockImpl) CalculateHash() string {
	code := strconv.Itoa(int(impl.High)) + impl.PreHash + impl.Data
	hash := sha256.Sum256([]byte(code))
	return hex.EncodeToString(hash[:])
}

func (impl BlockImpl) IsValidBlock(lastBlock IBlock) bool {
	if lastBlock.GetHigh()+1 != impl.GetHigh() {
		return false
	} else if lastBlock.GetHash() != impl.GetPreHash() {
		return false
	} else if impl.CalculateHash() != impl.GetHash() {
		return false
	}

	return true
}
