package main

import (
	"flag"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
)

var (
	httpAddr = flag.String("api", ":8080", "api server address")
	p2pAddr  = flag.String("p2p", ":7070", "p2p server address")
)

func main() {
	flag.Parse()

	var node INode = NewNode()

	http.HandleFunc("/blocks", node.HandleGetBlocks)
	http.HandleFunc("/blocks/add", node.HandleAddBlock)
	http.HandleFunc("/peers", node.HandleGetPeers)
	http.HandleFunc("/peers/add", node.HandleAddPeer)
	go func() {
		log.Println("Listen HTTP on", *httpAddr)
		logFatal("start api sever error", http.ListenAndServe(*httpAddr, nil))
	}()

	http.Handle("/", websocket.Handler(node.HandleBroadcast))
	log.Println("Listen P2P on", *p2pAddr)
	logFatal("start P2P server error", http.ListenAndServe(*p2pAddr, nil))
}

func logFatal(msg string, err error) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
