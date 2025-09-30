package context

import (
	"context"
	"time"
)

func legacyCall(userID int64) (int64, error) {
	time.Sleep(time.Second * 5)
	return userID + 100, nil
}

func safeProcessLegacyCall(ctx context.Context) {
	resC := make(chan int64, 1)
	errC := make(chan error, 1)

	go func() {
		res, err := legacyCall(100)
		if err != nil {
			errC <- err
		}
		resC <- res
	}()

	select {
	case res := <-resC:
		println("result:", res)
	case err := <-errC:
		println("error:", err)
	case <-ctx.Done():
		println("timeout:", ctx.Err())
	}
}
