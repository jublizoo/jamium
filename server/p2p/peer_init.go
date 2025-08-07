package p2p

import (
	"fmt"
	"net"
)

// When a peer alerts of a new peer joining
func (k * KVCluster) handlePeerInitAlert(alert InitAlert) error {
	peerAddr := alert.NewPeerAddress
	conn, err := net.Dial("tcp", peerAddr)
	if err != nil {
		fmt.Println("Could not connect to new peer (through alert)")
		return err
	}

	numPeers := len(k.peers)
	peer := Peer {
		Conn: conn,
		hashInterval: numPeers,
	}
	k.peers[peerAddr] = &peer

	return nil
}

func (k *KVCluster) sendPeerRemoval(peerAddress string) error {
	removalRequest := RemovePeer {
		PeerAddress: peerAddress,	
	}

	for _, peer := range k.peers {
		err := peer.SendRPC(*k.localAddr, removalRequest)
		if err != nil {
			fmt.Println("Failed to send removal request to peer", 
						peer.Conn.RemoteAddr().String, 
						"of peer", peerAddress)
			return err
		}
	}

	return nil
}

// When a peer requests to join the cluster
func (k *KVCluster) handlePeerInitRequest(fromPeer string, request InitRequest) error {
	alert := InitAlert {
		NewPeerAddress: fromPeer,
	}

	// List of peers that successfully received init request 
	alertedPeers := make([](*Peer), 0)
	for _, peer := range k.peers {
		if peer.Conn.RemoteAddr().String() == fromPeer {
			continue
		}
		err := peer.SendRPC(*k.localAddr, alert)
		if err != nil {
			response := InitResponse {
				Success: false,
			}	
			err := peer.SendRPC(*k.localAddr, response)
			if err != nil {
				return err
			}
			/* 	
				TODO: Alert other peers of this failure.
				They must remove this as a peer, and this may potentially trigger
				a key redistribution (if one of the other peers took on another init requet)
				A better solution: wait until peer is fully initialized to initialize
				any other peers
			 */
			for _, peer := range alertedPeers {
				k.sendPeerRemoval(peer.getRemoteAddr())
			}
			return err

				
		}

		alertedPeers = append(alertedPeers, peer)	

		// TODO: Wait for response of Peers when 
	}

	numPeers := k.getNumPeers()
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
