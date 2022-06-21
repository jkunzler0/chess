package process

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type gameResult struct {
	ctx       context.Context
	ctxCancel context.CancelFunc
	winID     string
	lossID    string
	// reporterID string
}

var pendResults = struct {
	sync.RWMutex
	m map[string]gameResult
}{m: make(map[string]gameResult)}

func ProcessGameResult(winID string, lossID string, rptID string) (bool, error) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pendResults.Lock()
	pendResults.m[rptID] = gameResult{ctx, cancel, winID, lossID}
	pendResults.Unlock()

	defer removeResult(rptID)

	if ok, _ := confirmedResult(winID, lossID, rptID); ok {
		return true, nil
	}

	select {
	case <-ctx.Done():
		return true, nil
	case <-time.After(10 * time.Second):
		return false, nil
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

func confirmedResult(winID string, lossID string, rptID string) (bool, error) {
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
