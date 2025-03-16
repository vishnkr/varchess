package wss

import (
	"context"
	"encoding/json"
	"fmt"
	"vc-server/vc-ws/internal/config"

	"github.com/redis/go-redis/v9"
)

var (
	redisClient *redis.Client
)

func InitRedis() {
    cfg := config.LoadConfig()
    redisClient = redis.NewClient(&redis.Options{
        Addr: cfg.RedisAddr,
    })
	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		fmt.Printf("Failed to connect to Redis: %v", err)
        return
	}
	fmt.Println("Connected to Redis")
    SubscribeToRedis()
}

func GetRedisClient() *redis.Client {
    return redisClient
}

func SubscribeToRedis() {
    pubsub := redisClient.Subscribe(context.Background(), "game_updates")
    ch := pubsub.Channel()

    go func() {
        for msg := range ch {
            var event Event
            if err := json.Unmarshal([]byte(msg.Payload), &event); err != nil {
                continue
            }

            game,err := hub.GetGame(event.GameID)
            if err == nil {
                game.Broadcast(event)
            }
        }
    }()
}

func PublishGameUpdate(redisClient *redis.Client, gameID string, eventType string, data interface{}) {
	message := map[string]interface{}{
		"gameId": gameID,
		"event":  eventType,
		"data":   data,
	}
	msgJSON, _ := json.Marshal(message)
	redisClient.Publish(context.TODO(), "game_events", msgJSON)
}