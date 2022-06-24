package p2p

import (
	"bufio"
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"strings"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	"github.com/multiformats/go-multiaddr"
)

// #######################################################################
// (Section 1) P2P Setup #################################################
// #######################################################################

var ghNotifier chan<- *GameHello
var gh GameHello

type GameHello struct {
	RCh   chan string // To read from the peer and write to the game thread
	WCh   chan string // To read from the game thread and write to the peer
	White bool        // True if we are white, false if we are black
}

type P2pConfig struct {
	GroupID    string
	ProtocolID string
	ListenHost string
	ListenPort int
}

func P2pSetup(cfg *P2pConfig, ghn chan<- *GameHello) error {

	ghNotifier = ghn

	// fmt.Printf("[*] Listening on: %s with port: %d\n", cfg.ListenHost, cfg.ListenPort)

	ctx := context.Background()
	r := rand.Reader

	// Create a new ECDSA key pair for this host
	xprv, _, err := crypto.GenerateECDSAKeyPair(r)
	if err != nil {
		panic(err)
	}

	// 0.0.0.0 will listen on any interface device
	sourceMultiAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d", cfg.ListenHost, cfg.ListenPort))

	// Construct a new libp2p Host
	host, err := libp2p.New(
		libp2p.ListenAddrs(sourceMultiAddr),
		libp2p.Identity(xprv),
	)
	if err != nil {
		panic(err)
	}

	// Set a stream handler that will be called when another peer initiates a connection with this peer
	host.SetStreamHandler(protocol.ID(cfg.ProtocolID), handleStream)

	// fmt.Printf("\n[*] Your Multiaddress Is: /ip4/%s/tcp/%v/p2p/%s\n", cfg.ListenHost, cfg.ListenPort, host.ID().Pretty())

	// Setup MDNS to discover other peers in the network
	peerChan := initMDNS(host, cfg.GroupID)
	// Block here until we discover a peer
	peer, ok := <-peerChan
	if !ok {
		panic("No peers found")
	}

	// If hosting, return to main and wait for a peer
	// Host will play as white
	if peer.ID == host.ID() {
		fmt.Println("Waiting for a peer...")
		gh.White = true
		return nil
	}

	fmt.Printf("Found peer: %+v, connecting\n", peer)
	if err := host.Connect(ctx, peer); err != nil {
		fmt.Println("Connection failed:", err)
		// TODO: retry on error
	}

	// Open a stream, this stream will be handled by handleStream at the other end
	stream, err := host.NewStream(ctx, peer.ID, protocol.ID(cfg.ProtocolID))

	// If failed to open a stream to peer, assume we are white/first player
	if err != nil {
		fmt.Println("Stream open failed", err)
		return err
	}

	fmt.Println("Connected to:", peer)
	handleStream(stream)

	return nil
}

// #######################################################################
// (Section 2) Read/Write ################################################
// #######################################################################

func handleStream(stream network.Stream) {
	fmt.Println("Got a new stream!")

	// Create a buffer stream for non blocking read/write to/from the peer
	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))
	// Create channels to send/receive moves to/from the game thread
	gh.RCh, gh.WCh = make(chan string, 1), make(chan string, 1)

	// Kick off the read/write routines for communicating with the peer
	go readStream(rw, gh.RCh)
	go writeStream(rw, gh.WCh)

	// Pass back the read/write channels to the game thread
	ghNotifier <- &GameHello{gh.RCh, gh.WCh, gh.White}
}

var ErrorStreamReset = errors.New("stream reset")

// Read from the connected peer and send to rch
func readStream(rw *bufio.ReadWriter, ch chan<- string) {
	// We expect a correctly formated input since they already processed their own move
	// 		So if its not a valid input, just panic for now
	//		TODO can be to ask them again for a valid input
	for {
		// Block here and wait for peer
		move, err := rw.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from buffer")
			panic(err)
		}
		if move == "quit" || move == "q" {
			ch <- move
			return
		}
		if move == "" || move == "\n" {
			fmt.Println("Empty buffer")
			panic(err)
		}
		// Remove the delimeter from the string
		move = strings.TrimSuffix(move, "\n")
		// Send their move to the game / main thread
		ch <- move
	}
}

// Write to the connected peer from wch
func writeStream(rw *bufio.ReadWriter, ch <-chan string) {

	for {
		// Block here until our move is sent on ch
		move, ok := <-ch
		if !ok {
			// If the channel has been closed, exit this go routine
			return
		}
		_, err := rw.WriteString(fmt.Sprintf("%s\n", move))
		if err != nil {
			fmt.Println("Error writing to buffer")
			panic(err)
		}
		err = rw.Flush()
		if err != nil {
			fmt.Println("Error flushing buffer")
			panic(err)
		}
	}
}

// #######################################################################
// (Section 3) MDNS Setup ################################################
// #######################################################################

type discoveryNotifee struct {
	PeerChan chan peer.AddrInfo
}

// Interface to be called when new  peer is found
func (n *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	n.PeerChan <- pi
}

// Initialize the MDNS service
func initMDNS(peerhost host.Host, rendezvous string) chan peer.AddrInfo {
	// Register with service so that we get notified about peer discovery
	n := &discoveryNotifee{}
	n.PeerChan = make(chan peer.AddrInfo)

	// An hour might be a long long period in practical applications. But this is fine for us
	ser := mdns.NewMdnsService(peerhost, rendezvous, n)
	if err := ser.Start(); err != nil {
		panic(err)
	}
	return n.PeerChan
}
