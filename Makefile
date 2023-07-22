SHELL := /bin/bash

.DEFAULT_GOAL := all
.PHONY: all
all: ## build pipeline
all: mod build

.PHONY: ci
ci: ## CI build pipeline
ci: all

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: clean
clean: ## remove files created during build pipeline
	$(call print-target)
	rm -rf dist
	rm -f coverage.*
	rm -f '"$(shell go env GOCACHE)/../golangci-lint"'
	go clean -i -cache -testcache -modcache -fuzzcache -x

.PHONY: mod
mod: ## go mod tidy
	$(call print-target)
	go mod tidy

.PHONY: build
build: ## go build
	$(call print-target)
	go build ./...

.PHONY: spell
spell: ## misspell
	$(call print-target)
	go get -u github.com/client9/misspell/cmd/misspell
	misspell -error -locale=US **/*.md

# .PHONY: test
# test: ## go test
# 	$(call print-target)
# 	go test -race -covermode=atomic -coverprofile=coverage.out -coverpkg=./... ./...
# 	go tool cover -html=coverage.out -o coverage.html

# .PHONY: diff
# diff: ## git diff
# 	$(call print-target)
# 	git diff --exit-code
# 	RES=$$(git status --porcelain) ; if [ -n "$$RES" ]; then echo $$RES && exit 1 ; fi


define print-target
    @printf "Executing target: \033[36m$@\033[0m\n"
endef

.PHONY: up-mongo
up-mongo: ## docker-compose up for mongo
	$(call print-target)
	docker compose -f docker-compose.mongo.yaml up  --build -d

.PHONY: down-mongo
down-mongo: ## docker-compose down for mongo
	$(call print-target)
	docker compose -f docker-compose.mongo.yaml down

.PHONY: up-influx
up-influx: ## docker-compose up for influx
	$(call print-target)
	docker compose -f docker-compose.influx.yaml up  --build -d

.PHONY: down-influx
down-influx: ## docker-compose down for influx
	$(call print-target)
	docker compose -f docker-compose.influx.yaml down

.PHONY: up-timescale
up-timescale: ## docker-compose up for timescale
	$(call print-target)
	docker compose -f docker-compose.timescale.yaml up  --build -d

.PHONY: down-timescale
down-timescale: ## docker-compose down for timescale
	$(call print-target)
	docker compose -f docker-compose.timescale.yaml down

.PHONY: up-loki
up-loki: ## docker-compose up for loki
	$(call print-target)
	docker compose -f docker-compose.loki.yaml up  --build -d

.PHONY: down-loki
down-loki: ## docker-compose down for loki
	$(call print-target)
	docker compose -f docker-compose.loki.yaml down

.PHONY: up-kafka
up-kafka: ## docker-compose up for kafka
	$(call print-target)
	docker compose -f docker-compose.kafka.yaml up  --build -d

.PHONY: down-kafka
down-kafka: ## docker-compose down for kafka
	$(call print-target)
	docker compose -f docker-compose.kafka.yaml down

.PHONY: logs-mongo
logs-mongo: ## docker-compose logs for mongo
	$(call print-target)
	docker compose -f docker-compose.mongo.yaml logs -f

.PHONY: logs-influx
logs-influx: ## docker-compose logs for influx
	$(call print-target)
	docker compose -f docker-compose.influx.yaml logs -f

.PHONY: logs-timescale
logs-timescale: ## docker-compose logs for timescale
	$(call print-target)
	docker compose -f docker-compose.timescale.yaml logs -f

.PHONY: logs-loki
logs-loki: ## docker-compose logs for loki
	$(call print-target)
	docker compose -f docker-compose.loki.yaml logs -f

.PHONY: logs-kafka
logs-kafka: ## docker-compose logs for kafka
	$(call print-target)
	docker compose -f docker-compose.kafka.yaml logs -f
