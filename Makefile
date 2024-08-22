build:
	@go build -o ./bin/teletradebot
run: build
	@./bin/teletradebot

