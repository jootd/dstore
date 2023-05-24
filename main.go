package main

import (
	"fmt"

	"github.com/jootd/dstore/p2p"
)

func OnPeer(p p2p.Peer) error {
	fmt.Println("doing some logic with the peer outside of TCPTransport")

	return nil
}

func main() {

	opts := p2p.TCPTransportOpts{
		ListenAddr:    ":3000",
		Decoder:       p2p.DefaultDecoder{},
		HandShakeFunc: p2p.NOPHandshakeFunc,
		OnPeer:        OnPeer,
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
