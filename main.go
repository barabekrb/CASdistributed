package main

import (
	"fmt"
	"log"

	"github.com/barabekrb/simple_blockchain/p2p"
)

func OnPeer(peer p2p.Peer) error {
	// fmt.Println("doing some logic with the peer outside of TCPTransport")
	peer.Close()
	return nil
}
func main() {
	tcpOpts := p2p.TCPTransportOptions{
		ListenAddr:    ":3000",
		HandshakeFunc: p2p.NOPHandhsakeFunc,
		Decoder:       p2p.DefaultDecoder{},
		OnPeer:        OnPeer,
	}

	tr := p2p.NewTCPTransport(tcpOpts)

	go func() {
		for {
			msg := <-tr.Consume()
			fmt.Printf("%+v\n", msg)
		}
	}()

	if er := tr.ListenAndAccept(); er != nil {
		log.Fatal(er)
	}

	select {}
}
