DB_URL := postgres://postgres:belajar@localhost:5432/user_service?sslmode=disable

.PHONY: migration
migration:
	migrate create -ext sql -dir database/migration -seq ${name}

.PHONY: migration-down
migration-down:
	migrate -database $(DB_URL) -path ./database/migration down

.PHONY: migration-up
migration-up:
	migrate -database $(DB_URL) -path ./database/migration up

.PHONY: migration-force
migration-force:
	migrate -database $(DB_URL) -path ./database/migration force $(version)