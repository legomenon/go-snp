postgres:
	docker run --name goSnp -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres

rundb:
	docker start goSnp

createdb:
	docker exec -it goSnp  createdb --username=root --owner=root goSnpDB 

dropdb:
	docker exec -it goSnp  dropdb goSnpDB

# migrate create -seq -ext=.sql -dir=./migrations create_snippet_table
# migrate create -seq -ext=.sql -dir=./migrations create_sessions_table

migrateup:
	migrate -path=./migrations -database=postgres://root:secret@localhost/goSnpDB?sslmode=disable up

migratedown:
	migrate -path=./migrations -database=postgres://root:secret@localhost/goSnpDB?sslmode=disable down
