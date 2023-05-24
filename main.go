package main

import (
	"fmt"

	"github.com/jootd/dstore/p2p"
)

func main() {

	transport := p2p.NewTCPTransport(p2p.TCPTransportOpts{
		ListenAddr: ":3000",
	})

	err := transport.ListenAndAccept()
	if err != nil {
		fmt.Println(err.Error())
	}
	select {}

}
