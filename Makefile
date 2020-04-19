
APP=go-covid-graph
BIN=bin/$(APP)

dev:
	@go run main.go

run: build
	@echo "Run build"
	$(BIN)

build b:
	@go build -o $(BIN)

linux l:
	@GOOS=linux go build -o $(BIN)

docker:
	@docker build --no-cache -t $(APP) .

docker-run: docker
	@docker run --rm --name ${APP} -p 5000:5000 $(APP)

