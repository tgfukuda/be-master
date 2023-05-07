# Part1: db schema

In the first lecture, I learned about how to design db schema and ER-diagram with [dbdiagram](https://dbdiagram.io/d).

It produces a MySQL, PostgreSQL, ER-diagram, and others easily from
[./dbdiagram.txt] and helps the generated ones with team members.

# Part2: run docker instance for database

In the second lecture, the teacher said run docker container with 
[postgre image](https://hub.docker.com/_/postgres/)

Note: password should be more complex one and
this is only for educational purpose.
```
$ sudo docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine
```

to execute a pre-defined sql file with commandline,

```
$ sudo docker cp db/generated/SimpleBank-postgre.sql postgre12:/tmp/
$ sudo docker exec -it postgre12 psql -U root -a -f /tmp/SimpleBank-postgre.sql
```

some tips

|MySql|Postgres|
|:-:|:-:|
| show databases | \l |
| show tables | \dt |
|desc <table_name>|select table_name,column_name,data_type from information_schema.columns where table_name='<table_name>'|

## GUI tools

[TablePlus](https://tableplus.com/)
