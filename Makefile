include ./.env
DBURL=postgres://$(DBUSER):$(DBPASS)@$(DBHOST):$(DBPORT)/$(DBNAME)?sslmode=disable
MIGRATIONPATH=db/migrations

migrate-create:
	migrate create -ext sql -dir $(MIGRATIONPATH) -seq create_$(NAME)_table

migrate-createUp:
	migrate -database $(DBURL) -path $(MIGRATIONPATH) up

migrate-createDown:
	migrate -database $(DBURL) -path $(MIGRATIONPATH) down $(s)

migrate-status:
	migrate -database $(DBURL) -path $(MIGRATIONPATH) version

migrate-force:
	migrate -database $(DBURL) -path $(MIGRATIONPATH) force $(v)