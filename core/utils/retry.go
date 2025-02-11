package utils

import (
	"context"
	"time"
)

// Retry executes the provided function with retries on failure
func Retry(ctx context.Context, attempts int, sleep time.Duration, fn func() error) error {
	var err error
	for i := 0; i < attempts; i++ {
		err = fn()
		if err == nil {
			// successful execution
			return nil
		}
		select {
		case <-ctx.Done():
			// return if context is done
			return ctx.Err()
		case <-time.After(sleep):
			// exponential backoff
			sleep *= 2
		}
	}
	return err // return the last error
}
