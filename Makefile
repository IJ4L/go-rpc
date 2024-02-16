postgres:
	docker run --name postgres-test --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

createdb:
	docker exec -it postgres-test createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres-test dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@127.0.0.1:5432/simple_bank?sslmode=disable" -verbose up

migrateup1:
	migrate -path db/migration -database "postgresql://root:ht%3E%21%2Aa%3FSsADeS4H7ug%7DHd~scS5Ef@simple-bank.crimymiiqkv0.ap-southeast-1.rds.amazonaws.com:5432/simple_bank" -verbose up 1

migratedown:
	migrate -path db/migration -database "postgresql://root:ht%3E%21%2Aa%3FSsADeS4H7ug%7DHd~scS5Ef@simple-bank.crimymiiqkv0.ap-southeast-1.rds.amazonaws.com:5432/simple_bank" -verbose down

migratedown1:
	migrate -path db/migration -database "postgresql://root:ht%3E%21%2Aa%3FSsADeS4H7ug%7DHd~scS5Ef@simple-bank.crimymiiqkv0.ap-southeast-1.rds.amazonaws.com:5432/simple_bank" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go simplebank.com/db/sqlgen Store

evans:
	evans -r -p 9090

.PHONY: postgres createdb dropdb migrateup migrateup1 migratedown migratedown1 sqlc test server mock