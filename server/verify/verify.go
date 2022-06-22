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

var pendReports = struct {
	sync.RWMutex
	m map[string]gameReport
}{m: make(map[string]gameReport)}

var ErrorTimeout = errors.New("timeout")

// VerifyMatch verifies a match between two users.
// If the match is verified, it returns true.
func VerifyMatch(gr GameResult) (bool, error) {

	// STAGE 1: Verify the match with using both players (i.e. diVerify).
	//
	// Step 1: Check if our match is already in the pending game results.
	// 			If it is, the match can be immediately verified; return true.
	// Step 2: Add our match to the pending results.
	// Step 3: Wait for the other player to complete step 1.
	// 			If the other player completes step 1, we can return.
	// 			If we wait more than some time, move to stage 2.

	ok, err := diVerify(gr)
	if ok {
		return true, err
	}

	// STAGE 2 (TODO): Verify the match with one player (i.e. monoVerify).
	//
	// Step 1: ...

	ok, _ = monoVerify(gr)
	return ok, err

}
