package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"vc-core/internal/models"

	"github.com/redis/go-redis/v9"
)



type WorkerPool struct {
	workers   []*Worker
	EventChannel chan Event
	wg        sync.WaitGroup
	redisClient  *redis.Client
}

type Worker struct {
	ID   int
	pool *WorkerPool
}



func NewWorkerPool(workerCount int,redisClient *redis.Client) *WorkerPool{
	wp := &WorkerPool{
		EventChannel: make(chan Event, 100),
		redisClient: redisClient,
	}
	wp.startPool(10)
	return wp
}

func (w *Worker) start() {
	defer w.pool.wg.Done()
	for event := range w.pool.EventChannel {
		w.processEvent(event)
	}
}

func (wp *WorkerPool) startPool(workerCount int) {
	for i := 0; i < workerCount; i++ {
		worker := &Worker{ID: i, pool: wp}
		wp.workers = append(wp.workers, worker)
		go worker.start()
	}
}

func (wp *WorkerPool) stopPool() {
	close(wp.EventChannel)
	wp.wg.Wait()
}

func (w *Worker) processJoin(event Event){
	r := w.pool.redisClient
	ctx := context.Background()
	var payload JoinPayload
	if err := json.Unmarshal([]byte(event.Data), &payload); err != nil {
		log.Println("Invalid move payload:", err)
		return
	}
	gameKey := fmt.Sprintf("game.%s", event.GameID)
	gameStateJSON, err := r.Get(ctx, gameKey).Result()
	if err == redis.Nil {
		log.Println("Game state not found:", event.GameID)
		return
	} else if err != nil {
		log.Println("Error fetching game state:", err)
		return
	}
	var activeGame models.ActiveGame
	if err := json.Unmarshal([]byte(gameStateJSON), &activeGame); err != nil {
		log.Println("Invalid game state format:", err)
		return
	}
	if activeGame.Players == nil {
		activeGame.Players = make(map[string]string)
	}
	if payload.Color != "w" && payload.Color != "b" {
		log.Println("Invalid color:", payload.Color)
		return
	}

	activeGame.Players[payload.Color] = event.UserID
	canStart := activeGame.Players["w"] != "" && activeGame.Players["b"] != ""
	if canStart{
		activeGame.State = models.InProgress
	}
	
	updatedGameStateJSON, err := json.Marshal(activeGame)
	if err != nil {
		log.Println("Error serializing updated game state:", err)
		return
	}

	if err := r.Set(ctx, gameKey, updatedGameStateJSON, 0).Err(); err != nil {
		log.Println("Error saving updated game state:", err)
		return
	}

	if canStart {
		startEvent := Event{
			GameID: event.GameID,
			UserID: "",
			Type:   Start,
			Data:   updatedGameStateJSON,
		}
		eventJSON, err := json.Marshal(startEvent)
		if err != nil {
			log.Println("Error serializing game start event:", err)
			return
		}
		if err := r.Publish(ctx, "game_updates", eventJSON).Err(); err != nil {
			log.Println("Error publishing game start event:", err)
		} else {
			log.Printf("Game %s started: %s", event.GameID, eventJSON)
		}
	}
	
}