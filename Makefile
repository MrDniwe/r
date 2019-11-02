setup.development: migrate.install start.dev-db migrate.development-db
	@echo Development setup completed

migrate.install:
	@echo Installing migrator...
	go get -v github.com/rubenv/sql-migrate/...

start.dev-db:
	@echo Starting development database...
	docker-compose up -d

migrate.development-db:
	@echo applying development migration...
	sql-migrate up -env="development"
