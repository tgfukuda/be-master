# BE master class

My backend educational repository based on [this course](https://github.com/techschool/simplebank).
Some addtional contents contains.
Each directory contains my summary of this course.

## Requirements
- docker: 23.0.4
- sqlc: v1.18.0
- migrate: dev
- statik
- protoc
- mockgen

### Setup postgres image

```
make postgres
```

### Migration

Migration

```
make createdb # if needed.
make migrateup
```

Create New migration file

```
make new_migration name=<name>
```

### Generate sql orm client

```
make sqlc
```

## Summaries

- [DB schema management](./db/README.md)
- [Testing](./TEST.md)
- [CRUD operations](./db/CRUD.md)
- [ACID property of DB](./db/ACID.md)
- [Web framework: Gin](./api/README.md)
- [Token Based Authentication](./api/AUTH.md)
- [Middlewares](./api/MIDLLEWARE.md)
- [gRPC](./gapi/README.md)
- [protoc](./proto/README.md)
- [Infrastructure](./terraform/README.md)
- [Task and Background Worker](./worker/README.md)
- [Continuous Integration/Continuous Development](./CICD.md)
