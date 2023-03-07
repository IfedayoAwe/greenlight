start-db:
	docker run \
		--name greenlight_db \
		-p 5432:5432 \
		-v greenlight_db-data:/var/lib/postgresql/data \
		-e POSTGRES_USER=postgres \
		-e POSTGRES_PASSWORD=654321 \
		-e POSTGRES_DB=greenlight \
		-d postgres:15-alpine

createdb:
	docker exec -it greenlight_db createdb --username=postgres --owner=postgres greenlight

dropdb:
	docker exec -it greenlight_db dropdb greenlight

migrateup:
    migrate -path ./migration -database "postgresql://postgres:654321@localhost:5432/greenlight?sslmode=disable" -verbose up

migratedown:
    migrate -path ./migration -database "postgresql://postgres:654321t@localhost:5432/greenlight?sslmode=disable" -verbose down

.PHONY: start-db createdb dropdb migrateup migratedown
