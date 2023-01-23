# Start docker

```bash
docker pull postgres
docker run --name=go-todo-rest-api -e POSTGRES_PASSWORD='1425' -p 5432:5432 -d --rm postgres
```

# Migrate

```bash

migrate create -ext sql -dir ./schema -seq init
migrate -path ./schema -database 'postgres://postgres:1425@localhost:5432/postgres?sslmode=disable' up
```

## Dirty

```bash
docker exec -it [container] /bin/bash
psql -U postgres
select * from schema_migrations;
update schema_migrations set version='000001', dirty=false;

```
