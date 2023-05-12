# ACID

## DB Transaction

DB transaction is a set of db operation.

It holds *ACID property*.

- Atomicity: either all operation complete successfully or the overall transaction fails (db unchanged).
- Consistency: the db state stayed valid after the transaction i.e. holds predefined rules, constraints, cascades, triggers.
- Isolation: concurrent transaction doesn't affect each other. (there're some levels)
- Durability: each transaction change is recorded persistently even in system failure.

Successful transaction will be
```
BEGIN; // tx started
:
COMMIT; // apply change
```

Failed transaction will be
```
BEGIN;
:
ROLLBACK; // abort transaction
```

### Composition rather than Inheritance

To extend a struct, [composition](https://www.geeksforgeeks.org/composition-in-golang/) is
a preferred way to inheritance in golang.

Note: https://hackthology.com/object-oriented-inheritance-in-go.html

Composition doesn't extend a class itself, embbeds a struct and implements a functionality on it.

We extends `Queries` in `db.go` because they are single queries and doesn't support tx.

### Go routine

To write db is easy and concurrency also easily causes bug unless handling concurrency carefully.

[go routine](https://golangbot.com/goroutines/) is a thread of go.

### Channels

[Channel](https://gobyexample.com/channels) is a way to connect concurrent thread
without explicit locking.

It provides us to share data with each channel.

```
c := make(chan ...)

go func() {

    c <- data // send some data to the main thread
}() // run go routine

d := <-c // receive something via channel

```

## DB Lock

If we get some rows without lock in a transaction concurrently,
it will result in an unexpected behavior.

The reason is that postgres returns the values
even on the other process updating some columns (though it may not complete).

To avoid it, we should [lock](https://www.postgresql.org/docs/current/explicit-locking.html)
the db with `FOR UPDATE` statement in postgres. https://stackoverflow.com/questions/18879584/postgres-select-for-update-in-functions

### Debugging deadlock

when `pq: deadlock detected` appears in some tests, we can use `context.WithValue` to debug it.

According to `context.WithValue`,

> WithValue returns a copy of parent in which the value associated with key is val.
>
> Use context Values only for request-scoped data that transits processes and APIs, not for passing optional parameters to functions.
>
> The provided key must be comparable and should not be of type string or any other built-in type to avoid collisions between packages using context. Users of WithValue should define their own types for keys. To avoid allocating when assigning to an interface{}, context keys often have concrete type struct{}. Alternatively, exported context key variables' static type should be a pointer or interface.

We defined
```go
var txKey = struct{}{}
```
in sqlc/store.go for this purpose.

To get it,
```go
value := ctx.Value(txKey)
```

The context lib of go provides what the ordering of transaction operations to us,
but still doesn't what happens in postgres.

### Postgres debugging statements

See https://wiki.postgresql.org/wiki/Lock_Monitoring for how to get it.

```sql
SELECT blocked_locks.pid     AS blocked_pid,
         blocked_activity.usename  AS blocked_user,
         blocking_locks.pid     AS blocking_pid,
         blocking_activity.usename AS blocking_user,
         blocked_activity.query    AS blocked_statement,
         blocking_activity.query   AS current_statement_in_blocking_process
   FROM  pg_catalog.pg_locks         blocked_locks
    JOIN pg_catalog.pg_stat_activity blocked_activity  ON blocked_activity.pid = blocked_locks.pid
    JOIN pg_catalog.pg_locks         blocking_locks 
        ON blocking_locks.locktype = blocked_locks.locktype
        AND blocking_locks.database IS NOT DISTINCT FROM blocked_locks.database
        AND blocking_locks.relation IS NOT DISTINCT FROM blocked_locks.relation
        AND blocking_locks.page IS NOT DISTINCT FROM blocked_locks.page
        AND blocking_locks.tuple IS NOT DISTINCT FROM blocked_locks.tuple
        AND blocking_locks.virtualxid IS NOT DISTINCT FROM blocked_locks.virtualxid
        AND blocking_locks.transactionid IS NOT DISTINCT FROM blocked_locks.transactionid
        AND blocking_locks.classid IS NOT DISTINCT FROM blocked_locks.classid
        AND blocking_locks.objid IS NOT DISTINCT FROM blocked_locks.objid
        AND blocking_locks.objsubid IS NOT DISTINCT FROM blocked_locks.objsubid
        AND blocking_locks.pid != blocked_locks.pid

    JOIN pg_catalog.pg_stat_activity blocking_activity ON blocking_activity.pid = blocking_locks.pid
   WHERE NOT blocked_locks.granted;
```

shows the locked and what locks and

```sql
SELECT a.datname,
        a.application_name,
         l.relation::regclass,
         l.transactionid,
         l.mode,
         l.locktype,
         l.GRANTED,
         a.usename,
         a.query,
         a.query_start,
         age(now(), a.query_start) AS "age",
         a.pid
FROM pg_stat_activity a
JOIN pg_locks l ON l.pid = a.pid
ORDER BY a.query_start;
```

shows which query causes a lock and whether it's granted.

There can be many reasons to cause deadlock, For example,

```sql
ALTER TABLE "transfers" ADD FOREIGN KEY ("from_account_id") REFERENCES "accounts" ("id");
```

has a foreign key constraint and any select statement of transfer get a lock.
If There's `FOR UPDATE` statements of `accounts`, it's deadlock.

### More precisely lock

A simple `FOR UPDATE` can be blocked by any related changes or reference like a foreign key constraint.
To tell we won't update such key or primary key, `FOR NO KEY UPDATE`.

### Another deadlock

Update like

```sql
BEGIN;
UPDATE accounts SET balance = ... WHERE id = 1;
UPDATE accounts SET balance = ... WHERE id = 2;
COMMIT;

BEGIN;
UPDATE accounts SET balance = ... WHERE id = 2;
UPDATE accounts SET balance = ... WHERE id = 1;
COMMIT;
```

can also cause a deadlock.

In this case, where `t1` is the first and `t2` is the second,

```
update id = 1 <- t1
update id = 2 <- t2
update id = 2 <- t1 -- cause deadlock!
```

because of the same Id of t2.
To fix it, we need to update accounts with ordering id and

```sql
BEGIN;
UPDATE accounts SET balance = ... WHERE id = 1;
UPDATE accounts SET balance = ... WHERE id = 2;
COMMIT;

BEGIN;
UPDATE accounts SET balance = ... WHERE id = 1;
UPDATE accounts SET balance = ... WHERE id = 2;
COMMIT;
```

The idea is making operations in the same order of keys no matter what we do.

## Isolation

Dive Deeper into Isolation part of ACID.

### Read Phenomena

There're something called *Read Phenomena*.
It is the interference between transactions running at the same time in a low level of Isolation.

- Dirty Read: A transaction reads data written by other concurrent **uncommitted** transaction.
- Non-Repeatable Read: A transaction reads the same row **twice**, and sees different value because other **committed** transaction modified it.
- Phantom Read: Like Non-Repeatable Read but it is a re-execution of reading multiple rows due to changes of insert or delete by other **committed** transaction.
- Serialization Anomaly: The result of a **group** of concurrent **commited** transactions is **impossible to achieve** if we try to run them sequentially in any order without overlapping.

### 4 standard isolation levels

[ANSI](https://ansi.org/) (American National Standards Institution) defines 4 levels of isolation.

The smaller indexed, The lower leveled.

1. Read Uncommited: Can see data written by uncommited transactions.
2. Read Committed: Only see data of committed transactions. (Dirty Read cannot occurs)
3. Repeatable Read: Same read query always returns same result no matter how may times executed and in the case some concurrent transaction has been committed.
4. Serializable: Can Achieve same result if execute transactions serially in some order instead of overlapping.

To work with isolation level,

mysql

```sql
select @@transaction_isolation; -- in this session. default will be Repeatable Read
select @@global.transaction_isolation; -- in the global
set session transaction isolation level read committed; -- to change unread committed, for example.
```

postgres

```sql
show transaction isolation level; -- Read Committed will be default value.
begin;  -- postgres allowed to set the transaction level for a specific 1 transaction, to do globally, https://stackoverflow.com/questions/62649971/how-to-change-transaction-isolation-level-globally.
set transaction isolation level read uncommitted;   -- set read uncommitted only for this transaction.
:
```

### What can Happens in each level?

o means "possible" and x means "impossible" in the below tables.

#### MySQL

|mysql|Dirty Read|Non-Repeatable Read|Phantom Read|Serialization Anomaly|
|:-:|:-:|:-:|:-:|:-:|
|Read Uncommited|o|o|o|o|
|Read Committed|x|o|o|o|
|Repeatable Read|x|x only for read but o|x only for read but o|o|
|Serializable|x|x|x|x|

Mysql's Repeatable Read prevents read-only queries from other transactions updates,
but once the transaction write anything, it will be affected.

Where `t1` and `t2` are transactions on Repeatable Read,

```sql
select balance from a where id = 1; <- t1 -- (1) assume 100
update a set balance = balance - 10 where id = 1; <- t2 -- updated by t2 and it've commited. it became 90
select balance from id = 1; <- t1 -- (2) same as (1) and 100
update a set balance = balance - 10 where id = 1; <- t1 -- updated by t1, it becomes 80
```

Mysql's Serializable automatically make all `SELECT` into `SELECT FOR SHARE` and it allows other transaction only to read.
By the lock, the above update cannot happens.
The second update query must wait for `t1`'s release of the lock.
However, if `t1` trys to query the forth update statement while `t2` waiting for `t1` complete,
it causes deadlock and `t1` will fail and `t2` will succeed.

For this reason, we must take care about
not only a timeout of wait but also deadlock in Serializable level.

Mysql handles Read Phenomena with *locking mechanism*.
If two transaction try to write simultaneously,
the later one wait for the other if its possible and
if impossible, the later one faces with dead lock error.

#### Postgres

|postgres|Dirty Read|Non-Repeatable Read|Phantom Read|Serialization Anomaly|
|:-:|:-:|:-:|:-:|:-:|
|Read Uncommited|x|o|o|o|
|Read Committed|x|o|o|o|
|Repeatable Read|x|x|x|o|
|Serializable|x|x|x|x|

Postgres allowed to set read uncommitted level but *it behaves like read committed*
because there's no circumstance to use read uncommitted, postgres said.

The lowest level of postgres is *Read Committed*.

See, https://www.postgresql.org/docs/current/transaction-iso.html.

Postgres adopts *dependency detection* to handle Read Phenomena and if it can happens,
the later transaction **fails** with hints like
> HINT: The transaction might succeed if retried.

### What should we do?

- Retry mechanism: There might be errors, timeout or deadlock.
- Read docs carefully: each database might implement isolation differently.

## Refereces

https://www.postgresql.org/docs/current/transaction-iso.html

https://dev.mysql.com/doc/refman/8.0/en/innodb-transaction-isolation-levels.html

https://youtu.be/4EajrPgJAk0
