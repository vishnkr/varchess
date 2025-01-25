package db

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
	"vc-core/internal/worker"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	client *mongo.Client
	dbName string
	RedisClient *redis.Client
}

func Connect(uri string,dbName string,redisAddr string) (*DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	clientOptions.SetMaxPoolSize(100)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	if err:= client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	log.Println("Connected to MongoDB")
	
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	if err := redisClient.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}
	log.Println("Connected to Redis")

	return &DB{client,dbName,redisClient}, nil
}

func (db *DB) Close() error {
	return db.client.Disconnect(context.Background())
}

func (db *DB) Collection(name string) *mongo.Collection {
	return db.client.Database(db.dbName).Collection(name)
}

func (db *DB) ListenAndPublishMessages(ctx context.Context,ch string) {
	sub := db.RedisClient.Subscribe(ctx, ch)
	defer sub.Close()
	for {
		msg, err := sub.ReceiveMessage(ctx)
		if err != nil {
			log.Printf("Error receiving Redis message: %v", err)
			return
		}

		ackMessage := fmt.Sprintf("Message received from channel %s: %s", msg.Channel, msg.Payload)
		if err := db.RedisClient.Publish(ctx, "validate_resp", ackMessage).Err(); err != nil {
			log.Printf("Error sending acknowledgment: %v", err)
		} else {
			log.Printf("Sent acknowledgment: %s", ackMessage)
		}
	}
}

func (db *DB) SetupRedisSubscriber(ctx context.Context,wp *worker.WorkerPool){
	ch := []string{"game_events","moves"}
	sub := db.RedisClient.Subscribe(ctx, ch...)
	defer sub.Close()

	for {
		msg, err := sub.ReceiveMessage(ctx)
		if err != nil {
			log.Println("Error receiving message:", err)
			continue
		}

		var event worker.Event
		if err := json.Unmarshal([]byte(msg.Payload), &event); err != nil {
			log.Println("Invalid event format:", err)
			continue
		}
		switch event.Type {
		case worker.Join:
			var joinData worker.JoinPayload
			if err := json.Unmarshal(event.Data, &joinData); err != nil {
				log.Println("Failed to unmarshal Join data:", err)
				continue
			}
			fmt.Printf("Join Event: %+v\n", joinData)

		case worker.Move:
			
			fmt.Printf("Move Event\n")

		default:
			log.Println("Unknown event type:", event.Type)
		}
		wp.EventChannel <- event
	}
}