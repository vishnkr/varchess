package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"vc-server/vc-core/internal/db"
	e "vc-server/vc-core/internal/event"
	"vc-server/vc-core/internal/models"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)



type WorkerPool struct {
	workers   []*Worker
	EventChannel chan e.Event
	wg        sync.WaitGroup
	db *db.DB
	
}

type Worker struct {
	ID   int
	pool *WorkerPool
}



func NewWorkerPool(workerCount int,db *db.DB) *WorkerPool{
	wp := &WorkerPool{
		EventChannel: make(chan e.Event, 100),
		db: db,
	}
	wp.startPool(workerCount)
	return wp
}

func (wp *WorkerPool) ConsumeEvent(event interface{}) {
    specificEvent, ok := event.(e.Event)
    if !ok {
        log.Println("Invalid event type received")
        return
    }
    
    wp.EventChannel <- specificEvent
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

func (w *Worker) processEvent(event e.Event) {
	switch event.Type {
	case e.Move:
		// move validation
		w.processMove(event)
	case e.Join:
		w.processJoin(event)
	case e.DrawAccept:
		// pdateGameState(event)
	case e.DrawOffer:
		// pdateGameState(event)
	case e.DrawReject:
		// pdateGameState(event)
	case e.Resign:
		// pdateGameState(event)
	case e.GameOver:
		// pdateGameState(event)
	default:
		fmt.Println("Unknown event type:", event.Type)
	}
}

func (w *Worker) processJoin(event e.Event){
	r := w.pool.db.RedisClient
	ctx := context.Background()
	var payload e.JoinPayload
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
		activeGame.Players = make(map[string]models.Player)
	}
	if payload.Color != "w" && payload.Color != "b" {
		log.Println("Invalid color:", payload.Color)
		return
	}
	player,err:= w.pool.getUserDetailsFromId(event.UserID)
	if err!=nil{
		return
	}
	activeGame.Players[payload.Color] = player
	playerW, okW := activeGame.Players["w"]
	playerB, okB := activeGame.Players["b"]
	canStart := okW && okB && playerW != playerB

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
		startEvent := e.Event{
			GameID: event.GameID,
			UserID: "",
			Type:   e.Start,
			Data:   updatedGameStateJSON,
		}
		eventJSON, err := json.Marshal(startEvent)
		if err != nil {
			log.Println("Error serializing game start event:", err)
			return
		}
		if err := r.Publish(ctx, e.SystemAction, eventJSON).Err(); err != nil {
			log.Println("Error publishing game start event:", err)
		} else {
			log.Printf("Game %s started: %s", event.GameID, eventJSON)
		}
	}
	
}

func (w *Worker) processMove(event e.Event){
	r := w.pool.db.RedisClient
	_  = context.Background()
	var payload e.MovePayload
	if err := json.Unmarshal([]byte(event.Data), &payload); err != nil {
		log.Println("Invalid move payload:", err)
		return
	}
	//validate move here 
	r.Publish(context.TODO(), e.UserAction, payload)
}

func (wp *WorkerPool) getUserDetailsFromId(userId string)(models.Player,error){
	dbClient := wp.db
	collection := dbClient.Collection("users")
	ctx := context.Background()
	//defer cancel()
	

	userIdHex, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return models.Player{}, fmt.Errorf("invalid user ID format: %w", err)
	}
	

	var user models.Player
	

	err = collection.FindOne(ctx, bson.M{"_id": userIdHex}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.Player{}, fmt.Errorf("user not found with ID: %s", userId)
		}
		return models.Player{}, fmt.Errorf("error finding user: %w", err)
	}
	
	return user, nil
	
}