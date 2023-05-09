# Part 4: CRUD
Create, Read, Update, Delete

Two approaches exists.

## With std lib

With low level standard packages like context, database/sql, log, time, ...etc.

- Pros: High performance.
- Cons: Quite boring (bind sql fields to variables, write raw sql) and easy to make mistake.

## With high level ORM(Object Relation Mapping)

Using gorm.

- Pros: All basic CRUD is already implemented and code will be very short.
- Cons: When network is conjucted, 3~5 times lower than std.

## Middle way approach: sqlx

Using sqlx library.

- Pros: As fast as std lib and binding to variable is already done.
- Cons: Code will be long and can cause mistake. Additionally, any error catched at runtime.

# Best way: sqlc

Using sqlc - metaprogramming of golang

- Pros: As fast as database/sql and easy to use. Automatic code gen. Catch sql error before code gen.
- Cons: only for postgres, mysql is experimental. (now, it seems to supports postgres, mysql, sqlite. see, https://docs.sqlc.dev/en/stable/tutorials/getting-started-mysql.html)

## sqlc installation

refer to https://github.com/kyleconroy/sqlc.

```
go install github.com/kyleconroy/sqlc/cmd/sqlc@latest
```

## sqlc configuration

There're 2 types of [configuration format](https://docs.sqlc.dev/en/stable/reference/config.html) now.

For version 1,

```yaml
version: "1"
packages:
  - name: "db"  # name of generated go package name
    path: "./db/sqlc" # path of generated package
    queries: "./db/query/" # referenced queruy dir
    schema: "./db/migration/" # schema def 
    engine: "postgresql" # postgres, mysql, ...
    emit_prepared_queries: false # for optimization.
    emit_interface: false # querier interfaces. useful when mock for tests.
    emit_exact_table_names: false # if false, sqlc make a struct name singular form. (table accounts -> struct Account)
    emit_empty_slices: false
    emit_exported_queries: false
    emit_json_tags: true # add json tag to struct
    emit_result_struct_pointers: false
    emit_params_struct_pointers: false
    emit_methods_with_db_argument: false
    emit_pointers_for_null_types: false
    emit_enum_valid_method: false
    emit_all_enum_values: false
    json_tags_case_style: "camel"
    output_batch_file_name: "batch.go"
    output_db_file_name: "db.go"
    output_models_file_name: "models.go"
    output_querier_file_name: "querier.go"
```

## Sqlc queries

Refer to [query reference](https://docs.sqlc.dev/en/stable/reference/query-annotations.html).
Meta programming with sql.
```sql
        name of func    returns single object
            V           V
-- name: CreateAccount :one
INSERT INTO authors (
  owner,
  balance,
  currency
) VALUES (
  $1, $2, $3
)
RETURNING *; <- tell postgres to return the values of all column after creation.
```

Ref: table of accounts
```sql
CREATE TABLE "accounts" (
  "id" bigserial PRIMARY KEY,
  "owner" varchar NOT NULL,
  "currency" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);
```

## Generated codes

If these're some mistakes, sqlc failed to compile queries!

1. models.go - golang struct binding of sql row.
2. db.go - db client initialized with sql.db or sql.tx.
3. queries - generated with the queries.

For more details, refer to document and ../db/query/.
