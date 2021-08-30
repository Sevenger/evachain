package evachain

import "github.com/sevenger/evachain/model"

var EvaChain = []*model.Block{GetGenesisBlock()}

func GetGenesisBlock() *model.Block {
	return model.NewBlock(0, "1630310110", "Hello, blockchain", "67343ee4fe230793916042affe72699e2860ed4e3c7c39ccebd4efb6c48e81b5", "")
}

func IsValidBlock(block, lastBlock *model.Block) bool {
	if block == nil || lastBlock == nil {
		return false
	} else if lastBlock.Index+1 != block.Index {
		return false
	} else if lastBlock.Hash != block.PreviewHash {
		return false
	} else if block.Hash != CalculateHashForBlock(block) {
		return false
	}
	return true
}
