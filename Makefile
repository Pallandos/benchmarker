APP_NAME=benchmarker
CMD_PATH=./cmd/benchmarker

.PHONY: build run clean

build:
	go build -o $(APP_NAME) $(CMD_PATH)

run: build	
	./$(APP_NAME)

clean:
	rm -f $(APP_NAME)