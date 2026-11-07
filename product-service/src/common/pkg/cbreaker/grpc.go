package cbreaker

import (
	"fmt"
	"product-service/src/common/log"
	"slices"
	"time"

	"github.com/sony/gobreaker/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewGrpc(name string) *gobreaker.CircuitBreaker[any] {
	setting := gobreaker.Settings{
		Name:        name,
		MaxRequests: 5,
		Interval:    1 * time.Minute,
		Timeout:     15 * time.Minute,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRate := float32(counts.TotalFailures) / float32(counts.Requests)
			return failureRate >= 0.8 && counts.Requests >= 5
		},
		IsSuccessful: func(err error) bool {
			if err == nil {
				return true
			}
			st, ok := status.FromError(err)
			if !ok {
				return false
			}
			statusCodeFail := []codes.Code{
				codes.Unavailable,
				codes.DeadlineExceeded,
				codes.Internal,
				codes.ResourceExhausted,
				codes.DataLoss,
				codes.Unknown,
			}
			return !slices.Contains(statusCodeFail, st.Code())
		},
		OnStateChange: func(name string, from, to gobreaker.State) {
			log.Logger.Info(fmt.Sprintf("[%s] state change: %s -> %s", name, from, to))
		},
	}
	return gobreaker.NewCircuitBreaker[any](setting)
}
