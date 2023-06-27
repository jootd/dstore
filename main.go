package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jootd/dstore/p2p"
)

func OnPeer(p p2p.Peer) error {
	fmt.Println("doing some logic with the peer outside of TCPTransport")

	return nil
}

func main() {

	tcpTransportOpts := p2p.TCPTransportOpts{
		ListenAddr:    ":3000",
		HandShakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
		// TODO: OnPeer Func
	}
	tcpTransport := p2p.NewTCPTransport(tcpTransportOpts)

	fileServerOpts := FileServerOpts{
		StorageRoot:       "3000_network",
		PathTransformFunc: CASPathTransformFunc,
		Transport:         tcpTransport,
	}

	fs := NewFileServer(fileServerOpts)

	go func() {
		time.Sleep(time.Second * 5)
		fs.Stop()

	}()

	if err := fs.Start(); err != nil {
		log.Fatal(err)

	}

	select {}
}
