package p2p

type RPC struct {
	From string
	// Encoding of Message type
	Payload []byte
}

type Message struct {
	Payload any
}



// Server facing messages
type TestMessage struct {
	Message    string
	TestNumber int
}

type Heartbeat struct{}

type InitRequest struct {}

type InitResponse struct {
	HashIndex int
}

type InitAlert struct {
	NewPeerAddress string
}

type PeerClose struct {
	// Peer that is closing
	ClosingPeer string
}



// Client facing messages

type Set struct {
	Key   string
	Value string
}

type Get struct {
	Key string
}
