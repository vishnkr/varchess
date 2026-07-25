package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"vc-server/chess"
	"vc-server/chess/variants"
	"vc-server/vc-core/internal/db"
	e "vc-server/vc-core/internal/event"
	"vc-server/vc-core/internal/models"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type WorkerPool struct {
	workers      []*Worker
	EventChannel chan e.Event
	wg           sync.WaitGroup
	db           *db.DB
}

type Worker struct {
	ID   int
	pool *WorkerPool
}

func NewWorkerPool(workerCount int, db *db.DB) *WorkerPool {
	wp := &WorkerPool{
		EventChannel: make(chan e.Event, 100),
		db:           db,
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
		w.processMove(event)
	case e.Join:
		w.processJoin(event)
	case e.DrawOffer:
		w.processDrawOffer(event)
	case e.DrawAccept:
		w.processDrawAccept(event)
	case e.DrawReject:
		w.processDrawReject(event)
	case e.Resign:
		w.processResign(event)
	case e.Hints:
		w.processHints(event)
	case e.GameOver:
	default:
		fmt.Println("Unknown event type:", event.Type)
	}
}

func (w *Worker) processJoin(event e.Event) {
	r := w.pool.db.RedisClient
	ctx := context.Background()
	var payload e.JoinPayload
	if err := json.Unmarshal([]byte(event.Data), &payload); err != nil {
		log.Println("Invalid join payload:", err)
		return
	}
	gameKey := fmt.Sprintf("game:%s", event.GameID)
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
	// Idempotent reconnect: player already seated — skip reassignment.
	for color, p := range activeGame.Players {
		if p.UserId.Hex() == event.UserID {
			log.Printf("Join ignored for reconnecting player %s as %s", event.UserID, color)
			if activeGame.State == models.InProgress {
				startEvent := e.Event{
					GameID: event.GameID,
					UserID: "",
					Type:   e.Start,
					Data:   json.RawMessage(gameStateJSON),
				}
				if eventJSON, err := json.Marshal(startEvent); err == nil {
					_ = r.Publish(ctx, e.SystemAction, eventJSON).Err()
				}
			}
			return
		}
	}
	if payload.Color != "w" && payload.Color != "b" {
		log.Println("Invalid color:", payload.Color)
		return
	}
	player, err := w.pool.getUserDetailsFromId(event.UserID)
	if err != nil {
		return
	}
	activeGame.Players[payload.Color] = player
	playerW, okW := activeGame.Players["w"]
	playerB, okB := activeGame.Players["b"]
	canStart := okW && okB && playerW.UserId != playerB.UserId

	if canStart {
		activeGame.State = models.InProgress
		activeGame.Turn = "w"
	}

	updatedGameStateJSON, err := json.Marshal(activeGame)
	if err != nil {
		log.Println("Error serialising updated game state:", err)
		return
	}
	if err := r.Set(ctx, gameKey, updatedGameStateJSON, 30*time.Minute).Err(); err != nil {
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
			log.Println("Error serialising game start event:", err)
			return
		}
		if err := r.Publish(ctx, e.SystemAction, eventJSON).Err(); err != nil {
			log.Println("Error publishing game start event:", err)
		} else {
			log.Printf("Game %s started", event.GameID)
		}
	}
}

