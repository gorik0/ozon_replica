init_db:

run_migration:
	 migrate -database postgres://goirk:gorik@localhost:5432/postgres?sslmode=disable -path ./db/migrations -verbose  up



grpc_gen:
	protoc --go_out=. --go_opt=paths=source_relative proto/gmodels/*.proto
	protoc --proto_path=proto --go_out=. --go-grpc_out=. proto/*.proto



run_auth:
	AUTH_JWT_SECRET_KEY=2132 \
  	CSRF_JWT_SECRET_KEY=12387 \
  	POSTGRES_DB=postgres \
  	POSTGRES_DB=postgres \
  	DB_PORT=5432 \
  	POSTGRES_PASSWORD=gorik \
  	POSTGRES_USER=goirk \
   	go run cmd/auth/main.go

run_products:
	AUTH_JWT_SECRET_KEY=2132 \
  	CSRF_JWT_SECRET_KEY=12387 \
  	POSTGRES_DB=postgres \
  	POSTGRES_DB=postgres \
  	DB_PORT=5432 \
  	POSTGRES_PASSWORD=gorik \
  	POSTGRES_USER=goirk \
   	go run cmd/products/main.go

client_auth:
	go run client/client.go