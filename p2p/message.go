package p2p

import "net"

// Message hold any arbitrated data that is being
// sent over the each transport between two nodes in
// network
type RPC struct {
	From    net.Addr
	Payload []byte
}
