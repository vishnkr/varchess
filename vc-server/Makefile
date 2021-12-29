CMD = ./cmd/main.go
BIN = vc-server

build:
	go build -o ${BIN} ${CMD}
run:
	./vc-server

build_run: build run
clean:
	go clean
	rm ./vc-server