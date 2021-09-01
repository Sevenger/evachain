package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/websocket"
	"io"
	"log"
	"net/http"
)

var sockets []*websocket.Conn

//handleBlocks 查看链数据
func handleBlocks(w http.ResponseWriter, r *http.Request) {
	bs, _ := json.Marshal(blockchain)
	w.Write(bs)
}

//handleAddBlock 添加区块
func handleAddBlock(w http.ResponseWriter, r *http.Request) {
	var params struct {
		Data string `json:"data"`
	}
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := decoder.Decode(&params); err != nil {
		logErr("decoder error", err)
	}

	block := GenerateNextBlock(params.Data)
	AddBlock(block)
	// 向其他p2p节点广播消息
	BoardCast(ResponseLatestMsg())
}

//handlePeers 查看p2p节点信息
func handlePeers(w http.ResponseWriter, r *http.Request) {
	var peers []string
	for _, socket := range sockets {
		peers = append(peers, socket.RemoteAddr().String())
	}

	peersJson, _ := json.Marshal(peers)
	w.Write(peersJson)
}

//handleAddPeer 添加p2p节点
func handleAddPeer(w http.ResponseWriter, r *http.Request) {
	var params struct {
		Peer string `json:"peer"`
	}
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	if err := decoder.Decode(&params); err != nil {
		log.Println("invalid peer", err)
		return
	}

	connectToPeers([]string{params.Peer})
}

//handleP2P websocket
func handleP2P(ws *websocket.Conn) {
	var (
		v    = &ResponseBlockchain{}
		peer = ws.LocalAddr().String()
	)
	sockets = append(sockets, ws)

	for {
		var msg []byte
		err := websocket.Message.Receive(ws, &msg)
		if err == io.EOF {
			log.Println("p2p节点关闭", peer)
			break
		}
		if err != nil {
			log.Println("无法接收p2p信息", peer, err.Error())
			break
		}

		log.Printf("接收信息[来自 %s]: %s.\n", peer, msg)
		if err = json.Unmarshal(msg, v); err != nil {
			log.Println("非法p2p信息")
			continue
		}

		switch v.Type {
		case queryLatest:
			v.Type = responseBlockchain
			bs := ResponseLatestMsg()
			log.Printf("responseLatestMsg: %s\n", bs)
			ws.Write(bs)

		case queryAll:
			v.Type = responseBlockchain
			d, _ := json.Marshal(blockchain)
			v.Data = string(d)
			bs, _ := json.Marshal(v)
			log.Printf("responseChainMsg: %s", bs)
			ws.Write(bs)

		case responseBlockchain:
			handleBlockchainResponse(msg)
		}
	}
}

func handleBlockchainResponse(msg []byte) {
	var receivedBlocks = []*Block{}

	if err := json.Unmarshal(msg, &receivedBlocks); err != nil {
		log.Println("非法区块链", err)
		return
	}

	receivedBlock := receivedBlocks[len(receivedBlocks)-1]
	heldBlock := GetLatestBlock()
	// 接收的区块的下标必须比当前的大
	if receivedBlock.Index > heldBlock.Index {
		// 接收的区块的previousHash必须是当前的hash
		if heldBlock.Hash == receivedBlock.PreviousHash {
			blockchain = append(blockchain, receivedBlock)
			// 接收的区块长度为1时查询其他节点
		} else if len(receivedBlocks) == 1 {
			BoardCast([]byte(fmt.Sprintf("{\"type\": %d}", queryAll)))
		} else {
			// 检查是否需要替换链
			ReplaceChain(receivedBlocks)
		}
	} else {
		log.Println("接收的区块链更短，不执行操作")
	}
}

func connectToPeers(peers []string) {
	for _, peer := range peers {
		if peer == "" {
			continue
		}
		ws, err := websocket.Dial(peer, "", peer)
		if err != nil {
			log.Println("dial to peer", err)
			continue
		}

		go handleP2P(ws)
	}
}

//BoardCast 向所有P2P节点发送信息
func BoardCast(msg []byte) {
	for i, socket := range sockets {
		if _, err := socket.Write(msg); err != nil {
			logMsgf("peer[%s] disconnect", socket.RemoteAddr().String())
			sockets = append(sockets[:i], sockets[i+1:]...)
		}
	}
}

const (
	queryAll = iota
	queryLatest
	responseBlockchain
)

type ResponseBlockchain struct {
	Type int    `json:"type,omitempty"`
	Data string `json:"data,omitempty"`
}

func ResponseLatestMsg() []byte {
	blockJson, _ := json.Marshal(GetLatestBlock())
	response := &ResponseBlockchain{
		Type: responseBlockchain,
		Data: string(blockJson),
	}
	responseJson, _ := json.Marshal(response)
	return responseJson
}