// processMove validates the move against the chess engine and, if legal,
// persists the new game state and broadcasts the move to connected clients.
func (w *Worker) processMove(event e.Event) {
	r := w.pool.db.RedisClient
	ctx := context.Background()

	var payload e.MovePayload
	if err := json.Unmarshal([]byte(event.Data), &payload); err != nil {
		log.Println("Invalid move payload:", err)
		return
	}

	gameKey := fmt.Sprintf("game:%s", event.GameID)
	gameStateJSON, err := r.Get(ctx, gameKey).Result()
	if err == redis.Nil {
		log.Println("Game not found for move:", event.GameID)
		return
	} else if err != nil {
		log.Println("Error fetching game for move:", err)
		return
	}

	var activeGame models.ActiveGame
	if err := json.Unmarshal([]byte(gameStateJSON), &activeGame); err != nil {
		log.Println("Invalid game state:", err)
		return
	}
	if activeGame.State != models.InProgress {
		log.Println("Move on non-active game:", event.GameID)
		return
	}

	// Verify it is the player's turn.
	expectedColor := "w"
	if activeGame.Turn == "b" {
		expectedColor = "b"
	}
	playerColor := ""
	for color, p := range activeGame.Players {
		if p.UserId.Hex() == event.UserID {
			playerColor = color
			break
		}
	}
	if playerColor != expectedColor {
		log.Printf("Move rejected: wrong turn (expected %s, got player %s)", expectedColor, playerColor)
		return
	}

	// Reconstruct the engine from the initial config and replay all moves.
	engine, err := chess.NewEngine(activeGame.Config)
	if err != nil {
		log.Println("Failed to create engine:", err)
		return
	}

	// Replay recorded moves to reach current position (also rebuilds variant state).
	for _, moveStr := range activeGame.Moves {
		var mj chess.MoveJSON
		if err := json.Unmarshal([]byte(moveStr), &mj); err != nil {
			log.Println("Invalid recorded move:", err)
			return
		}
		m, err := mj.ToMove()
		if err != nil {
			log.Println("Invalid recorded move data:", err)
			return
		}
		if _, err := engine.PerformMove(m); err != nil {
			log.Printf("Replay failed at move %s: %v", moveStr, err)
			return
		}
	}

	// Validate and apply the incoming move.
	// Clients may send incomplete MoveJSON (missing classic_move_type/capture);
	// resolve against legal moves by from/to/piece when needed.
	incomingMove, err := payload.Move.ToMove()
	if err != nil {
		log.Println("Invalid incoming move:", err)
		return
	}
	resolved := resolveLegalMove(engine, incomingMove)
	if resolved == nil {
		log.Printf("Illegal move from %s in game %s: %+v", event.UserID, event.GameID, incomingMove)
		return
	}
	incomingMove = *resolved

	moveResult, err := engine.PerformMove(incomingMove)
	if err != nil {
		log.Printf("Illegal move from %s in game %s: %v", event.UserID, event.GameID, err)
		return
	}

	// Append resolved move to history.
	resolvedJSON := incomingMove.ToJSON()
	moveJSONBytes, err := json.Marshal(resolvedJSON)
	if err != nil {
		log.Println("Failed to serialise move:", err)
		return
	}
	activeGame.Moves = append(activeGame.Moves, string(moveJSONBytes))

	// Flip turn.
	if activeGame.Turn == "w" {
		activeGame.Turn = "b"
	} else {
		activeGame.Turn = "w"
	}

	// Persist variant state for reconnects / client UI (wormhole cooldowns, etc.).
	var variantStateJSON json.RawMessage
	if state := engine.VariantState(); state != nil {
		if b, err := json.Marshal(state); err == nil {
			activeGame.VariantState = b
			variantStateJSON = b
		}
	} else {
		activeGame.VariantState = nil
	}

	if moveResult.IsGameOver && moveResult.Result != nil {
		activeGame.State = models.Complete
		w.persistCompletedGame(ctx, &activeGame, moveResult)
	}

	updatedJSON, err := json.Marshal(activeGame)
	if err != nil {
		log.Println("Error serialising game state:", err)
		return
	}
	if err := r.Set(ctx, gameKey, updatedJSON, 30*time.Minute).Err(); err != nil {
		log.Println("Error saving game state:", err)
		return
	}

	// Always broadcast the confirmed move (incl. check flag) so clients can
	// apply it and highlight the checked king. Game-over is a follow-up event.
	confirmedPayload, _ := json.Marshal(e.MovePayload{
		Move:         resolvedJSON,
		IsCheck:      moveResult.IsCheck,
		VariantState: variantStateJSON,
	})
	outEvent := e.Event{
		GameID: event.GameID,
		UserID: event.UserID,
		Type:   e.Move,
		Data:   confirmedPayload,
	}
	eventJSON, err := json.Marshal(outEvent)
	if err != nil {
		log.Println("Error serialising move event:", err)
		return
	}
	r.Publish(ctx, e.SystemAction, eventJSON)

	if moveResult.IsGameOver && moveResult.Result != nil {
		resultJSON, _ := json.Marshal(moveResult.Result)
		overEvent := e.Event{
			GameID: event.GameID,
			UserID: event.UserID,
			Type:   e.GameOver,
			Data:   resultJSON,
		}
		overJSON, err := json.Marshal(overEvent)
		if err != nil {
			log.Println("Error serialising game-over event:", err)
			return
		}
		r.Publish(ctx, e.SystemAction, overJSON)
	}
}

