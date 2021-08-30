package evachain

import (
	"crypto/sha256"
	"fmt"
	"github.com/sevenger/evachain/model"
)

func CalculateHash(index int, previousHash string, data model.Data) string {
	dataByte := fmt.Sprintf("%d%s%v", index, previousHash, data)
	hash := sha256.Sum256([]byte(dataByte))
	return fmt.Sprintf("%x", hash)
}

func CalculateHashForBlock(block *model.Block) string {
	return CalculateHash(block.Index, block.PreviewHash, block.Data)
}
