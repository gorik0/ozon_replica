init_db:

run_migration:
	 migrate -database postgres://goirk:gorik@localhost:5432/postgres?sslmode=disable -path ./db/migrations -verbose  up
