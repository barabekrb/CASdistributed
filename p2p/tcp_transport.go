package p2p

import (
	"fmt"
	"net"
)

// TCPPeer represents the remote node over TCP established connection
type TCPPeer struct {
	// conn is the underlying connection of the peer
	connection net.Conn

	// if we dial and retrieve connection => outbound == true
	// if we accept and retrieve connection => outbound == false
	outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		connection: conn,
		outbound:   outbound,
	}
}

func (p *TCPPeer) Close() error {
	return p.connection.Close()
}

type TCPTransportOptions struct {
	ListenAddr    string
	HandshakeFunc HandshakeFunc
	Decoder       Decoder
	OnPeer        func(Peer) error
}

type TCPTransport struct {
	TCPTransportOptions
	listener net.Listener
	rpcch    chan RPC
}

// Consume emplements the Transport interface, which will return read-only chan
// for reading the incoming messages recieved from another peer in the network
func (t *TCPTransport) Consume() <-chan RPC {
	return t.rpcch
}

func NewTCPTransport(ops TCPTransportOptions) *TCPTransport {
	return &TCPTransport{
		TCPTransportOptions: ops,
		rpcch:               make(chan RPC),
	}
}

func (t *TCPTransport) ListenAndAccept() error {
	var er error
	t.listener, er = net.Listen("tcp", t.ListenAddr)
	if er != nil {
		return er
	}
	go t.startAcceptLoop()

	return nil
}

func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, er := t.listener.Accept()
		if er != nil {
			fmt.Printf("TCP accept error %s\n", er)
		}
		fmt.Printf("new incomming connection %+v\n", conn)

		go t.handelConn(conn)
	}
}

type Temp struct {
}

func (t *TCPTransport) handelConn(conn net.Conn) {
	var er error
	defer func() {
		fmt.Printf("dropping peer connection: %s", er)
		conn.Close()
	}()

	peer := NewTCPPeer(conn, true)

	if er := t.HandshakeFunc(peer); er != nil {
		return
	}

	if t.OnPeer != nil {
		if er := t.OnPeer(peer); er != nil {
			return
		}
	}

	// Read loop
	rpc := RPC{}
	// buf := make([]byte, 2000)
	for {
		er := t.Decoder.Decode(conn, &rpc)
		if er != nil {
			fmt.Printf("TCP reader error: %s/n", er)
			return
		}

		rpc.From = conn.RemoteAddr()
		t.rpcch <- rpc
		fmt.Printf("message: %+v\n", rpc)
	}

}
