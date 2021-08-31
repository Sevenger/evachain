package main

import (
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"golang.org/x/net/websocket"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	queryAll = iota
	queryLast
)

var GenesisBlock = &Block{
	Index:        0,
	TimeStamp:    1630392503,
	PreviousHash: "",
	Hash:         "280074b880e45acfced3772b80d83a08fd90e0affd69e15dcc583c66e40863d4",
	Data:         "Hello, eva-chain!",
}

var (
	EvaChain     = []*Block{GenesisBlock}
	socket       []*websocket.Conn
	httpAddr     = flag.String("api", ":3001", "api server address")
	p2pAddr      = flag.String("p2p", ":6001", "p2p server address")
	initialPeers = flag.String("peers", "ws://localhost:6001", "initial peers")
)

type Block struct {
	Index        int64  `json:"index,omitempty"`
	TimeStamp    int64  `json:"time_stamp,omitempty"`
	PreviousHash string `json:"previous_hash,omitempty"`
	Hash         string `json:"hash,omitempty"`
	Data         string `json:"data,omitempty"`
}

func (b *Block) ToString() string {
	return fmt.Sprintf("Index: %d, TimeStamp: %d, PreviousHash: %s, Hash: %s, Data: %s",
		b.Index, b.TimeStamp, b.PreviousHash, b.Hash, b.Data)
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

func GenerateNextBlock(previousHash string, data string) *Block {
	index := int64(len(EvaChain) + 1)
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
	return CalculateHash(block.Index, block.Hash, block.Data)
}

func IsValidBlock(newBlock, lastBlock *Block) bool {
	if newBlock.Index != lastBlock.Index+1 {
		return false
	} else if newBlock.PreviousHash != lastBlock.Hash {
		return false
	} else if newBlock.Hash != CalculateHashForBlock(newBlock) {
		return false
	}
	return true
}

func connectToPeers(peersAddr []string) {
	for _, peer := range peersAddr {
		if peer == "" {
			continue
		}
		ws, err := websocket.Dial(peer, "", peer)
		if err != nil {
			log.Fatalln(err)
		}
		initConnection(ws)
	}
}

func initConnection(ws *websocket.Conn) {

}
func main() {
	flag.Parse()
}
