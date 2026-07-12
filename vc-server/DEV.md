# VarChess — Local Development Guide

## Prerequisites

- Go 1.21+
- Docker (for Redis and MongoDB)
- `vc-cli` binary (see below)

---

## First-time setup

```bash
cd vc-server
cp .env.example .env   # fill in secrets if needed; defaults work for local dev
make deps-up           # starts Redis on :6379 and MongoDB on :27017
make build-cli         # produces bin/vc-cli
```

---

## Fast variant testing — no infrastructure

The engine REPL runs entirely in memory. It needs no Redis, MongoDB, or running
servers. Use it whenever you are developing or debugging move validation logic.

```bash
# Standard chess
./bin/vc-cli engine repl

# Antichess
./bin/vc-cli engine repl --variant antichess

# 3-Check
./bin/vc-cli engine repl --variant ncheck --target-checks 3

# Custom FEN
./bin/vc-cli engine repl --variant checkmate \
  --fen "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1"
```

### REPL commands

| Command | Action |
|---------|--------|
| `board` | Print the current board |
| `moves` | List all legal moves with numbers |
| `move <n>` | Play move number `<n>` from the last listing |
| `move e2e4` | Play by algebraic squares |
| `move e7e8Q` | Promotion — append piece letter |
| `fen` | Print current FEN |
| `variant` | Print variant type and serialised state |
| `quit` | Exit |

---

## Running all chess engine tests

These tests do not need any running services:

```bash
cd vc-server
make test-chess
# or with verbose output:
go test ./chess/... -v -count=1
```

---

## Full-stack two-terminal multiplayer test

Open **four terminal windows** (or panes):

### Terminal 1 — vc-core

```bash
cd vc-server
make run-core
```

### Terminal 2 — vc-ws

```bash
cd vc-server
make run-ws
```

### Terminal 3 — Player 1 (White)

```bash
cd vc-server

# Sign up once (skip if already done)
curl -s -X POST localhost:5000/signup \
  -H 'Content-Type: application/json' \
  -d '{"username":"alice","email":"alice@example.com","password":"secret"}'

# Login and get a token
TOKEN=$(./bin/vc-cli play login --user alice --pass secret --token-out -)

# Create a game (prints the game ID)
GAME=$(./bin/vc-cli play create --variant ncheck --target-checks 3 --token $TOKEN)
echo "Game ID: $GAME"

# Share the game ID with terminal 4, then join
./bin/vc-cli play join $GAME --color w --token $TOKEN
```

### Terminal 4 — Player 2 (Black)

```bash
cd vc-server

# Sign up once
curl -s -X POST localhost:5000/signup \
  -H 'Content-Type: application/json' \
  -d '{"username":"bob","email":"bob@example.com","password":"secret"}'

TOKEN=$(./bin/vc-cli play login --user bob --pass secret --token-out -)

# Replace GAME_ID with the ID printed in terminal 3
./bin/vc-cli play join GAME_ID --color b --token $TOKEN
```

Once both players join, the `start` event fires and each terminal prompts for
moves on your turn.

### In-game commands

| Command | Action |
|---------|--------|
| `moves` | List legal moves |
| `e2e4` | Play by algebraic squares |
| `resign` | Resign the game |
| `board` | Reprint the board |
| `quit` | Disconnect |

---

## Ports reference

| Port | Service |
|------|---------|
| 5000 | vc-core HTTP (REST + auth) |
| 8080 | vc-ws WebSocket |
| 6379 | Redis |
| 27017 | MongoDB |

---

## Environment variables (`.env`)

| Variable | Description | Default |
|----------|-------------|---------|
| `SERVER_PORT` | vc-core listen port | `5000` |
| `DB_URI` | MongoDB connection string | `mongodb://localhost:27017` |
| `DB_NAME` | MongoDB database name | `varchess` |
| `REDIS_ADDR` | Redis address | `localhost:6379` |
| `REDIS_PASSWORD` | Redis password (optional) | _(empty)_ |
| `JWT_SECRET_KEY` | JWT signing secret | _(required)_ |
| `ENVIRONMENT` | If non-empty, missing `.env` is tolerated | _(empty)_ |

---

## Makefile targets

```
make deps-up      Start Redis + MongoDB
make deps-down    Stop all Docker services
make run-core     Run vc-core on :5000
make run-ws       Run vc-ws on :8080
make test-chess   Run chess engine/variant tests
make test-all     Run all tests
make build-cli    Build bin/vc-cli
make clean        Remove build artifacts
make dev          Print the two-terminal quick-start
make help         Show all targets
```
