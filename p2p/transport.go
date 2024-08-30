package p2p

// Peer is an interface that represents thr remote node
type Peer interface {
	Close() error
}

// Transport is anithing that handels the communicstion between nodes in the network
//This can be of the form (TCP, UDP, websockets...)
type Transport interface {
	ListenAndAccept() error
	Consume() <-chan RPC
}
