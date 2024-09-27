



run_migration_up:
	 migrate -database postgres://goirk:gorik@localhost:5432/postgres?sslmode=disable -path ./db/migrations -verbose  up

run_migration_down:
	 migrate -database postgres://goirk:gorik@localhost:5432/postgres?sslmode=disable -path ./db/migrations -verbose  down

run_migration_force:
	 migrate -database postgres://goirk:gorik@localhost:5432/postgres?sslmode=disable -path ./db/migrations -verbose  force $(V)



grpc_gen: #generate grpc stuff
	protoc --go_out=. --go_opt=paths=source_relative proto/gmodels/*.proto
	protoc --proto_path=proto --go_out=. --go-grpc_out=. proto/*.proto



run_auth:# run auth service (for local)
	AUTH_JWT_SECRET_KEY=a \
  	CSRF_JWT_SECRET_KEY=a \
  	POSTGRES_DB=postgres \
  	POSTGRES_DB=postgres \
  	DB_PORT=5432 \
  	POSTGRES_PASSWORD=gorik \
  	POSTGRES_USER=goirk \
   	go run cmd/auth/main.go

run_products:# run products service (for local)
	AUTH_JWT_SECRET_KEY=a \
  	CSRF_JWT_SECRET_KEY=a \
  	POSTGRES_DB=postgres \
  	POSTGRES_DB=postgres \
  	DB_PORT=5432 \
  	POSTGRES_PASSWORD=gorik \
  	POSTGRES_USER=goirk \
   	go run cmd/products/main.go

run_order: # run order service (for local)
	AUTH_JWT_SECRET_KEY=a \
  	CSRF_JWT_SECRET_KEY=a \
  	POSTGRES_DB=postgres \
  	POSTGRES_DB=postgres \
  	DB_PORT=5432 \
  	POSTGRES_PASSWORD=gorik \
  	POSTGRES_USER=goirk \
   	go run cmd/orders/main.go

run_main: # run main service (for local)
	AUTH_JWT_SECRET_KEY=a \
  	CSRF_JWT_SECRET_KEY=a \
  	POSTGRES_DB=postgres \
  	POSTGRES_DB=postgres \
  	DB_PORT=5432 \
  	POSTGRES_PASSWORD=gorik \
  	POSTGRES_USER=goirk \
   	go run cmd/main/main.go




client_auth: #Test grpc client for auth_service
	go run client/auth/client.go

client_products:  #Test grpc client for products_service
	go run client/products/client.go

client_orders:  #Test grpc client for orders_service
	go run client/order/client.go


run_all: run_products run_auth run_order run_main