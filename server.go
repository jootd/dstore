package main

import (
	"fmt"
	"log"

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
	}
}

func (s *FileServer) Start() error {

	if err := s.Transport.ListenAndAccept(); err != nil {
		return err
	}

	if !(len(s.BootstrapNodes) == 0) {
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
			fmt.Println(string(msg.Payload))

		case <-s.quitch:
			fmt.Println("gracefull shutdown")
			return
		}

	}

}
