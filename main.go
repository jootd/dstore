package main

import (
	"fmt"
	"log"

	"github.com/jootd/dstore/p2p"
)

func OnPeer(p p2p.Peer) error {
	fmt.Println("doing some logic with the peer outside of TCPTransport")

	return nil
}

func makeServer(listenAddr string, nodes ...string) *FileServer {
	tcpTransportOpts := p2p.TCPTransportOpts{
		ListenAddr:    listenAddr,
		HandShakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
		// TODO: OnPeer Func
	}
	tcpTransport := p2p.NewTCPTransport(tcpTransportOpts)

	fileServerOpts := FileServerOpts{
		StorageRoot:       listenAddr + "_network",
		PathTransformFunc: CASPathTransformFunc,
		BootstrapNodes:    nodes,
		Transport:         tcpTransport,
	}

	s := NewFileServer(fileServerOpts)
	tcpTransport.OnPeer = s.OnPeer
	return s
}

func main() {

	s1 := makeServer(":3000")
	s2 := makeServer(":4001", ":3000")

	go func() {
		log.Fatal(s1.Start())
	}()

	go func() {

		log.Fatal(s2.Start())
	}()

	select {}

}
