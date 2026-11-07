package database

import (
	"user-service/env"
	"user-service/src/common/log"

	"github.com/redis/go-redis/v9"
)

func NewRedis() *redis.ClusterClient {
	redisCluster := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{
			env.Conf.Redis.AddrNode1,
			env.Conf.Redis.AddrNode2,
			env.Conf.Redis.AddrNode3,
			env.Conf.Redis.AddrNode4,
			env.Conf.Redis.AddrNode5,
			env.Conf.Redis.AddrNode6,
		},
		Password: env.Conf.Redis.Password,

		// Tambahkan ini agar IP internal Docker dialihkan ke port localhost
		NewClient: func(opt *redis.Options) *redis.Client {
			switch opt.Addr {
			case "172.48.0.11:6379":
				opt.Addr = env.Conf.Redis.AddrNode1
			case "172.48.0.12:6379":
				opt.Addr = env.Conf.Redis.AddrNode2
			case "172.48.0.13:6379":
				opt.Addr = env.Conf.Redis.AddrNode3
			case "172.48.0.14:6379":
				opt.Addr = env.Conf.Redis.AddrNode4
			case "172.48.0.15:6379":
				opt.Addr = env.Conf.Redis.AddrNode5
			case "172.48.0.16:6379":
				opt.Addr = env.Conf.Redis.AddrNode6
			}
			return redis.NewClient(opt)
		},
	})

	return redisCluster
}

func CloseRedis(redis *redis.ClusterClient) {
	if err := redis.Close(); err != nil {
		log.Logger.Error(err.Error())
		return
	}
}
