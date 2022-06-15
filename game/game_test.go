package game

import (
	"bufio"
	"strings"
	"testing"
)

func TestHotseat(t *testing.T) {

	gs := NewGameState()
	gs.reader = bufio.NewReader(strings.NewReader("f2f3\ne7e5\ng2g4\nd8h4\nq\n"))
	HotseatGame(gs)

	gs = NewGameState()
	gs.reader = bufio.NewReader(strings.NewReader("q\n"))
	HotseatGame(gs)

}

func TestP2pGame(t *testing.T) {

	gs := NewGameState()
	gs.reader = bufio.NewReader(strings.NewReader("e7e5\nd8h4\n"))
	gs.white = false
	gs.rch, gs.wch = make(chan string, 1), make(chan string, 1)

	go func() {
		gs.rch <- "f2f3"
		<-gs.wch
		gs.rch <- "g2g4"
		<-gs.wch
	}()

	P2pGame(gs)
}
