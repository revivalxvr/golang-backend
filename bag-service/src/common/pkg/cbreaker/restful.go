package cbreaker

import (
	"slices"
	"time"

	"bag-service/src/common/errors"
	"bag-service/src/common/log"
	"fmt"

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
				successCode := []int{200, 201, 202, 204, 404}
				return slices.Contains(successCode, e.HttpCode)
			}
			return false
		},
		OnStateChange: func(name string, from, to gobreaker.State) {
			log.Logger.Info(fmt.Sprintf("[%s] state change: %s -> %s", name, from, to))
		},
	}
	return gobreaker.NewCircuitBreaker[any](setting)
}
