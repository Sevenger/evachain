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
func handleBlocks(w http.ResponseWriter, _ *http.Request) {
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
		logMsg("invalid block data:", err)
		return
	}

	block := GenerateNextBlock(params.Data)
	AddBlock(block)
	// 向其他p2p节点广播消息
	BoardCast(ResponseLatestMsg())
}

//handlePeers 查看p2p节点信息
func handlePeers(w http.ResponseWriter, _ *http.Request) {
	var peers []string
	for _, socket := range sockets {
		peers = append(peers, socket.RemoteAddr().String())
	}

	bs, _ := json.Marshal(peers)
	w.Write(bs)
}

//handleAddPeer 添加p2p节点
func handleAddPeer(_ http.ResponseWriter, r *http.Request) {
	var params struct {
		Peer string `json:"peer"`
	}
	logMsg("jia jiedian")
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	if err := decoder.Decode(&params); err != nil {
		logMsg("invalid peer data:", err)
		return
	}

	connectToPeer(params.Peer)
}

//handleP2P websocket
func handleP2P(ws *websocket.Conn) {
	var (
		res  = &Response{}
		peer = ws.LocalAddr().String()
	)
	sockets = append(sockets, ws)

	for {
		var msg []byte
		err := websocket.Message.Receive(ws, &msg)
		if err == io.EOF {
			logMsgf("peer[%s] closed", peer)
			break
		}
		if err != nil {
			logMsgf("can't receive msg from peer[%s]:", peer, err.Error())
			break
		}

		logMsgf("received[from %s]: %s.\n", peer, msg)
		if err = json.Unmarshal(msg, res); err != nil {
			log.Println("invalid msg")
			continue
		}

		switch res.Type {
		case queryLastBlock:
			res.Type = responseBlockchain
			bs := ResponseLatestMsg()

			logMsg("responseLatestMsg:", bs)
			ws.Write(bs)

		case queryAllBlock:
			res.Type = responseBlockchain
			d, _ := json.Marshal(blockchain)
			res.Data = string(d)
			bs, _ := json.Marshal(res)

			logMsg("responseChainMsg:", bs)
			ws.Write(bs)

		case responseBlockchain:
			handleBlockchainResponse([]byte(res.Data))
		}
	}
}

//handleBlockchainResponse 处理response
func handleBlockchainResponse(msg []byte) {
	var receivedBlocks = []*Block{}

	if err := json.Unmarshal(msg, &receivedBlocks); err != nil {
		log.Println("invalid chain:", err)
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
			BoardCast([]byte(fmt.Sprintf("{\"type\": %d}", queryAllBlock)))
		} else {
			// 检查是否需要替换链
			ReplaceChain(receivedBlocks)
		}
	} else {
		log.Println("接收的区块链更短，不执行操作")
	}
}

//connectToPeer 链接p2p节点
func connectToPeer(peer string) {
	ws, err := websocket.Dial(peer, "", peer)
	if err != nil {
		logMsg("invalid peer:", err)
		return
	}

	go handleP2P(ws)

	//log.Println("query latest block.")
}

//BoardCast 向所有P2P节点发送信息
func BoardCast(msg []byte) {
	for i, socket := range sockets {
		if _, err := socket.Write(msg); err != nil {
			logMsgf("peer[%s] disconnect", socket.RemoteAddr().String())
			// 去除节点
			sockets = append(sockets[:i], sockets[i+1:]...)
		}
	}
}

const (
	queryAllBlock = iota
	queryLastBlock
	responseBlockchain
)

type Response struct {
	Type int    `json:"type,omitempty"`
	Data string `json:"data,omitempty"`
}

func ResponseLatestMsg() []byte {
	d, _ := json.Marshal(GetLatestBlock())
	response := &Response{
		Type: responseBlockchain,
		Data: string(d),
	}
	bs, _ := json.Marshal(response)
	return bs
}
