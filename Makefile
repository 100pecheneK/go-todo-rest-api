up:
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"

migrate:
	@echo "Starting migrations" 
	@migrate -path ./schema -database 'postgres://postgres:1425@localhost:5432/postgres?sslmode=disable' up
	@echo "Migration competed!"

migrate_down:
	@migrate -path ./schema -database 'postgres://postgres:1425@localhost:5432/postgres?sslmode=disable' down



down:
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Done!"


