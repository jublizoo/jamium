package p2p

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
)

// When a peer alerts of a new peer joining
func (k * KVCluster) handlePeerInitAlert(alert InitAlert) error {
	peerAddr := alert.newPeerAddress
	conn, err := net.Dial("tcp", peerAddr)
	if err != nil {
		fmt.Println("Could not connect to new peer (through alert)")
		return err
	}

	numPeers = len(k.peers)
	peer := Peer {
		Conn: conn,
		hashInterval: numPeers,
	}
	k.peers[peerAddr] = peer

	return nil
}

// When a peer requests to join the cluster
func (k *KVCluster) handlePeerInitRequest(fromPeer string, request InitRequest) {
	alert := InitAlert {
		newPeerAddress: fromPeer,
	}
	buf := new(bytes.Buffer)
	gob.NewEncoder(buf).Encode(alert)
	for _, peer := range k.peers {
		peer.SendRPC(InitAlert)
	}
	response := InitResponse {
		HashIndex: numPeers,
	}

}


