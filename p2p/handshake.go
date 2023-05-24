package p2p

import (
	"errors"
)

var ErrInvalidHandShake = errors.New("invalid handshake")

type HandShakeFunc func(Peer) error

func NOPHandshakeFunc(Peer) error {

	return nil
}
