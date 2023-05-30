# DB schema

In the first lecture, I learned about how to design db schema and ER-diagram with [dbdiagram](https://dbdiagram.io/d).

It produces a MySQL, PostgreSQL, ER-diagram, and others easily from
[dbdiagram.txt](./dbdiagram.txt) and helps the generated ones with team members.

`dbml2sql` helps us generate Schema on local. Install `npm i -g @dbml/cli` and `make db_schema`.

# Run docker instance for database

In the second lecture, the teacher said run docker container with 
[postgre image](https://hub.docker.com/_/postgres/)

To execute a pre-defined sql file with commandline,

```
$ sudo docker cp db/generated/SimpleBank-postgre.sql postgre12:/tmp/
$ sudo docker exec -it postgres12 psql -U root -a -f /tmp/SimpleBank-postgre.sql
```

Some tips

|MySql|Postgres|
|:-:|:-:|
| show databases | \l |
| show tables | \dt |
|desc <table_name>| \d <table_name>|

In postgres, `select table_name,column_name,data_type from information_schema.columns where table_name='<table_name>';`
also show us the table details.

## GUI tools

[TablePlus](https://tableplus.com/)

# Migration

## Installation

See, https://github.com/golang-migrate/migrate/tree/master/cmd/migrate
With go toolchain,
```
$ go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```
will install migrate at `$GOPATH/bin/migrate`.
(My linux failed to install with above instruction...)

## Initialization

We can use migrate to generate empty migration file with

```
$ migrate create -ext sql -dir db/migration/ -seq <desc_of_migrate>
```

## Best practice

There should be 2 types of sql file for migration.

First is `.up.sql`. The file modify db schema to go ahead.

To get the latest version of schema, apply
```
000001_xxxx.up.sql -> 000002_xxxx.up.sql -> ...
```

Second is `.down.sql` and it's the contrary.

Therefore, `000nnn_xxx_yyy.up.sql` should be the reverse transformation
against `000nnn_xxx_yyy.down.sql`.

Note: Need to write migrations manually.

## Migration in code

https://github.com/golang-migrate/migrate#use-in-your-go-project.

# Flow of add new table to schema

1. Run `migrate create -ext sql -dir db/migration/ -seq <desc_of_migrate>` to add empty [up and down script](./migration/). Make sure the migration can be applied by `make migrateup` (`make postgres` and `make createdb` if needed).
2. Add queries inside [queries](./query/).
3. Run `make sqlc` for [go bindings](./sqlc/)
4. Run `make mock` to add [mock](./mock/). Make sure `make test` passes.
5. Do work with the table.

# Documantation

[DBDoc](https://dbdocs.io/) powerfully supports documentation.

Docs generated from dbml, we can manage and visualize the schema with a dbml file and `dbdocs build docs/db.dbml`.

```bash
$ dbdocs login # login dbdocs. use email or github.
$ dbdocs build docs/db.dbml # build documentation
$ dbdocs password --set <password> --project <prj_name> # set password to the built schema doc
$ dbdocs remove <prj_name> # remove a project
```
