package goroku

import (
	"fmt"
	"log"
	"net/url"
	"os"

	"golang.org/x/net/context"
	"gopkg.in/redis.v2"
)

type redisDB string

// MustRedis Must is a helper that wraps a call to a function returning (*redis.Client, error) and panics if the error is non-nil
func MustRedis(ctx context.Context) (*redis.Client) {
	db, ok := Redis(ctx)
	if !ok {
		panic("not found redis")
	}
	return db
}

// Redis returns the connected redis client for the default name
func Redis(ctx context.Context) (*redis.Client, bool) {
	key := redisDB("default")
	db, ok := ctx.Value(key).(*redis.Client)
	return db, ok
}

// OpenRedis open Redis connections in the context's default
func OpenRedis(ctx context.Context) context.Context {
	addr, password := GetHerokuRedisAddr()
	client := redis.NewTCPClient(&redis.Options{
		Addr:     addr,
		Password: password,
	})
	ctx = context.WithValue(ctx, redisDB("default"), client)
	return ctx
}

// GetHerokuRedisAddr returns heroku env
func GetHerokuRedisAddr() (addr string, password string) {
	addr = fmt.Sprintf("%s:%d", "localhost", 6379)
	password = ""

	redisURL := os.Getenv("REDISTOGO_URL")
	if redisURL == "" {
		fmt.Println("local: redis", addr, password)
		return
	}

	redisInfo, err := url.Parse(redisURL)
	if err != nil {
		return
	}

	addr = redisInfo.Host
	if redisInfo.User != nil {
		password, _ = redisInfo.User.Password()
	}
	return
}

// CloseRedis closes Redis connections in the context's default
func CloseRedis(ctx context.Context) context.Context {
	client, _ := Redis(ctx)
	if client == nil {
		return ctx
	}

	if err := client.Close(); err != nil {
		log.Println("redis close error:", err)
	}

	return context.WithValue(ctx, redisDB("default"), nil)
}
