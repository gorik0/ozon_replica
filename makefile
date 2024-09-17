init_db:

run_migration:
	 migrate -database postgres://goirk:gorik@localhost:5432/postgres?sslmode=disable -path ./db/migrations -verbose  up



grpc_gen:
	protoc --go_out=. --go_opt=paths=source_relative proto/gmodels/*.proto
	protoc --proto_path=proto --go_out=. --go-grpc_out=. proto/*.proto


