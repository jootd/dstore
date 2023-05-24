package main

import (
	"fmt"

	"github.com/jootd/dstore/p2p"
)

func main() {

	opts := p2p.TCPTransportOpts{
		ListenAddr:    ":3000",
		Decoder:       p2p.DefaultDecoder{},
		HandShakeFunc: p2p.NOPHandshakeFunc,
		OnPeer:        func(p2p.Peer) error { return fmt.Errorf("failed the onpeer func") },
	}

	tr := p2p.NewTCPTransport(opts)

	go func() {
		for {
			msg := <-tr.Consume()
			fmt.Printf("%+v\n", msg)
		}

	}()

	err := tr.ListenAndAccept()
	if err != nil {
		fmt.Println(err.Error())
	}
	select {}

}
