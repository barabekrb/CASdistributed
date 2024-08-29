package p2p

import (
	"fmt"
	"net"
	"sync"
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

type TCPTransportOptions struct {
	ListenAddr    string
	HandshakeFunc HandshakeFunc
	Decoder       Decoder
}

type TCPTransport struct {
	TCPTransportOptions
	listener net.Listener

	mu    sync.RWMutex
	peers map[net.Addr]Peer
}

func NewTCPTransport(ops TCPTransportOptions) *TCPTransport {
	return &TCPTransport{
		TCPTransportOptions: ops,
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
	peer := NewTCPPeer(conn, true)

	if er := t.HandshakeFunc(peer); er != nil {
		conn.Close()
		fmt.Printf("TCP handshake error: %s\n", er)
		return
	}
	// Read loop
	msg := &Message{}
	// buf := make([]byte, 2000)
	for {
		if er := t.Decoder.Decode(conn, msg); er != nil {
			fmt.Printf("TCP error: %s/n", er)
			continue
		}

		msg.From = conn.RemoteAddr()

		fmt.Printf("message: %+v\n", msg)
	}

}
