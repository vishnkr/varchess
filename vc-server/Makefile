CMD = ./cmd/main.go
BIN = vc-server

build:
	go build -o ${BIN} ${CMD}
run:
	./${BIN}

build_run: build run
clean:
	go clean
	rm ./{BIN}