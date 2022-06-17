package game

import (
	"bufio"
	"strings"
	"testing"
)

func TestHotseatGame(t *testing.T) {

	g, err := InitHotseat()
	if err != nil {
		t.Error(err)
	}
	g.reader = bufio.NewReader(strings.NewReader("f2f3\ne7e5\ng2g4\nd8h4\nq\n"))
	g.PlayHotseat()

	g, err = InitHotseat()
	if err != nil {
		t.Error(err)
	}
	g.reader = bufio.NewReader(strings.NewReader("q\n"))
	g.PlayHotseat()

}

func TestP2pGame(t *testing.T) {

	rch, wch := make(chan string, 1), make(chan string, 1)
	g, err := InitP2P(P2PParams{
		YouStart:  false,
		ReadChan:  rch,
		WriteChan: wch})
	if err != nil {
		panic(err)
	}
	g.reader = bufio.NewReader(strings.NewReader("q\n"))

	go func() {
		rch <- "f2f3"
		<-wch
		rch <- "g2g4"
		<-wch
	}()

	g.PlayP2P()
}
