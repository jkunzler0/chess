package game

import (
	"bufio"
	"strings"
	"testing"
)

func TestHotseat(t *testing.T) {

	reader := bufio.NewReader(strings.NewReader("f2f3\ne7e5\ng2g4\nd8h4\nq\n"))
	HotseatGame(reader)

	reader = bufio.NewReader(strings.NewReader("q\n"))
	HotseatGame(reader)

}

func TestP2pGame(t *testing.T) {

	// P2pGame(rch <-chan string, wch chan<- string, yourColor bool, stdin *bufio.Reader) {

	rch, wch := make(chan string, 1), make(chan string, 1)

	reader := bufio.NewReader(strings.NewReader("e7e5\nd8h4\n"))

	go func() {
		rch <- "f2f3"
		<-wch
		rch <- "g2g4"
		<-wch
	}()

	P2pGame(rch, wch, false, reader)
}
