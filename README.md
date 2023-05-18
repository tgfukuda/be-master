# BE master class

The original course found [here](https://github.com/techschool/simplebank).
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
