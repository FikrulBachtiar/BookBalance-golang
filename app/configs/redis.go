package configs

import (
	"log"

	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

type ConnectionRedis struct {
	Addr     string
	Password string
	DB       int
}

func InitRedis(ctx context.Context, con *ConnectionRedis) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: con.Addr,
		Password: con.Password,
		DB: con.DB,
	});

	_, err := client.Ping(ctx).Result();
	if err != nil {
		log.Fatal(err);
	}

	return client;
}