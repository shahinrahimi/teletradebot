build:
	CGO_ENABLED=1 @go build -o ./bin/teletradebot
run: build
	@./bin/teletradebot

test: 
	go test ./store

