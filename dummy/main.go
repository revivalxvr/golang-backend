package main

import (
	"context"
	"dummy/producer"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {

	producer.SendRealLogTraffic()
}
func flushAllNodes() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Daftar port yang di-mapping ke host Mac kamu
	ports := []string{"5371", "5372", "5373", "5374", "5375", "5376"}

	for _, port := range ports {
		rdb := redis.NewClient(&redis.Options{
			Addr:     "127.0.0.1:" + port,
			Password: "rahasia",
		})

		err := rdb.FlushAll(ctx).Err()
		if err != nil {
			fmt.Printf("Gagal flush node di port %s: %v\n", port, err)
		} else {
			fmt.Printf("Berhasil flush node di port %s\n", port)
		}
		_ = rdb.Close()
	}
}
