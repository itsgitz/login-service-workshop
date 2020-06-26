run: up

up:
	docker-compose up --build

down:
	docker-compose down -v