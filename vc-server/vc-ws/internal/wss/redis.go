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
            if err != nil {
                fmt.Printf("Game not found: %v\n", event.GameID)
                continue
            }
            switch event.Type {
            case Move, Join, Resign, DrawOffer, DrawAccept, DrawReject, GameOver:
                game.Broadcast(event)

            default:
                fmt.Printf("Unhandled event type: %v\n", event.Type)
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

func PublishMoveForValidation(redisClient *redis.Client, event Event) {
    msgJSON, err := json.Marshal(event)
    if err != nil {
        fmt.Printf("Failed to marshal move: %v\n", err)
        return
    }
    redisClient.XAdd(context.TODO(), &redis.XAddArgs{
        Stream: "move_validation",
        Values: map[string]interface{}{"data": msgJSON},
    })
}


func PublishChatMessage(gameID, playerID, message string) {
    msg := map[string]interface{}{
        "gameId": gameID,
        "player": playerID,
        "message": message,
    }
    msgJSON, _ := json.Marshal(msg)
    redisClient.Publish(context.Background(), "game_chat", msgJSON)
}