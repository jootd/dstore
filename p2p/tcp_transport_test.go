package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPTransport(t *testing.T) {
	addr := ":4000"
	tr := NewTCPTransport(TCPTransportOpts{ListenAddr: addr})
	assert.Equal(t, tr.ListenAddr, addr)

	assert.Nil(t, tr.ListenAndAccept())

	select {}

}
