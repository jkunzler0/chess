package game

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"testing"
)

func TestHotseat(t *testing.T) {

	r, w, err := os.Pipe()
	if err != nil {
		log.Fatal(err)
	}

	var move string
	reader := bufio.NewReader(r)

	// Test hotseat with fool's mate
	move = "f2f3\ne7e5\ng2g4\nd8h4\n"
	w.WriteString(fmt.Sprintf("%s\n", move))
	HotseatGame(reader)

	// Test hotseat with quiting
	move = "q\n"
	w.WriteString(fmt.Sprintf("%s\n", move))
	HotseatGame(reader)

	// var stdin bytes.Buffer
	// stdin.Write([]byte("f2f3\ne7e5\ng2g4\nd8h4\nq\n"))
	// stdin.Write([]byte("q\n"))
	// stdin.Write([]byte(""))
	// reader := bufio.NewReader(&stdin)

}

func TestP2pGame(t *testing.T) {

	// P2pGame(rch <-chan string, wch chan<- string, yourColor bool, stdin *bufio.Reader) {

	rch, wch := make(chan string, 1), make(chan string, 1)

	r, w, err := os.Pipe()
	if err != nil {
		log.Fatal(err)
	}

	var move string
	reader := bufio.NewReader(r)

	// Test p2p with fool's mate
	move = "e7e5\nd8h4\n"
	w.WriteString(fmt.Sprintf("%s\n", move))

	go func() {
		rch <- "f2f3"
		<-wch
		rch <- "g2g4"
		<-wch
	}()

	P2pGame(rch, wch, false, reader)
}
