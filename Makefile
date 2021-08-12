
.Phony: help build test docker install

help:
	@echo "|-----------------------------------------------------------------|"
	@echo "| The following are some of the valid targets for this Makefile:  |"
	@echo "|-----------------------------------------------------------------|"
	@echo "|																 |"
	@echo "|  - make all (the default if no target is provided)              |"
	@echo "|  - make test (to running test before we build)		        	 |"
	@echo "|  - make install (to download and check all dependency)		     |"
	@echo "|  - make run (to run service & postgresql using docker compose)	 |"
	@echo "|  - make stop (to stop docker compose)	 						 |"
	@echo "|  - make docker (to create docker image)	         			 |"
	@echo "|  - make build (to download & check dependency and build binary) |"
	@echo "|																 |"
	@echo "|-----------------------------------------------------------------|"

all: install test build

test:
	@echo "=================================================================================="
	@echo "Coverage Test"
	@echo "=================================================================================="
	go fmt ./... && go test -coverprofile coverage.cov -cover ./... # use -v for verbose
	@echo "\n"
	@echo "=================================================================================="
	@echo "All Package Coverage"
	@echo "=================================================================================="
	go tool cover -func coverage.cov

docker:
	@docker build -t referral_service -f Dockerfile .

run:
	@docker-compose up -d

stop:
	@docker-compose down

install:
	go mod tidy
	go mod download
	go mod verify

build: install
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-s -w" -o referral_service github.com/candraalim/be_tsel_candra/cmd/app
