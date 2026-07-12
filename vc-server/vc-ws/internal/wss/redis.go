package wss

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"vc-server/vc-ws/internal/config"

	"github.com/redis/go-redis/v9"
)

var (
	redisClient *redis.Client
	hub         *Hub
)

func InitRedis(cfg *config.Config) {
	redisClient = redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})
	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		fmt.Printf("Failed to connect to Redis: %v\n", err)
		return
	}
	fmt.Println("Connected to Redis")
	hub = NewHub(redisClient)
	subscribeToRedis()
}

func GetRedisClient() *redis.Client {
	return redisClient
}

func subscribeToRedis() {
	pubsub := redisClient.Subscribe(context.Background(), SystemAction)
	ch := pubsub.Channel()

	go func() {
		for msg := range ch {
			var event Event
			if err := json.Unmarshal([]byte(msg.Payload), &event); err != nil {
				log.Printf("Failed to decode system event: %v", err)
				continue
			}
			hub.HandleSystemEvent(event)
		}
	}()
}

func PublishGameUpdate(rc *redis.Client, gameID string, eventType string, data interface{}) {
	message := map[string]interface{}{
		"gameId": gameID,
		"event":  eventType,
		"data":   data,
	}
	msgJSON, _ := json.Marshal(message)
	rc.Publish(context.TODO(), UserAction, msgJSON)
}

func PublishChatMessage(gameID, playerID, message string) {
	msg := map[string]interface{}{
		"gameId":  gameID,
		"player":  playerID,
		"message": message,
	}
	msgJSON, _ := json.Marshal(msg)
	redisClient.Publish(context.Background(), "game_chat", msgJSON)
}
