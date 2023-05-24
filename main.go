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
	}

	transport := p2p.NewTCPTransport(opts)

	err := transport.ListenAndAccept()
	if err != nil {
		fmt.Println(err.Error())
	}
	select {}

}
