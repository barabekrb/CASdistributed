package main

import (
	"log"

	"github.com/barabekrb/simple_blockchain/p2p"
)

func main() {
	tcpOpts := p2p.TCPTransportOptions{
		ListenAddr:    ":3000",
		HandshakeFunc: p2p.NOPHandhsakeFunc,
		Decoder:       p2p.DefaultDecoder{},
	}

	tr := p2p.NewTCPTransport(tcpOpts)

	if er := tr.ListenAndAccept(); er != nil {
		log.Fatal(er)
	}

	select {}
}
