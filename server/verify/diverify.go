package verify

import (
	"context"
	"fmt"
	"time"
)

// Verify a match with data from both players.
func diVerify(gr GameResult) (bool, error) {

	// Add this game report to the pending reports
	// 		Include a context with cancel function in the report
	// 		The cancel function will be called if the game report is confirmed by the second player
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	pendResults.Lock()
	pendResults.m[gr.RptID] = gameReport{ctx, cancel, gr.WinID, gr.LossID}
	pendResults.Unlock()

	// Remove this game result from the pending results before leaving this function
	defer removeResult(gr.RptID)

	// If our report matches a pending report, then we return true and
	if ok, _ := reportMatch(gr.WinID, gr.LossID, gr.RptID); ok {
		return true, nil
	}

	select {
	case <-ctx.Done():
		return false, nil
	case <-time.After(10 * time.Second):
		return false, ErrorTimeout
	}
}

func removeResult(rptID string) error {
	pendResults.Lock()
	defer pendResults.Unlock()
	fmt.Println(pendResults.m)
	fmt.Println(pendResults.m[rptID].ctx)
	delete(pendResults.m, rptID)
	return nil
}

func reportMatch(winID string, lossID string, rptID string) (bool, error) {
	for key, element := range pendResults.m {
		if key != rptID && element.winID == winID && element.lossID == lossID {

			pendResults.Lock()
			defer pendResults.Unlock()

			element.ctxCancel()
			pendResults.m[rptID].ctxCancel()

			return true, nil
		}
	}
	return false, nil
}
