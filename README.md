# BE master class

## Requirements
- docker: 23.0.4
- sqlc: v1.18.0
- migrate: dev

## setup postgres image

```
make postgres
```

## migration

```
make createdb # if needed.
make migrateup
```

## gen sql orm client

```
make sqlc
```