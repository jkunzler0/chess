package p2p

import (
	"bufio"
	"context"
	"crypto/rand"
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

type GameHello struct {
	Rw    *bufio.ReadWriter
	White bool
}

type P2pConfig struct {
	GroupID    string
	ProtocolID string
	ListenHost string
	ListenPort int
}

func P2pSetup(cfg *P2pConfig, ghn chan<- *GameHello) error {

	ghNotifier = ghn
	fmt.Printf("[*] Listening on: %s with port: %d\n", cfg.ListenHost, cfg.ListenPort)

	ctx := context.Background()
	r := rand.Reader

	// Create a new RSA key pair for this host
	prvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		panic(err)
	}

	// 0.0.0.0 will listen on any interface device
	sourceMultiAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d", cfg.ListenHost, cfg.ListenPort))

	// Construct a new libp2p Host
	host, err := libp2p.New(
		libp2p.ListenAddrs(sourceMultiAddr),
		libp2p.Identity(prvKey),
	)
	if err != nil {
		panic(err)
	}

	// Set a stream handler that will be called when another peer initiates a connection with this peer
	host.SetStreamHandler(protocol.ID(cfg.ProtocolID), handleStream)

	fmt.Printf("\n[*] Your Multiaddress Is: /ip4/%s/tcp/%v/p2p/%s\n", cfg.ListenHost, cfg.ListenPort, host.ID().Pretty())
	fmt.Println(host.ID())

	// Setup MDNS to discover other peers in the network
	peerChan := initMDNS(host, cfg.GroupID)
	// Block here until we discover a peer
	peer, ok := <-peerChan
	if !ok {
		panic("No peers found")
	}

	// If hosting return to main and wait for a peer
	if peer.ID == host.ID() {
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

	// Create a buffer stream for non blocking read and write
	ghNotifier <- &GameHello{
		bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream)),
		false,
	}
	// p2pGame(rw, Black)
	fmt.Println("End of p2p")
	return nil

	// Wait here for now
	// select {}
}

// #######################################################################
// (Section 2) Read/Write ################################################
// #######################################################################

func handleStream(stream network.Stream) {
	fmt.Println("Got a new stream!")
	// Create a buffer stream for non blocking read and write
	ghNotifier <- &GameHello{
		bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream)),
		true,
	}
	// p2pGame(rw, White)
	// fmt.Println(stream)
}

func ReadStream(rw *bufio.ReadWriter) string {
	fmt.Println("Waiting for opponent...")
	// ReadString will block until the delimiter is entered
	// We expect a correctly formated input since they already processed their own move
	// 		So if its not a valid input, just panic for now
	//		TODO can be to ask them again for a valid input
	move, err := rw.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading from buffer")
		panic(err)
	}
	if move == "" || move == "\n" {
		fmt.Println("Empty buffer")
		panic(err)
	}
	// Remove the delimeter from the string
	move = strings.TrimSuffix(move, "\n")
	// move = strings.ReplaceAll(move, " ", "")
	fmt.Println("Their move: ", move)
	return move
}

func WriteStream(rw *bufio.ReadWriter, move string) {
	// Write to stream
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
