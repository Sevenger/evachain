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

	http.HandleFunc("/blocks", handleBlocks)
	http.HandleFunc("/blocks/add", handleAddBlock)
	http.HandleFunc("/peers", handlePeers)
	http.HandleFunc("/peers/add", handleAddPeer)
	go func() {
		logMsg("Listen HTTP on", *httpAddr)
		logFatal("start api sever error", http.ListenAndServe(*httpAddr, nil))
	}()

	http.Handle("/", websocket.Handler(handleP2P))
	logMsg("Listen P2P on", *p2pAddr)
	logFatal("start P2P server error", http.ListenAndServe(*p2pAddr, nil))
}

func logFatal(msg string, err error) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}

func logMsg(msg ...interface{}) {
	log.Println(msg...)
}

func logMsgf(f string, v ...interface{}) {
	log.Printf(f, v...)
}
