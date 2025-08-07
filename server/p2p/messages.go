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

type InitRequest struct{}

// Response to initialization request of new peer.
// When Peer A send an InitRequest Peer B to join the cluster,
// Peer B sends back and InitResponse.
type InitResponse struct {
	Success   bool
	HashIndex int
}

type InitAlert struct {
	NewPeerAddress string
}

// TODO: Different type of alert for a peer closing  
// (intentional, couldn't initialize, multiple missed requests)
// than for a peer that is temporarily down (single missed request).
// If a peer is temporarily down, we switch to the replica
// For now, we assume the peer will either recover, or the information
// will only exist on the replica (no other peer in the cluster will duplicate the replica)

type PeerClose struct {
	// Peer that is closing
	ClosingPeerAddress string
}



// Client facing messages

type Set struct {
	Key   string
	Value string
}

type Get struct {
	Key string
}
