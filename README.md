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

# API 
| METHOD |              URI              | DESCRIPTION |
| ------ | ----------------------------- | ----------- |
| POST   | /auth/sign-up                 | signUp      |
| POST   | /auth/sign-in                 | signIn      |
| POST   | /api/lists/                   | createList  |
| GET    | /api/lists/                   | getAllLists |
| GET    | /api/lists/:id                | getListById |
| PUT    | /api/lists/:id                | updateList  |
| DELETE | /api/lists/:id                | deleteList  |
| POST   | /api/lists/:id/items/         | createItem  |
| GET    | /api/lists/:id/items/         | getAllItems |
| GET    | /api/lists/:id/items/:item_id | getItemById |
| PUT    | /api/lists/:id/items/:item_id | updateItem  |
| DELETE | /api/lists/:id/items/:item_id | deleteItem  |
| GET    | /api/items/:id                | getItemById |
| PUT    | /api/items/:id                | updateItem  |
| DELETE | /api/items/:id                | deleteItem  |
