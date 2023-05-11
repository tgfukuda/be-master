# part5: DB Transaction

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

# part 6: DB Lock

## working with TDD

Test Driven Development.

Tests first, Improve it gradually.

## DB Lock
