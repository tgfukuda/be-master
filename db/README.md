# DB schema

In the first lecture, I learned about how to design db schema and ER-diagram with [dbdiagram](https://dbdiagram.io/d).

It produces a MySQL, PostgreSQL, ER-diagram, and others easily from
[./dbdiagram.txt] and helps the generated ones with team members.

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
|desc <table_name>|select table_name,column_name,data_type from information_schema.columns where table_name='<table_name>'|

## GUI tools

[TablePlus](https://tableplus.com/)

# Migration

## installation

See, https://github.com/golang-migrate/migrate/tree/master/cmd/migrate
With go toolchain,
```
$ go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```
will install migrate at `$GOPATH/bin/migrate`.
(My linux failed to install with above instruction...)

## best practice

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
