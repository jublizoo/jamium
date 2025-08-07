package p2p

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"net"
	"strconv"
)

type Peer struct {
	net.Conn
	isReplica bool

	// Which interval in the unit circle does this peer use for consistent hashing? (0-indexed)
	hashInterval int
}

func (p *Peer) getRemoteAddr() string {
	return p.Conn.RemoteAddr().String()
}

func (p *Peer) Write(b []byte) (int, error) {
	return p.Conn.Write(b)
}

func (p *Peer) SendRPC(from string, payload any) error {
	msg := Message{
		Payload: payload,	
	}		
	msg_buf := new(bytes.Buffer)
	if err := gob.NewEncoder(msg_buf).Encode(msg); err != nil {
		fmt.Println("Failed to encode RPC Message:", msg)
		return err
	}
	rpc := RPC {
		From: from, 
		Payload: msg_buf.Bytes(),
	}
	rpc_buf := new(bytes.Buffer)
	if err := gob.NewEncoder(rpc_buf).Encode(rpc); err != nil {
		fmt.Println("Failed to encode RPC struct:", rpc)
		return err
	}
	_, err := p.Write(rpc_buf.Bytes())

	return err
}

func (p *Peer) handleDisconnect(err error) {
	return 
}

type KVCluster struct {
	localAddr *string
	peers   map[string](*Peer)
	rpcch   chan *RPC
	closech chan struct{}
}

func createKVCluster() (*KVCluster, error) {
	localAddr, err := getLocalAddr()
	if err != nil {
		return nil, err
	}

	cluster := KVCluster{
		localAddr: &localAddr,
		peers:   make(map[string](*Peer)),
		rpcch:   make(chan *RPC),
		closech: make(chan struct{}),
	}

	return &cluster, nil
}

func (k *KVCluster) getNumPeers() int {
	return len(k.peers)
}

func (k *KVCluster) getLocalAddr() (string, error) {
	go listenToSelf()
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		return "", err
	}
	return conn.LocalAddr().String(), nil
}

func listenToSelf() error {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		return err
	}
	_, err = listener.Accept()
	if err != nil {
		return err
	}

	return nil
}


// TODO: Get remote address
func (k *KVCluster) sendHeartbeat(peer *Peer) {
	localAddr := "PLACEHOLDER"
	peer.SendRPC(localAddr, Heartbeat{})
}
 
// TODO: Does not wrap in RPC currently
func (k *KVCluster) sendHeartbeatOld(peer *Peer) {
	buf := new(bytes.Buffer)
	gob.NewEncoder(buf).Encode(Heartbeat{})
	peer.Write(buf.Bytes())
}

func (k *KVCluster) ListenAndAccept(port int) error {
	portStr := ":" + strconv.Itoa(port)
	listener, err := net.Listen("tcp", portStr)
	if err != nil {
		fmt.Println("Error when listening on port", portStr, ":", err)
		return err
	}
	go k.acceptNewConnections(listener)

	return nil
}

func (k *KVCluster) acceptNewConnections(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error while accepting peer:", err)
			return
		}
		go k.handleConnection(conn)
	}
}

func (k *KVCluster) Dial(address string) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	k.handleConnection(conn)

	return nil
}

func (k *KVCluster) handleConnection(conn net.Conn) {
	peer := Peer{
		Conn: conn,
	}

	k.peers[conn.RemoteAddr().String()] = &peer

	for {
		var buf []byte
		_, err := conn.Read(buf)
		if err != nil {
			fmt.Println(err)
			peer.handleDisconnect(err)
			return
		}

		rpc := RPC{}
		decoder := gob.NewDecoder(bytes.NewReader(buf))
		decoder.Decode(&rpc)
		rpc.From = conn.RemoteAddr().String()
		k.rpcch <- &rpc
	}
}

func (k *KVCluster) RPCLoop() {
	for {
		select {
		case rpc := <- k.rpcch:
			var msg Message
			decoder := gob.NewDecoder(bytes.NewReader(rpc.Payload))
			decoder.Decode(&msg)
			handleMessage(msg)		
		case <- k.closech:
			// TODO: alert peers of closing? Alternatively they could auto-detect
			return
		}
	}
}

func handleMessage(msg Message) error {
	switch m := msg.Payload.(type) {
	case TestMessage:
		fmt.Println("Received test message with message", m.Message, "and test number", m.TestNumber)
	default:
		err := errors.New("Message not recognized")
		fmt.Println("Error while handling message:", err)
		return err
	}

	return nil
}

// Request for KV
func (k * KVCluster) handleClientRequest(key string) error {
	return nil
}

func (k *KVCluster) onPeerClose() {
	
}

func (k *KVCluster) Close() {
	k.closech <- struct{}{}
}

func (k *KVCluster) localEventLoop() {
	for {
		select {
		case <-k.closech:
			k.Close()
			return
		}
	}
}

