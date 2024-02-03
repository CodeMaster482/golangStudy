.PHONY: build test clean run up down migrate

GOCMD=go
GORUN=$(GOCMD) run
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOCLEAN=$(GOCMD) clean
BINARY_NAME=server

API_MAIN=./cmd/atm/main.go
MIGRATION=./cmd/migration/main.go

all: build up migrate run 

build:
	$(GOBUILD) -o $(BINARY_NAME) $(API_MAIN)

run:
	$(BINARY_NAME)

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	sudo docker system prune -f

up:
	docker-compose up --remove-orphans

down:
	docker-compose down

migrate:
	$(GORUN) $(MIGRATION)