func (w *Worker) processResign(event e.Event) {
	r := w.pool.db.RedisClient
	ctx := context.Background()
	gameKey := fmt.Sprintf("game:%s", event.GameID)

	gameStateJSON, err := r.Get(ctx, gameKey).Result()
	if err != nil {
		return
	}
	var activeGame models.ActiveGame
	if err := json.Unmarshal([]byte(gameStateJSON), &activeGame); err != nil {
		return
	}

	var winner string
	for color, p := range activeGame.Players {
		if p.UserId.Hex() == event.UserID {
			if color == "w" {
				winner = "black"
			} else {
				winner = "white"
			}
			break
		}
	}
	activeGame.State = models.Complete
	result := chess.MoveResult{
		IsGameOver: true,
		Result:     &chess.EngineResult{Winner: winner, Reason: "resignation"},
	}
	w.persistCompletedGame(ctx, &activeGame, result)

	resultJSON, _ := json.Marshal(result.Result)
	outEvent := e.Event{
		GameID: event.GameID,
		UserID: event.UserID,
		Type:   e.GameOver,
		Data:   resultJSON,
	}
	eventJSON, _ := json.Marshal(outEvent)
	r.Publish(ctx, e.SystemAction, eventJSON)
}

func (w *Worker) loadActiveGame(ctx context.Context, gameID string) (*models.ActiveGame, string, error) {
	r := w.pool.db.RedisClient
	gameKey := fmt.Sprintf("game:%s", gameID)
	gameStateJSON, err := r.Get(ctx, gameKey).Result()
	if err != nil {
		return nil, gameKey, err
	}
	var activeGame models.ActiveGame
	if err := json.Unmarshal([]byte(gameStateJSON), &activeGame); err != nil {
		return nil, gameKey, err
	}
	return &activeGame, gameKey, nil
}

func (w *Worker) isPlayerInGame(activeGame *models.ActiveGame, userID string) bool {
	for _, p := range activeGame.Players {
		if p.UserId.Hex() == userID {
			return true
		}
	}
	return false
}

// processDrawOffer relays a draw offer to the opponent via the WS hub.
func (w *Worker) processDrawOffer(event e.Event) {
	r := w.pool.db.RedisClient
	ctx := context.Background()
	activeGame, _, err := w.loadActiveGame(ctx, event.GameID)
	if err != nil || activeGame.State != models.InProgress {
		return
	}
	if !w.isPlayerInGame(activeGame, event.UserID) {
		return
	}
	outEvent := e.Event{
		GameID: event.GameID,
		UserID: event.UserID,
		Type:   e.DrawOffer,
		Data:   nil,
	}
	eventJSON, err := json.Marshal(outEvent)
	if err != nil {
		return
	}
	r.Publish(ctx, e.SystemAction, eventJSON)
}

// processDrawAccept ends the game as a draw and notifies both clients.
func (w *Worker) processDrawAccept(event e.Event) {
	r := w.pool.db.RedisClient
	ctx := context.Background()
	activeGame, gameKey, err := w.loadActiveGame(ctx, event.GameID)
	if err != nil || activeGame.State != models.InProgress {
		return
	}
	if !w.isPlayerInGame(activeGame, event.UserID) {
		return
	}

	activeGame.State = models.Complete
	result := chess.MoveResult{
		IsGameOver: true,
		Result:     &chess.EngineResult{Winner: "draw", Reason: "agreement"},
	}
	w.persistCompletedGame(ctx, activeGame, result)

	updatedJSON, err := json.Marshal(activeGame)
	if err == nil {
		_ = r.Set(ctx, gameKey, updatedJSON, 30*time.Minute).Err()
	}

	resultJSON, _ := json.Marshal(result.Result)
	// Broadcast draw-accept first so FE can clear UI, then game-over.
	acceptEvent := e.Event{
		GameID: event.GameID,
		UserID: event.UserID,
		Type:   e.DrawAccept,
		Data:   resultJSON,
	}
	if b, err := json.Marshal(acceptEvent); err == nil {
		r.Publish(ctx, e.SystemAction, b)
	}
	overEvent := e.Event{
		GameID: event.GameID,
		UserID: event.UserID,
		Type:   e.GameOver,
		Data:   resultJSON,
	}
	if b, err := json.Marshal(overEvent); err == nil {
		r.Publish(ctx, e.SystemAction, b)
	}
}

// processDrawReject relays a declined draw offer to the offering player.
func (w *Worker) processDrawReject(event e.Event) {
	r := w.pool.db.RedisClient
	ctx := context.Background()
	activeGame, _, err := w.loadActiveGame(ctx, event.GameID)
	if err != nil || activeGame.State != models.InProgress {
		return
	}
	if !w.isPlayerInGame(activeGame, event.UserID) {
		return
	}
	outEvent := e.Event{
		GameID: event.GameID,
		UserID: event.UserID,
		Type:   e.DrawReject,
		Data:   nil,
	}
	eventJSON, err := json.Marshal(outEvent)
	if err != nil {
		return
	}
	r.Publish(ctx, e.SystemAction, eventJSON)
}

