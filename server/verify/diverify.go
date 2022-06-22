package verify

import (
	"context"
	"fmt"
	"time"
)

// Verify a match with data from both players.
func diVerify(gr GameResult) (bool, error) {

	// Step 1: Check if our match is already in the pending game results.
	if ok, _ := reportMatch(gr.WinID, gr.LossID, gr.RptID); ok {
		return true, nil
	}

	// Step 2: Add our match to the pending results.
	// 		Include a context with cancel function in the report
	// 		The cancel function will be called if the game report is verified by the second player
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pendReports.Lock()
	pendReports.m[gr.RptID] = gameReport{ctx, cancel, gr.WinID, gr.LossID}
	pendReports.Unlock()

	// Remove this game result from the pending results before leaving this function
	defer removeReport(gr.RptID)

	// Step 3: Wait for the other player to complete step 1.
	select {
	case <-ctx.Done():
		// If the context is canceled,
		//		then the match will be verified by the other player
		return false, nil
	case <-time.After(10 * time.Second):
		// If the context is not canceled after some time,
		//		then return and attempt to another method for verification
		return false, ErrorTimeout
	}
}

func removeReport(rptID string) error {
	pendReports.Lock()
	defer pendReports.Unlock()
	fmt.Println(pendReports.m)
	fmt.Println(pendReports.m[rptID].ctx)
	delete(pendReports.m, rptID)
	return nil
}

func reportMatch(winID string, lossID string, rptID string) (bool, error) {
	for key, element := range pendReports.m {
		if key != rptID && element.winID == winID && element.lossID == lossID {

			// Cancel the context of the pending report
			//		to signal that the match will be verified
			pendReports.Lock()
			defer pendReports.Unlock()
			element.ctxCancel()

			return true, nil
		}
	}
	return false, nil
}
