package wss

import (
	"context"
	"encoding/json"
	"fmt"
	"vc-ws/internal/config"

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
}

func GetRedisClient() *redis.Client {
    return redisClient
}

func PublishToRedis(gameID string, msg []byte) {
    redisClient.Publish(context.Background(), "game."+gameID, msg)
}

func SubscribeToRedis() {
    pubsub := redisClient.Subscribe(context.Background(), "*_resp")
    ch := pubsub.Channel()

    go func() {
        for msg := range ch {
            var event struct {
                GameID  string `json:"game_id"`
                Payload string `json:"payload"`
            }
            if err := json.Unmarshal([]byte(msg.Payload), &event); err != nil {
                continue
            }

            game,err := hub.GetGame(event.GameID)
            if err != nil {
                game.Broadcast([]byte(event.Payload))
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