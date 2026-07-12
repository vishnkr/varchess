.PHONY: deps-up deps-down run-core run-ws test-chess test-all build-cli clean dev help

# ── Infrastructure ────────────────────────────────────────────────────────────

## Start Redis and MongoDB in the background (required by vc-core and vc-ws).
deps-up:
	docker compose up redis db -d

## Stop all Docker services.
deps-down:
	docker compose down

# ── Services (run in separate terminals) ──────────────────────────────────────

## Run vc-core REST API server (port 5000).
run-core:
	go run ./vc-core

## Run vc-ws WebSocket server (port 8080).
run-ws:
	go run ./vc-ws

# ── Tests ─────────────────────────────────────────────────────────────────────

## Run chess engine and variant tests (no infrastructure needed).
test-chess:
	go test ./chess/... -v -count=1

## Run all tests in the monorepo.
test-all:
	go test ./... -count=1

# ── CLI ───────────────────────────────────────────────────────────────────────

## Build the vc-cli binary to bin/vc-cli.
build-cli:
	@mkdir -p bin
	go build -o bin/vc-cli ./cmd/vc-cli
	@echo "Built: bin/vc-cli"

## Remove build artifacts.
clean:
	rm -rf bin/

# ── Quick start ───────────────────────────────────────────────────────────────

## Print the two-terminal quick-start workflow.
dev:
	@echo ""
	@echo "Quick-start workflow:"
	@echo "  1.  make deps-up            # start Redis + MongoDB"
	@echo "  2.  make run-core           # terminal 1: REST API on :5000"
	@echo "  3.  make run-ws             # terminal 2: WS server on :8080"
	@echo "  4.  make build-cli          # build ./bin/vc-cli"
	@echo ""
	@echo "Engine REPL (no infra needed):"
	@echo "  ./bin/vc-cli engine repl --variant antichess"
	@echo "  ./bin/vc-cli engine repl --variant ncheck --target-checks 3"
	@echo ""
	@echo "Two-terminal multiplayer:"
	@echo "  # Sign up (once per user):"
	@echo "  curl -s -X POST localhost:5000/signup \\"
	@echo "    -H 'Content-Type: application/json' \\"
	@echo "    -d '{\"username\":\"alice\",\"email\":\"a@b.com\",\"password\":\"secret\"}'"
	@echo ""
	@echo "  # Terminal A (player 1 — white):"
	@echo "  TOKEN=\$$(./bin/vc-cli play login --user alice --pass secret --token-out -)"
	@echo "  GAME=\$$(./bin/vc-cli play create --variant ncheck --target-checks 3 --token \$$TOKEN)"
	@echo "  echo \"Share this game ID with player 2: \$$GAME\""
	@echo "  ./bin/vc-cli play join \$$GAME --color w --token \$$TOKEN"
	@echo ""
	@echo "  # Terminal B (player 2 — black) — replace GAME_ID with the printed ID:"
	@echo "  TOKEN=\$$(./bin/vc-cli play login --user bob --pass secret --token-out -)"
	@echo "  ./bin/vc-cli play join GAME_ID --color b --token \$$TOKEN"
	@echo ""

## Show available targets.
help:
	@grep -E '^##' Makefile | sed 's/## //'
