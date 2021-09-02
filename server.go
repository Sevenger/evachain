package main

import (
	"encoding/json"
	"golang.org/x/net/websocket"
	"io"
	"log"
	"net/http"
	"strings"
)

var sockets []*websocket.Conn

//handleBlocks 查看链数据
func handleBlocks(w http.ResponseWriter, _ *http.Request) {
	bs, _ := json.Marshal(EvaChain)
	w.Write(bs)
}

//handleAddBlock 添加区块
func handleAddBlock(_ http.ResponseWriter, r *http.Request) {
	var params struct {
		Data string `json:"data"`
	}
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := decoder.Decode(&params); err != nil {
		logMsg("invalid block data:", err)
		return
	}

	logMsgf("add block[data: %s]", params.Data)
	block := GenerateNextBlock(params.Data)
	AddBlock(block)
	// 向其他p2p节点广播消息
	BoardCast(ResponseLatestMsg())
}

//handlePeers 查看p2p节点信息
func handlePeers(w http.ResponseWriter, _ *http.Request) {
	var peers []string
	for _, socket := range sockets {
		if socket.IsClientConn() {
			peers = append(peers, strings.Replace(socket.LocalAddr().String(), "ws://", "", 1))
		} else {
			peers = append(peers, socket.Request().RemoteAddr)
		}
	}

	bs, _ := json.Marshal(peers)
	w.Write(bs)
}

//handleAddPeer 添加p2p节点
func handleAddPeer(_ http.ResponseWriter, r *http.Request) {
	var params struct {
		Peer string `json:"peer"`
	}
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	if err := decoder.Decode(&params); err != nil {
		logMsg("invalid peer data:", err)
		return
	}

	logMsgf("add peer[%s]", params.Peer)
	connectToPeer(params.Peer)
}

//connectToPeer 链接p2p节点
func connectToPeer(peer string) {
	ws, err := websocket.Dial(peer, "", peer)
	if err != nil {
		logMsg("invalid peer:", err)
		return
	}

	go handleP2P(ws)
}

//handleP2P websocket
func handleP2P(ws *websocket.Conn) {
	var (
		msg  = &Msg{}
		peer = ws.LocalAddr().String()
	)
	sockets = append(sockets, ws)

	// 新连接时查询该节点的链
	log.Println("query latest block")
	ws.Write(QueryLatestMsg())

	for {
		var receive []byte
		err := websocket.Message.Receive(ws, &receive)
		if err == io.EOF {
			logMsgf("peer[%s] closed", peer)
			break
		}
		if err != nil {
			logMsgf("can't receive msg from peer[%s]:", peer, err.Error())
			break
		}

		logMsgf("received[from %s]: %s.\n", peer, receive)
		if err = json.Unmarshal(receive, msg); err != nil {
			log.Println("invalid received msg")
			continue
		}

		switch msg.Type {
		case queryLastBlock:
			bs := ResponseLatestMsg()

			logMsgf("responseLatestMsg: %s\n", bs)
			ws.Write(bs)

		case queryAllBlock:
			bs := ResponseAllMsg()

			logMsgf("responseChainMsg: %s\n", bs)
			ws.Write(bs)

		case responseBlockchain:
			handleBlockchainResponse([]byte(msg.Data))
		}
	}
}

//handleBlockchainResponse 处理response
func handleBlockchainResponse(msg []byte) {
	var receivedChain Chain
	if err := json.Unmarshal(msg, &receivedChain); err != nil {
		log.Println("invalid chain:", err)
		return
	}

	receivedBlock := receivedChain[len(receivedChain)-1]
	heldBlock := GetLatestBlock()
	// 接收的区块的下标比当前的大时，决定是否同步链
	if receivedBlock.Index > heldBlock.Index {
		switch {
		case receivedBlock.PreviousHash == heldBlock.Hash: // 接收的区块刚好是当前节点的下一个
			logMsgf("add block[data: %s] from received", receivedBlock.Data)
			EvaChain = append(EvaChain, receivedBlock)

		case len(receivedChain) == 1: // 因为传进来的是最后一个节点，所以请求所有链
			BoardCast(QueryAllMsg())

		default: // 检查是否需要替换链
			ReplaceChain(receivedChain)
		}
	} else {
		log.Println("received chain not longer than current chain")
	}
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
	queryLastBlock = iota
	queryAllBlock
	responseBlockchain
)

type Msg struct {
	Type int    `json:"type"`
	Data string `json:"data"`
}

func QueryLatestMsg() []byte {
	msg := &Msg{Type: queryLastBlock}
	bs, _ := json.Marshal(msg)
	return bs
}

func QueryAllMsg() []byte {
	msg := &Msg{Type: queryAllBlock}
	bs, _ := json.Marshal(msg)
	return bs
}

func ResponseLatestMsg() []byte {
	d, _ := json.Marshal(EvaChain[len(EvaChain)-1:])
	msg := &Msg{
		Type: responseBlockchain,
		Data: string(d),
	}
	bs, _ := json.Marshal(msg)
	return bs
}

func ResponseAllMsg() []byte {
	d, _ := json.Marshal(EvaChain)
	msg := &Msg{
		Type: responseBlockchain,
		Data: string(d),
	}
	bs, _ := json.Marshal(msg)
	return bs
}
