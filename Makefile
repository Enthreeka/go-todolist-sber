server:
	go run cmd/server/main.go

docker-up:
	docker compose -f docker-compose.yaml up --build

docker-down:
	docker compose -f docker-compose.yaml down

swagger:
	swag init -g cmd/server/main.go

test:
	go test -v ./internal/...