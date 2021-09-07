package main

import (
	"golang.org/x/net/websocket"
	"net/http"
)

type INode interface {
	HandleGetBlocks(w http.ResponseWriter, r *http.Request)
	HandleAddBlock(w http.ResponseWriter, r *http.Request)
	HandleGetPeers(w http.ResponseWriter, r *http.Request)
	HandleAddPeer(w http.ResponseWriter, r *http.Request)
	ConnectPeer(addr string)
	HandleBroadcast(peer *websocket.Conn)
	Broadcast(data MetaData)
}
