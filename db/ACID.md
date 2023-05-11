# DB Transaction

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

## Composition rather than Inheritance

To extend a struct, [composition](https://www.geeksforgeeks.org/composition-in-golang/) is
a preferred way to inheritance.

Note: https://hackthology.com/object-oriented-inheritance-in-go.html

Composition doesn't extend a class itself, embbeds a struct and implements a functionality on it.

We extends `Queries` in `db.go` because they are single queries and doesn't support tx.

## Go routine

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

# DB Lock

## working with TDD

Test Driven Development.

Tests first, Improve it gradually.

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

#### postgres debugging statements

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
