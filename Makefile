build:
	docker-compose build
up-db:
	docker-compose up -d db
	docker-compose up -d pma
up-go:
	docker-compose up -d go-ex