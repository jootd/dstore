package p2p

import "net"

type Peer interface {
	Send([]byte) error
	RemoteAddr() net.Addr
	Close() error
}

type Transport interface {
	ListenAndAccept() error
	Consume() <-chan RPC
	Close() error
	Dial(string) error
}
