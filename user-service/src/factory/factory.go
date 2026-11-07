package factory

import (
	"user-service/src/cache"
	"user-service/src/publisher"
	"user-service/src/repository"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Factory struct {
	UserRepository         repository.User
	CredentialRepository   repository.Credential
	RefreshTokenRepository repository.RefreshToken
	UserBlockRepository    repository.UserBlock
	CacheRepository        cache.Cache
	RabbitMQPublisher      *publisher.RabitMQ
	RestfullLogPublisher   *publisher.Kafka
	GrpcLogPublisher       *publisher.Kafka
}

// fucntion ini akan di panggil di main.go
func NewFactory(db *gorm.DB, rc *redis.ClusterClient, pr *publisher.RabitMQ, kf *publisher.Kafka, gl *publisher.Kafka) *Factory {
	return &Factory{
		UserRepository:         repository.NewUser(db), //panggil function NewUser dari repository
		RefreshTokenRepository: repository.NewRefreshToken(db),
		CredentialRepository:   repository.NewCredential(db),
		UserBlockRepository:    repository.NewUserBlocked(db),
		CacheRepository:        cache.NewCache(rc),
		RabbitMQPublisher:      pr,
		RestfullLogPublisher:   kf,
		GrpcLogPublisher:       gl,
	}
}
