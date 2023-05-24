package p2p

import "net"

type Decoder interface {
	Decode(net.Conn, Temp) error
}
