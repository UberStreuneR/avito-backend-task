build:
	docker-compose up --build -d
down:
	docker-compose down -v
test:
	go test ./tests/services
run-db:
	docker-compose up -d postgres