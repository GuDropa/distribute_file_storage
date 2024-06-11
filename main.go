package main

import (
	"fmt"
	"log"

	"github.com/GuDropa/distribute_file_storage/p2p"
)

func OnPeer(peer p2p.Peer) error {
	fmt.Println("doing some logic with the peer outside of TCPTransport")
	return nil
}

func main() {

	tcpOpts := p2p.TCPTransportOpts{
		ListenAddress: ":3000",
		Decoder:       p2p.DefaultDecoder{},
		ShakeHands:    p2p.NOPHandshakeFunc,
		OnPeer:        OnPeer,
	}

	tr := p2p.NewTCPTransport(tcpOpts)

	fmt.Printf("Running on port %v\n", tr.ListenAddress)

	// IMPORTANT:
	// The go statement in Go is used to launch a new goroutine.
	// A goroutine is a lightweight thread managed by the Go runtime.
	// When you use the go statement, the function call that follows
	// it is executed concurrently in a separate goroutine.
	go func() {
		for {
			msg := <-tr.Consume()
			fmt.Printf("%v\n", msg)
		}
	}()

	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}

	select {}
}
