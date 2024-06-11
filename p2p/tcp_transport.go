package p2p

import (
	"fmt"
	"net"
)

// TCPPeer represents the remote node over a TCP established connection.
type TCPPeer struct {
	// conn is the underlying connection of the peer
	conn net.Conn

	// If we dial and retrieve a conn => outbound = true
	// If we accept and retrieve a conn => outbound = false
	outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

// Close implements the close interface
func (p *TCPPeer) Close() error {
	return p.conn.Close()
}

type TCPTransportOpts struct {
	ListenAddress string
	ShakeHands    HandshakeFunc
	Decoder       Decoder
	OnPeer        func(Peer) error
}

type TCPTransport struct {
	TCPTransportOpts
	listener net.Listener

	// IMPORTANT:
	// The chan keyword in Go is used to create a channel,
	// which is a data structure that allows goroutines to
	// communicate and synchronize their execution.
	rpcch chan RPC
}

// Anything in the module that starts with a uppercase letter
// is public and can be seen and used with import in other modules
func NewTCPTransport(t TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: t,
	}
}

// Consume implenets the Transport interface, which will return read-only
// channels for reading the incoming messages received from another peer
// in the network.
func (t *TCPTransport) Consume() <-chan RPC {
	return t.rpcch
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error

	t.listener, err = net.Listen("tcp", t.ListenAddress)
	if err != nil {
		return err
	}

	go t.startAcceptLoop()

	return nil
}

func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Printf("TCP accept error: %s\n", err)
		}
		go t.handleConn(conn)
	}
}

func (t *TCPTransport) handleConn(conn net.Conn) {

	var err error

	defer func() {
		fmt.Printf("Dropping peer connection: %s\n", err)
		conn.Close()
	}()

	peer := NewTCPPeer(conn, true)

	if err := t.ShakeHands(peer); err != nil {
		fmt.Printf("TCP handshake error: %s\n", err)
		conn.Close()
		return
	}

	if t.OnPeer != nil {
		if err := t.OnPeer(peer); err != nil {
			return
		}
	}

	//Read loop
	rpc := RPC{}

	// buf := make([]byte, 2000)
	for {
		// n, err := conn.Read(buf)
		// if err != nil {
		// 	fmt.Printf("TCP error: %s\n", err)
		// }

		if err := t.Decoder.Decode(conn, &rpc); err != nil {
			continue
		}

		rpc.From = conn.RemoteAddr()

		t.rpcch <- rpc
	}

}
