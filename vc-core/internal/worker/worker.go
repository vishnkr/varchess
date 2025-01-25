package worker

import (
	"fmt"
	"sync"

	"github.com/redis/go-redis/v9"
)



type WorkerPool struct {
	workers   []*Worker
	EventChannel chan Event
	wg        sync.WaitGroup
}

type Worker struct {
	ID   int
	pool *WorkerPool
}

type Event struct {
	GameID   string
	UserID   string
	EventType string
	Data      string
}

func NewWorkerPool(workerCount int) *WorkerPool{
	wp := &WorkerPool{
		EventChannel: make(chan Event, 100),
	}
	wp.startPool(10)
	return wp
}

func (w *Worker) start() {
	defer w.pool.wg.Done()
	for event := range w.pool.EventChannel {
		processEvent(event)
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

func processEvent(event Event) {
	switch event.EventType {
	case "move":
		// move validation
	case "update":
		// pdateGameState(event)
	default:
		fmt.Println("Unknown event type:", event.EventType)
	}
}
func processMoveAndPublish(redisClient *redis.Client, gameID string, move string, currentBoard string) error {
    /* Process the move (this could involve validation, etc.)
    updatedBoard, err := processMove(move, currentBoard)
    if err != nil {
        return fmt.Errorf("failed to process move: %w", err)
    }

    // After processing the move, we need to publish the updated game state
    updatedGameState := GameState{
        GameID: gameID,
        Board:  updatedBoard,
    }

    // Serialize the updated game state into a message (JSON, for example)
    message, err := json.Marshal(updatedGameState)
    if err != nil {
        return fmt.Errorf("failed to serialize game state: %w", err)
    }

    // Publish the updated game state to the Redis channel
    err = redisClient.Publish(context.Background(), "game-updates", message).Err()
    if err != nil {
        return fmt.Errorf("failed to publish updated game state to Redis: %w", err)
    }

    log.Printf("Published updated game state to Redis: %s", message)*/
    return nil
}