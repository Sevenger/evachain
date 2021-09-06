package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
)

type Node struct {
	blockchain IBlockchain
	peers      []*websocket.Conn
}

func (n *Node) HandleGetBlocks(w http.ResponseWriter, request *http.Request) {
	bs, _ := json.Marshal(n.blockchain)
	w.Write(bs)
}

func (n *Node) HandleGetPeers(writer http.ResponseWriter, request *http.Request) {
	var addrs []string
	for _, peer := range n.peers {
		if peer.IsClientConn() {
			addrs = append(addrs, peer.LocalAddr().String())
		} else {
			addrs = append(addrs, peer.Request().RemoteAddr)
		}
	}
	bs, _ := json.Marshal(addrs)
	writer.Write(bs)
}

func (n *Node) HandleAddBlock(writer http.ResponseWriter, request *http.Request) {
	var params struct {
		Data string `json:"data"`
	}
	decoder := json.NewDecoder(request.Body)
	defer request.Body.Close()
	if err := decoder.Decode(&params); err != nil {
		return
	}

	block := n.blockchain.GenerateNextBlock(params.Data)
	n.blockchain.AddBlock(block)
	n.Broadcast(n.QueryLatestBlockResponse())
}

func (n *Node) HandleAddPeer(w http.ResponseWriter, r *http.Request) {
	var params struct {
		Addr string `json:"addr"`
	}
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	if err := decoder.Decode(&params); err != nil {
		return
	}
	n.ConnectPeer(params.Addr)
}

func (n *Node) ConnectPeer(addr string) {
	peer, err := websocket.Dial(addr, "", addr)
	if err != nil {
		return
	}
	go n.HandleBroadcast(peer)
}

func (n *Node) HandleBroadcast(peer *websocket.Conn) {
	var msg Message
	n.peers = append(n.peers, peer)
	peer.Write(n.QueryLatestBlockMsg())
	for {
		var metaData MetaData
		if err := websocket.Message.Receive(peer, &metaData); err != nil {
			fmt.Println("error happened")
			break
		}
		if err := json.Unmarshal(metaData, &msg); err != nil {
			fmt.Println("continue")
			continue
		}

		fmt.Println(msg)

		switch msg.Type {
		case MsgQueryLatestBlock:
			peer.Write(n.QueryLatestBlockResponse())

		case MsgQueryAllBlock:
			peer.Write(n.QueryAllBlockResponse())

		case MsgResponse:
			var receivedChain IBlockchain
			if err := json.Unmarshal(metaData, &receivedChain); err != nil {
				log.Println(err)
				continue
			}

			receivedBlock := receivedChain.GetLatestBlock()
			heldBlock := n.blockchain.GetLatestBlock()
			if receivedBlock.GetHigh() > heldBlock.GetHigh() {
				switch {
				case receivedBlock.GetPreviousHash() == heldBlock.GetPreviousHash():

				case receivedChain.GetBlockHigh() == 1:
					n.Broadcast(n.QueryAllBlockMsg())

				default:
					n.blockchain = IBlockchain(receivedChain)
				}
			}
		}
	}
}

func (n *Node) Broadcast(data MetaData) {
	for i, peer := range n.peers {
		if _, err := peer.Write(data); err != nil {
			fmt.Println(err)
			n.peers = append(n.peers[:i], n.peers[i+1:]...)
		}
	}
}

func NewNode() *Node {
	return &Node{
		blockchain: NewBlockChain(),
		peers:      []*websocket.Conn{},
	}
}

const (
	MsgQueryAllBlock = iota + 1
	MsgQueryLatestBlock
	MsgResponse
)

type Message struct {
	Type int    `json:"type"`
	Data string `json:"data"`
}

type MetaData = []byte

func (*Node) QueryAllBlockMsg() MetaData {
	msg := &Message{Type: MsgQueryAllBlock}
	return MsgToMeta(msg)
}

func (n *Node) QueryAllBlockResponse() MetaData {
	d, _ := json.Marshal(n.blockchain)
	msg := &Message{
		Type: MsgResponse,
		Data: string(d),
	}
	return MsgToMeta(msg)
}

func (n *Node) QueryLatestBlockMsg() MetaData {
	msg := &Message{Type: MsgQueryLatestBlock}
	return MsgToMeta(msg)
}

func (n *Node) QueryLatestBlockResponse() MetaData {
	blockchain := n.blockchain.(*BlockChain)
	blockchain.Blocks = blockchain.Blocks[len(blockchain.Blocks)-1:]
	d, _ := json.Marshal(blockchain)
	msg := &Message{
		Type: MsgResponse,
		Data: string(d),
	}
	return MsgToMeta(msg)
}

func MsgToMeta(msg *Message) MetaData {
	meta, _ := json.Marshal(msg)
	return meta
}
