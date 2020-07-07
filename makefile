run: up

up:
	docker-compose up --build

down:
	docker-compose rm -f
	docker-compose down -v
	docker system prune -f
