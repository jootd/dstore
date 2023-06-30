package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/jootd/dstore/p2p"
)

type FileServerOpts struct {
	StorageRoot       string
	PathTransformFunc PathTransformFunc
	Transport         p2p.Transport
	TCPTransportOpts  p2p.TCPTransportOpts
	BootstrapNodes    []string
}

type FileServer struct {
	FileServerOpts
	store  *Store
	quitch chan struct{}

	peerLock sync.Mutex
	peers    map[string]p2p.Peer
}

func NewFileServer(opts FileServerOpts) *FileServer {
	storeOpts := StoreOpts{
		Root:              opts.StorageRoot,
		PathTransformFunc: opts.PathTransformFunc,
	}
	return &FileServer{
		FileServerOpts: opts,
		store:          NewStore(storeOpts),
		quitch:         make(chan struct{}),
		peers:          make(map[string]p2p.Peer),
	}
}

func (s *FileServer) Start() error {

	if err := s.Transport.ListenAndAccept(); err != nil {
		return err
	}

	if len(s.BootstrapNodes) != 0 {
		s.BootstrapNetwork()
	}

	s.loop()

	return nil
}

func (s *FileServer) BootstrapNetwork() error {
	for _, addr := range s.BootstrapNodes {
		go func(addr string) {
			fmt.Printf("Attempting to connect %s\n", addr)
			if err := s.Transport.Dial(addr); err != nil {
				log.Println(err.Error())
			}

		}(addr)

	}

	return nil

}

type Payload struct {
	Key  string
	Data []byte
}

func (s *FileServer) broadcast(p *Payload) error {
	peers := []io.Writer{}

	for _, peer := range s.peers {
		peers = append(peers, peer)
	}

	mw := io.MultiWriter(peers...)
	return gob.NewEncoder(mw).Encode(p)
}

func (s *FileServer) OnPeer(p p2p.Peer) error {
	s.peerLock.Lock()
	defer s.peerLock.Unlock()
	s.peers[p.RemoteAddr().String()] = p
	log.Printf("connected with remote %s ", p.RemoteAddr())

	return nil
}

func (s *FileServer) StoreData(key string, r io.Reader) error {

	if err := s.store.Write(key, r); err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	_, err := io.Copy(buf, r)

	if err != nil {
		return err
	}

	p := &Payload{
		Key:  key,
		Data: buf.Bytes(),
	}

	return s.broadcast(p)

}

func (s *FileServer) Stop() {
	close(s.quitch)
}

func (s *FileServer) loop() {
	defer func() {
		fmt.Println("File Server Stopped")
		s.Transport.Close()
	}()

	for {
		select {
		case msg := <-s.Transport.Consume():
			var p Payload
			if err := gob.NewDecoder(bytes.NewReader(msg.Payload)).Decode(&p); err != nil {
				log.Fatal(err)
			}

			fmt.Printf("receivg msg: %+v", p)
		case <-s.quitch:
			fmt.Println("gracefull shutdown")
			return
		}

	}

}
