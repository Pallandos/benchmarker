APP_NAME=benchmarker
BIN_PATH=./bin
CMD_PATH=./cmd/benchmarker

.PHONY: build run clean

build:
	go build -o $(BIN_PATH)/$(APP_NAME) $(CMD_PATH)

run: build	
	$(BIN_PATH)/$(APP_NAME)

clean:
	rm -f $(BIN_PATH)/$(APP_NAME)