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

## Add Queries

See [user.sql](../db/query/user.sql) and generated one.

There're some constraint and we need to modify [account.go](./account.go) like commit 9fad23b282b7419a85173759a2b71b6121288d29.
We can handle sql errors with more detail.

```go
if pqErr, ok := err.(*pq.Error); ok {
    switch pqErr.Code.Name() {
    case "foreign_key_violation", "unique_violation":
        ctx.JSON(http.StatusForbidden, errorResponse(err))
        return
    }
}
```

## Hashed Password

[Bcrypt](https://auth0.com/blog/hashing-in-action-understanding-bcrypt/) (Blowfish Cipher and Crypt for the password hash function of UNIX) often
used for password hashing management.

### Why Bcrypt?

- We must avoid to store naked password for clear security reasons
- There're some attacks like rainbow table now a days for the extreamly fast computation.

Bcrypt provides us a completely different output with [salt](https://en.wikipedia.org/wiki/Salt_(cryptography)) that used in hash iterations
even if the passwords themself are the same.

The format is as follows.

- First: Hash Algorism Identifier. 2A means bcrypt.
- Second: Cost means key expantion rounds. 10 means 2^10 = 1024 rounds.
- Third: 16bytes of salt used to calculate each iteration. (22 chars in [base64](https://base64.guru/learn/what-is-base64))
- Forth: 24 bytes of hash itself. (31 chars in base64)

```
$2A$10$NQO...16bytes...ADP4...24bytes...YW
 |  |          |                  |
ALG COST     SALT               HASH
```

The implementation in golang is https://github.com/golang/crypto/blob/master/bcrypt/bcrypt.go.

### Use bcrypt

Tips: According to https://stackoverflow.com/questions/61283248/format-errors-in-go-s-v-or-w, use `%w` for formatting error.

```
// returns the bcrypt hash of the password
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password %w", err)
	}
	return string(hashedPassword), nil
}

// returns the password is correct or not
func CheckPassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
```

Then we can store it like [user.go](./user.go).

## Resources

- https://bcrypt.online/
- https://en.wikipedia.org/wiki/Dictionary_attack
- https://base64.guru/learn/what-is-base64
