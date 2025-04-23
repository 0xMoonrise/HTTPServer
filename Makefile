all:
	go run .

up:
	goose -dir sql/schema/ postgres $(DB_URL) up

down:
	goose -dir sql/schema/ postgres $(DB_URL) down

status:
	goose -dir sql/schema/ postgres $(DB_URL) status
