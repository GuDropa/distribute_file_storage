package p2p

// Peer is an interface that represents the remote node.
type Peer interface {
	Close() error
}

// Transport is anything that handles the communication
// between the nodes in the network. This can be of the
// form (TCP, UDP, websockets, ...)
type Transport interface {
	ListenAndAccept() error

	// IMPORTANT:
	// The "<-" prefix is user to dictate that your only
	// allowed to read from the channel created.
	Consume() <-chan RPC
}
