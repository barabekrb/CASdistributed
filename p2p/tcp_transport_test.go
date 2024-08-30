package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPTransport(t *testing.T) {
	tcpTrOPt := TCPTransportOptions{ListenAddr: ":4000"}
	listenAddr := ":4000"
	tr := NewTCPTransport(tcpTrOPt)

	assert.Equal(t, tr.ListenAddr, listenAddr)

	// Server
	assert.Nil(t, tr.ListenAndAccept())
}
