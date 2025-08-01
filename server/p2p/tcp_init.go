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
func (k *KVCluster) handlePeerInitRequest(fromPeer string, request InitRequest) error {
	alert := InitAlert {
		newPeerAddress: fromPeer,
	}
	buf := new(bytes.Buffer)
	gob.NewEncoder(buf).Encode(alert)
	for _, peer := range k.peers {
		err := peer.SendRPC(InitAlert)
		if err != nil {
			response := InitResponse {
				success: false,
			}	
			err := peer.SendRPC(response)
			/* 	
				TODO: Alert other peers of this failure.
				They must remove this as a peer, and this may potentially trigger
				a key redistribution (if one of the other peers took on another init requet)
				A better solution: wait until peer is fully initialized to initialize
				any other peers
			 */
		}
		// TODO: Wait for response of Peers when 
	}
	response := InitResponse {
		HashIndex: numPeers,
	}

	newPeer := k.peers[fromPeer]
	err := newPeer.SendRPC(response)
	if err != nil {
		fmt.Println("Error while responding to Peer init request:", err)
		return err
	}

	return nil
}
