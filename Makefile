.PHONY: build
build:
	go build ./cmd/apiserver

.PHONY: pgrun
pgrun:
	docker run --name testpg -p 25432:5432 -e POSTGRES_USER=test -e POSTGRES_PASSWORD=test -e POSTGRES_DB=test -d --rm postgres:14.1-alpine

.PHONY: migup
migup:
	migrate -path migrations -database "postgres://localhost:25432/test?user=test&password=test&sslmode=disable" up

.PHONY: migdown
migdown:
	migrate -path migrations -database "postgres://localhost:25432/test?user=test&password=test&sslmode=disable" down -all

.PHONY: pgstop
pgstop:
	docker stop testpg


.PHONY: test
test:
	migrate -path migrations -database "postgres://localhost:25432/test?user=test&password=test&sslmode=disable" down -all
	migrate -path migrations -database "postgres://localhost:25432/test?user=test&password=test&sslmode=disable" up
	go test -v -race -timeout 30s ./...

.DEFAULT_GOAL :=build
