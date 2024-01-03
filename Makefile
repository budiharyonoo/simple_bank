# Init Postgres in Docker
postgres:
	docker run --restart unless-stopped --name postgres16 -d -e POSTGRES_PASSWORD=root -e POSTGRES_USER=root -p 5432:5432 postgres:16-alpine

# Create Database in Postgres
createdb:
	docker exec -it postgres16 createdb --username=root simple_bank

# Drop Database in Postgres
dropdb:
	docker exec -it postgres16 dropdb --username=root --force simple_bank

# Run migration up
migrateup:
	migrate -path db/migrations -database "postgresql://root:root@localhost:5432/simple_bank?sslmode=disable" -verbose up

# Run migration down
migratedown:
	migrate -path db/migrations -database "postgresql://root:root@localhost:5432/simple_bank?sslmode=disable" -verbose down

# Run all test file
test:
	go test -v -cover ./...

# Run development server
server:
	go run main.go

# Run SQLC generate
sqlc:
	sqlc generate

# Run mockgen
mockdb:
	mockgen -package mockdb -destination db/mock/store.go github.com/budiharyonoo/simple_bank/db/sqlc Store

.PHONY: postgres createdb dropdb migrateup migratedown test server sqlc mockdb
