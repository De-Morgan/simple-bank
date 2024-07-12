DB_URL=postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable
postgres:
	docker run --name --network bank-network postgres16 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:16.3-alpine3.20

createdb:
	docker exec -it postgres16 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres16 dropdb simple_bank

migrateup:
	migrate -path db/migration -database "${DB_URL}" --verbose up

migrateup1:
	migrate -path db/migration -database "${DB_URL}" --verbose up 1

migratedown:
	migrate -path db/migration -database "${DB_URL}" --verbose down

migratedown1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" --verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

mock:
	mockgen -package mock_db  -destination db/mock/store_mock.go github.com/morgan/simplebank/db/sqlc Store

server:
	go run main.go

docker_image_build:
	docker build -t simplebank:latest .

docker_create_container:
	docker run --name simplebank --network bank-network -p 8080:8080 -e GIN_MODE=release -e DB_SOURCE="postgresql://root:secret@postgres16:5432/simple_bank?sslmode=disable" simplebank:latest 


proto:
	rm -f pb/*.go
	rm -f doc/swagger/*.swagger.json
	protoc --proto_path=protos --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	 --grpc-gateway_out ./pb \
    --grpc-gateway_opt paths=source_relative \
	--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=simple_bank \
    protos/*.proto

evans:
	 evans --host localhost --port 50051 -r repl
	

redis:
	docker run --name redis -p 6379:6379 redis:7-alpine