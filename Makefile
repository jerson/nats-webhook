REGISTRY?=registry.gitlab.com/pardacho/nats-webhook
APP_VERSION?=latest
BUILD?=go build -ldflags="-w -s"

default: build

build: format lint
	$(BUILD) -o nats-webhook main.go

deps:
	go mod download

test:
	go test ./...

format:
	go fmt ./...

vet:
	go vet ./...

lint:
	golint -set_exit_status -min_confidence 0.3 ./...

registry: registry-build registry-push

registry-build:
	docker build --pull -t $(REGISTRY):$(APP_VERSION) .

registry-push:
	docker push $(REGISTRY):$(APP_VERSION)

stop:
	docker-compose stop

dev:
	docker-compose build
	docker-compose up -d
	clear
	@echo ""
	@echo "starting command line:"
	@echo "** when finish exist and run: make stop**"
	@echo ""
	docker-compose exec webhook sh