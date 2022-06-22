package verify

import (
	"context"
	"errors"
	"sync"
)

type GameResult struct {
	WinID  string
	LossID string
	RptID  string
}

type gameReport struct {
	ctx       context.Context
	ctxCancel context.CancelFunc
	winID     string
	lossID    string
}

var pendResults = struct {
	sync.RWMutex
	m map[string]gameReport
}{m: make(map[string]gameReport)}

var ErrorTimeout = errors.New("timeout")

// VerifyMatch verifies a match between two users.
// If the match is verified, it returns true.
func VerifyMatch(gr GameResult) (bool, error) {

	// Stage 1: Check if the match is already in the pending results.
	// 			If it is, the match can be immediately verified; return true.
	// 			If it is not, the match must be verified by the second player.
	// 			If the second player does not respond in time, move to stage 2.

	// Verify the match with using both players (i.e. diVerify)
	ok, err := diVerify(gr)
	if ok {
		return true, err
	}

	// Stage 2 (TODO): The winner must prove by themselves that the match is valid

	// Verify the match with one player (i.e. monoVerify)
	ok, err = monoVerify(gr)
	return ok, err

}
