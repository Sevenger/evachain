package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const difficulty = 5

type IBlock interface {
	Pow()
}

func main() {
	var b IBlock = B{}
	b.Pow()
}

type B struct {
	Index int64
}

func (B) Pow() {

}

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
	Nonce        string `json:"nonce"`
	Difficulty   int64  `json:"difficulty"`
}

func NewBlock(index, timeStamp int64, previousHash, hash, data, nonce string, difficulty int64) *Block {
	return &Block{
		Index:        index,
		TimeStamp:    timeStamp,
		PreviousHash: previousHash,
		Hash:         hash,
		Data:         data,
		Nonce:        nonce,
		Difficulty:   difficulty,
	}
}

func GenerateBlock(lastBlock *Block, data string) *Block {
	newBlock := NewBlock(lastBlock.Index+1, time.Now().Unix(), lastBlock.Hash, "", data, "", difficulty)
	for i := int64(0); ; i++ {
		nonce := strconv.FormatInt(i, 16)
		newBlock.Nonce = nonce
		hash := CalculateHashForBlock(newBlock)
		if !IsValidHash(hash, newBlock.Difficulty) {
			fmt.Printf("\rCalculateHash[%s], do more work!", hash)
		} else {
			fmt.Println()
			logMsgf("CalculateHash[%s], done work!", hash)
			newBlock.Hash = hash
			break
		}
	}
	return newBlock
}

func GenerateNextBlock(data string) *Block {
	return GenerateBlock(GetLatestBlock(), data)
}

//CalculateHash 计算哈希规则：下标+前块哈希+数据+随机数
func CalculateHash(index int64, previousHash string, data string, nonce string) string {
	sb := strings.Builder{}
	sb.WriteString(strconv.Itoa(int(index)))
	sb.WriteString(previousHash)
	sb.WriteString(data)
	sb.WriteString(nonce)
	hash := sha256.Sum256([]byte(sb.String()))
	return hex.EncodeToString(hash[:])
}

func CalculateHashForBlock(block *Block) string {
	return CalculateHash(block.Index, block.PreviousHash, block.Data, block.Nonce)
}

//IsValidHash 计算有difficulty个前导0判断哈希是否合法
func IsValidHash(hash string, difficulty int64) bool {
	prefix := strings.Repeat("0", int(difficulty))
	return strings.HasPrefix(hash, prefix)
}

func IsValidBlock(newBlock, lastBlock *Block) bool {
	if newBlock.Index != lastBlock.Index+1 {
		logMsg("invalid index!")
		return false
	} else if newBlock.PreviousHash != lastBlock.Hash {
		logMsg("invalid previous hash!")
		return false
	} else if newBlock.Hash != CalculateHashForBlock(newBlock) {
		logMsg("invalid hash!")
		return false
	}
	return true
}

func IsValidChain(chain Chain) bool {
	if *chain[0] != *GenesisBlock {
		logMsg("is not the same chain!")
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

func AddBlock(newBlock *Block) {
	lastBlock := GetLatestBlock()
	if IsValidBlock(newBlock, lastBlock) {
		EvaChain = append(EvaChain, newBlock)
	} else {
		logMsgf("invalid block[%+v]\n", *newBlock)
		logMsgf("last block[%+v]\n", *lastBlock)
	}
}

func ReplaceChain(chain Chain) {
	if IsValidChain(chain) && len(chain) > len(EvaChain) {
		logMsg("replace the chain")
		EvaChain = chain
		BoardCast(ResponseLatestMsg())
	}
}
