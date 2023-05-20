# BE master class

My backend educational repository based on [this course](https://github.com/techschool/simplebank).
Some addtional contents contains.
Each directory contains my summary of this course.

## Requirements
- docker: 23.0.4
- sqlc: v1.18.0
- migrate: dev

### setup postgres image

```
make postgres
```

### migration

```
make createdb # if needed.
make migrateup
```

### gen sql orm client

```
make sqlc
```

## Summaries

- [DB schema management](./db/README.md)
- [Testing](./db//TEST.md)
- [CRUD operations](./db/CRUD.md)
- [ACID property of DB](./db/ACID.md)
- [Web framework: Gin](./api/README.md)
- [Token Based Authentication](./api/AUTH.md)
- [Middlewares](./api/MIDLLEWARE.md)
