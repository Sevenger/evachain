package main

import (
	"flag"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
	"strings"
)

var (
	httpAddr     = flag.String("api", ":3001", "api server address")
	p2pAddr      = flag.String("p2p", ":6001", "p2p server address")
	initialPeers = flag.String("peers", "ws://localhost:6001", "initial peers")
)

func main() {
	flag.Parse()

	connectToPeers(strings.Split(*initialPeers, ","))

	http.HandleFunc("/blocks", handleBlocks)
	http.HandleFunc("/blocks/add", handleAddBlock)
	http.HandleFunc("/peers", handlePeers)
	http.HandleFunc("/peers/add", handleAddPeer)
	go func() {
		logMsg("Listen HTTP on", *httpAddr)
		logErr("start api sever error", http.ListenAndServe(*httpAddr, nil))
	}()

	http.Handle("/", websocket.Handler(handleP2P))
	logMsg("Listen P2P on", *p2pAddr)
	logErr("start P2P server error", http.ListenAndServe(*p2pAddr, nil))
}

func logErr(msg interface{}, err error) {
	if err != nil {
		log.Fatalln(msg, err.Error())
	}
}

func logMsg(msg ...interface{}) {
	log.Println(msg...)
}

func logMsgf(f string, v ...interface{}) {
	log.Printf(f, v...)
}
