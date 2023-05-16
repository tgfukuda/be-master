# Authentication

We must provides a authentication feature for the bank project.

First of all, we defines Users table to identify each user.

## Migration file: Up

See [migration file](../db/migration/000002_add_user.up.sql).

dbdiagram helps us generate these schema, but skip this time for the important.

Users tables of

```sql
CREATE TABLE "users" (
    "username" varchar PRIMARY KEY,
    "hashed_password" varchar NOT NULL, -- we MUST NOT store passwords directly. generally with solt
    "full_name" varchar NOT NULL,
    "email" varchar UNIQUE NOT NULL, -- to communicate with user
    "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',  -- it's recommended to rotate password in a certain period. default is very past time for ease develop.
    "created_at" timestamptz NOT NULL DEFAULT (now())
);
```

have a constraint between `accounts` of

```sql
ALTER TABLE "accounts" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");   -- acccounts.owner *->1 user.username
```

It means that there must be one `username` of `user` for each `owner` of `accounts`.

Additionally, a `user` can hold `account`s but their currency should differ.
To express that in db context, There're two ways.

First is creating a unique index.

```sql
CREATE UNIQUE INDEX ON "accounts" ("owner", "currency"); -- 
```

Second is adding a unique constraint.
Basically, it will automatically create **the same unique composite index as above**.

```sql
ALTER TABLE "accounts" ADD CONSTRAINT "owner_currency_key" UNIQUE ("owner", "currency");
```

However this migration will fail if some accounts entry already exist.
It will make the db dirty state, and down migration will be rejected.
We need to fix it by hand in such situation.

```
simple_bank=# select * from schema_migrations;          
 version | dirty 
---------+-------
       2 | t
(1 row)

simple_bank=# update schema_migrations set dirty = false;
UPDATE 1
simple_bank=# select * from schema_migrations;
 version | dirty 
---------+-------
       2 | f
(1 row)
```

Then we can `make migratedown` (the data will be lost) and now we can `make migrateup`.

## Migration file: Down

See [migration file](../db/migration/000002_add_user.down.sql).

We need to reverse operation of up.

First, drop constraint

```sql
ALTER TABLE IF EXISTS "accounts" DROP CONSTRAINT IF EXISTS "owner_currency_key"; -- defined in up
```

Second, `accounts` is

```
simple_bank=# \d accounts
                                       Table "public.accounts"
   Column   |           Type           | Collation | Nullable |               Default                
------------+--------------------------+-----------+----------+--------------------------------------
 id         | bigint                   |           | not null | nextval('accounts_id_seq'::regclass)
 owner      | character varying        |           | not null | 
 balance    | bigint                   |           | not null | 
 currency   | character varying        |           | not null | 
 created_at | timestamp with time zone |           | not null | now()
Indexes:
    "accounts_pkey" PRIMARY KEY, btree (id)
    "owner_currency_key" UNIQUE CONSTRAINT, btree (owner, currency)
    "accounts_owner_idx" btree (owner)
Foreign-key constraints:
    "accounts_owner_fkey" FOREIGN KEY (owner) REFERENCES users(username)
Referenced by:
    TABLE "entries" CONSTRAINT "entries_account_id_fkey" FOREIGN KEY (account_id) REFERENCES accounts(id)
    TABLE "transfers" CONSTRAINT "transfers_from_account_id_fkey" FOREIGN KEY (from_account_id) REFERENCES accounts(id)
    TABLE "transfers" CONSTRAINT "transfers_to_account_id_fkey" FOREIGN KEY (to_account_id) REFERENCES accounts(id)
```

have a constraint related to `user` of `"accounts_owner_fkey"` and so we need to drop it.

```sql
ALTER TABLE IF EXISTS "accounts" DROP CONSTRAINT IF EXISTS "accounts_owner_fkey";
```

Then, simply drop the table `users`.