// processHints returns legal destination squares for a selected piece (requester only).
func (w *Worker) processHints(event e.Event) {
	r := w.pool.db.RedisClient
	ctx := context.Background()

	var payload e.HintsPayload
	if err := json.Unmarshal(event.Data, &payload); err != nil {
		log.Println("Invalid hints payload:", err)
		return
	}

	activeGame, _, err := w.loadActiveGame(ctx, event.GameID)
	if err != nil || activeGame.State != models.InProgress {
		return
	}
	if !w.isPlayerInGame(activeGame, event.UserID) {
		return
	}

	engine, err := w.engineAtCurrentPosition(activeGame)
	if err != nil {
		log.Println("hints engine:", err)
		return
	}

	tos := make([]int, 0)
	for _, lm := range engine.GetLegalMoves() {
		if lm.From != payload.From {
			continue
		}
		// For wormhole teleports, hint the portal entrance (what the player clicks),
		// not the exit square the piece ends on.
		if portal, ok := variants.TeleportPortalOf(lm); ok {
			tos = append(tos, portal)
			continue
		}
		tos = append(tos, lm.To)
	}

	result := e.HintsResultPayload{From: payload.From, To: tos}
	resultJSON, err := json.Marshal(result)
	if err != nil {
		return
	}
	outEvent := e.Event{
		GameID: event.GameID,
		UserID: event.UserID,
		Type:   e.Hints,
		Data:   resultJSON,
	}
	eventJSON, err := json.Marshal(outEvent)
	if err != nil {
		return
	}
	r.Publish(ctx, e.SystemAction, eventJSON)
}

func (w *Worker) engineAtCurrentPosition(activeGame *models.ActiveGame) (chess.Engine, error) {
	engine, err := chess.NewEngine(activeGame.Config)
	if err != nil {
		return nil, err
	}
	for _, moveStr := range activeGame.Moves {
		var mj chess.MoveJSON
		if err := json.Unmarshal([]byte(moveStr), &mj); err != nil {
			return nil, err
		}
		m, err := mj.ToMove()
		if err != nil {
			return nil, err
		}
		if _, err := engine.PerformMove(m); err != nil {
			return nil, err
		}
	}
	return engine, nil
}

func (w *Worker) persistCompletedGame(ctx context.Context, activeGame *models.ActiveGame, result chess.MoveResult) {
	dbClient := w.pool.db
	collection := dbClient.Collection("games")

	players := make(map[string]primitive.ObjectID)
	playerNames := make(map[string]string)
	for color, p := range activeGame.Players {
		players[color] = p.UserId
		playerNames[color] = p.Name
	}

	winner := ""
	reason := ""
	if result.Result != nil {
		winner = result.Result.Winner
		reason = result.Result.Reason
	}

	now := primitive.NewDateTimeFromTime(time.Now())
	game := models.Game{
		ShortID:     activeGame.ID,
		Players:     players,
		PlayerNames: playerNames,
		Config:      activeGame.Config,
		Moves:       activeGame.Moves,
		Result:      models.GameResult{Winner: winner, Reason: reason},
		CreatedAt:   now,
		EndedAt:     &now,
	}
	if _, err := collection.InsertOne(ctx, game); err != nil {
		log.Println("Failed to persist completed game:", err)
	}
}

func (wp *WorkerPool) getUserDetailsFromId(userId string) (models.Player, error) {
	dbClient := wp.db
	collection := dbClient.Collection("users")
	ctx := context.Background()

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

// resolveLegalMove finds a legal move matching from/to (and piece when set).
// Lets thin clients omit classic_move_type / capture and still play.
// Wormhole: clients send To = portal entrance; legal moves use To = exit.
func resolveLegalMove(engine chess.Engine, incoming chess.Move) *chess.Move {
	var fromToMatch *chess.Move
	for _, lm := range engine.GetLegalMoves() {
		if lm.From != incoming.From {
			continue
		}
		toMatch := lm.To == incoming.To
		if !toMatch {
			if portal, ok := variants.TeleportPortalOf(lm); ok && portal == incoming.To {
				toMatch = true
			}
		}
		if !toMatch {
			continue
		}
		if incoming.Piece == 0 || lm.Piece == incoming.Piece {
			matched := lm
			return &matched
		}
		matched := lm
		fromToMatch = &matched
	}
	return fromToMatch
}
