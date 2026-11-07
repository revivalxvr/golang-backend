package middleware

import (
	"user-service/src/cache"
	"user-service/src/factory"
	"user-service/src/publisher"
	"user-service/src/repository"
)

type Middleware struct {
	restFullLogPublisher   *publisher.Kafka
	cacheRepository        cache.Cache
	userBlockLogRepository repository.UserBlock
}

// function ini akan di panggil di server/restful.go
func NewMiddleware(f *factory.Factory) *Middleware {
	return &Middleware{
		restFullLogPublisher:   f.RestfullLogPublisher,
		cacheRepository:        f.CacheRepository,
		userBlockLogRepository: f.UserBlockRepository,
	}
}
