package p2p

import "fmt"

// The address param is the address to connect to, or nil if it is the first peer in the cluster
func addKVCluster(address *string) *KVCluster {
	return nil		
}

func createPeers() error {
	numPeers := 5

	for i := range numPeers {
		var leaderPeerAddress *string
		if i == 0 {
			leaderPeer := addKVCluster(nil)
			leaderPeerAddressStr, err := leaderPeer.getLocalAddr()
			if err != nil {
				return err
			}
			leaderPeerAddress := &leaderPeerAddressStr
			fmt.Println(leaderPeerAddress)
		} else {
			addKVCluster(leaderPeerAddress)
		}
	}

	return nil
}
