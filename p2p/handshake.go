package p2p

// Handshakefunc ...
type HandshakeFunc func(Peer) error

func NOPHandhsakeFunc(Peer) error { return nil }
