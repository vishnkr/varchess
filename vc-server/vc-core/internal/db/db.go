package db

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
	"vc-server/vc-core/internal/config"
	"vc-server/vc-core/internal/worker"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	client *mongo.Client
	dbName string
	RedisClient *redis.Client
}

func Connect(cfg *config.Config) (*DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(cfg.DBURI)
	clientOptions.SetMaxPoolSize(100)
	client, err := mongo.Connect(ctx, clientOptions)
	collections := []string{"games", "users", "templates"}
	if err != nil {
		return nil, err
	}

	if err:= client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	log.Println("Connected to MongoDB")
	
	if err := ensureCollections(client, cfg.DBName, collections); err != nil {
		return nil, fmt.Errorf("failed to ensure collections: %w", err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
		Password: cfg.RedisPassword,
	})
	if err := redisClient.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}
	log.Println("Connected to Redis")

	return &DB{client,cfg.DBName,redisClient}, nil
}

func ensureCollections(client *mongo.Client, dbName string, collections []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := client.Database(dbName)
	existingCollections, err := db.ListCollectionNames(ctx, bson.M{})
	if err != nil {
		return err
	}

	existingSet := make(map[string]struct{})
	for _, coll := range existingCollections {
		existingSet[coll] = struct{}{}
	}

	for _, coll := range collections {
		if _, exists := existingSet[coll]; !exists {
			if err := db.CreateCollection(ctx, coll); err != nil {
				return fmt.Errorf("error creating collection %s: %w", coll, err)
			}
			log.Printf("created collection: %s\n", coll)
		} else {
			log.Printf("collection %s already exists\n", coll)
		}
	}

	return nil
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
	ch := []string{worker.UserAction}
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
func (db *DB) ListenForMoveValidation(ctx context.Context, wp *worker.WorkerPool) {
    for {
        res, err := db.RedisClient.XRead(ctx, &redis.XReadArgs{
            Streams: []string{"move_validation", "$"},
            Count:   1,
            Block:   0,
        }).Result()

        if err != nil {
            log.Println("Error reading from move_validation:", err)
            continue
        }

        for _, stream := range res {
            for _, msg := range stream.Messages {
                var event worker.Event
                if err := json.Unmarshal([]byte(msg.Values["data"].(string)), &event); err != nil {
                    log.Println("Invalid move data:", err)
                    continue
                }
                wp.EventChannel <- event
            }
        }
    }
}
