package p2p

import (
	"errors"
	"fmt"
	"log"
	"net"
)

type TCPPeer struct {

	// conn is underlying connection of peer
	net.Conn

	// dial, retrieve a conn => outbound == true
	// accept,  retrieve a conn => outbound == false
	outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		Conn:     conn,
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

	// mu sync.RWMutex
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

	log.Printf("TCP Transport listening on port %s\n", t.ListenAddr)

	return nil

}

func (t *TCPTransport) startAcceptLoop() {

	for {
		conn, err := t.listener.Accept()
		if errors.Is(err, net.ErrClosed) {
			return
		}
		if err != nil {
			fmt.Printf("TCP accept error: %s\n", err)
		}

		fmt.Println("TCP ACCEPT")

		go t.handleConn(conn, false)

	}

}

func (t *TCPTransport) Close() error {
	return t.Close()
}

func (t *TCPTransport) Dial(addr string) error {

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	go t.handleConn(conn, true)

	return nil

}

// send implements the Peer interface
func (p *TCPPeer) Send(msg []byte) error {

	if _, err := p.Conn.Write(msg); err != nil {
		return err
	}

	return nil
}

// close implements the Peer interface

func (t *TCPTransport) handleConn(conn net.Conn, outbound bool) {
	var err error

	defer func() {
		fmt.Printf("dropping peer connection: %s", err)
		conn.Close()
	}()

	peer := NewTCPPeer(conn, outbound)

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
		err = t.Decoder.Decode(conn, &rpc)
		if err != nil {
			return
		}
		// fmt.Printf("TCP error: %s\n", err)

		rpc.From = conn.RemoteAddr()
		t.rpcch <- rpc

		fmt.Printf("message: %+v\n", rpc)
	}
}
