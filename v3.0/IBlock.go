package main

type IBlockchain interface {
	AddBlock(block IBlock)
	GetGenesisBlock() IBlock
	GetBlock(high int) IBlock
	GetLatestBlock() IBlock
	GenerateNextBlock(data string) IBlock
	GetBlockHigh() int
	IsValidChain() bool
	ReplaceChain(chain IBlockchain)
}

type IBlock interface {
	GetHigh() int64
	GetTimeStamp() int64
	GetPreviousHash() string
	GetHash() string
	GetData() string
	GetString() string
	CalculateHash() string
	IsValidBlock(lastBlock IBlock) bool
}
