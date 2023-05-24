package p2p

import (
	"fmt"
	"net"
	"sync"
)

type TCPPeer struct {

	// conn is underlying connection of peer
	conn net.Conn

	// dial, retrieve a conn => outbound == true
	// accept,  retrieve a conn => outbound == false
	outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

type TCPTransportOpts struct {
	ListenAddr    string
	HandShakeFunc HandShakeFunc
	Decoder       Decoder
	OnPeer        func(Peer) error
}

type TCPTransport struct {
	TCPTransportOpts
	listener net.Listener
	rpcch    chan RPC

	mu sync.RWMutex
	// peers map[net.Addr]Peer
}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
		rpcch:            make(chan RPC),
	}
}

// consume implements the Transport interface, which will return read-only channel
// fro reading the incoming messages received from another peer in the network
func (t *TCPTransport) Consume() <-chan RPC {
	return t.rpcch
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error

	t.listener, err = net.Listen("tcp", t.ListenAddr)
	if err != nil {
		return err
	}

	go t.startAcceptLoop()

	// select {}

	return nil

}

func (t *TCPTransport) startAcceptLoop() {

	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Printf("TCP accept error: %s\n", err)

		}

		go t.handleConn(conn)

	}

}

// close implements the Peer interface
func (p *TCPPeer) Close() error {
	return p.conn.Close()
}

func (t *TCPTransport) handleConn(conn net.Conn) {
	var err error

	defer func() {
		fmt.Printf("dropping peer connection: %s", err)
		conn.Close()

	}()

	fmt.Printf("new incoming connection %+v\n", conn)
	peer := NewTCPPeer(conn, true)

	if err = t.HandShakeFunc(peer); err != nil {
		fmt.Printf("handshake failed: %s\n", err)
		return

	}

	if t.OnPeer != nil {
		if err = t.OnPeer(peer); err != nil {
			fmt.Printf("error in  onpeer method %s", err)
			peer.Close()
			return
		}
	}

	// Read Loop
	rpc := RPC{}
	for {
		if err := t.Decoder.Decode(conn, &rpc); err != nil {
			fmt.Printf("TCP error: %s\n", err)
			continue
		}

		rpc.From = conn.RemoteAddr()
		t.rpcch <- rpc

		fmt.Printf("message: %+v\n", rpc)
	}
}
