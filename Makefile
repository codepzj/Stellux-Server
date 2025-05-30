start:
	docker compose -f docker-compose.development.yaml up -d

rebuild:
	docker compose -f docker-compose.development.yaml down
	rm -rf data/mongo
	mkdir -p data/mongo
	docker compose -f docker-compose.development.yaml up -d

stop:
	docker compose -f docker-compose.development.yaml down