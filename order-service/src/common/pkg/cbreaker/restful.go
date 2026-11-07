package cbreaker

import (
	"net/http"
	"slices"
	"time"

	"fmt"
	"order-service/src/common/errors"
	"order-service/src/common/log"

	"github.com/sony/gobreaker/v2"
)

func NewRestfulBreaker(name string) *gobreaker.CircuitBreaker[any] {
	setting := gobreaker.Settings{
		Name:        name,
		MaxRequests: 5,
		Interval:    1 * time.Minute,
		Timeout:     1 * time.Minute,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRate := float32(counts.TotalFailures) / float32(counts.Requests)
			return failureRate >= 0.8 && counts.Requests >= 5
		},
		IsSuccessful: func(err error) bool {
			if err == nil {
				return true
			}
			if e, ok := err.(*errors.Response); ok {
				statusCodeFail := []int{
					http.StatusInternalServerError,
					http.StatusBadGateway,
					http.StatusServiceUnavailable,
					http.StatusGatewayTimeout,
				}
				return !slices.Contains(statusCodeFail, e.HttpCode)
			}
			return false
		},
		OnStateChange: func(name string, from, to gobreaker.State) {
			log.Logger.Info(fmt.Sprintf("[%s] state change: %s -> %s", name, from, to))
		},
	}
	return gobreaker.NewCircuitBreaker[any](setting)
}
