# Go testing

go test file typically put the same directory as the target functionality and suffix with _test.

In ./sqlc, account_test.go means to test `account.sql.go`.
`main_test.go` organizes all test and has some helpers.

**Note**:
In main_test.go,
```
_ "github.com/lib/pq"
```
importing with name _ is special import to tell go not to remove this deps

## Recommended packages
- [pq](https://github.com/lib/pq)
- [testify](https://github.com/stretchr/testify#installation)

## init() and fuggy testing

```
// special function called first of all
func init() {
	rand.Seed(time.Now().UnixNano()) // If Seed is not called, the generator behaves as if seeded by Seed(1).
}

func RandInt(...) {
    ...
}
```

## Mocking

There can be two way to test functionalities related to db.
First is use db directly and Second is Mock db. It depends on a situation, but Mock seems to be better.

### Why mock database?

- Independent tests: to avoid conflicts
- Faster tests: reduce a lot of time talking with database
- 100% coverage: easily setup edge cases i.e. unexpected errors, a connection lost

### How to mock?

1. Use Fake DB - In memory: develop the same interface as real db.
    - Pros: very simple and easy
    - Cons: bothered by the codes to be used only for test.
2. Use DB Stub: GOMOCK

We use [mockgen](https://github.com/golang/mock) to mock.
Install as the Instruction.

Mockgen has two mode called Source mode and Reflect mode.
Source mode is more suitable for prroduction but we use Reflect for simplicity.

1. Declare the Store as an interface. sqlc interface generation is useful and embedd it with composition.
    ```go
    type Store interface {
        Querier
        TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
        SafeDeleteAccountTx(ctx context.Context, arg SafeDeleteAccountTxParams) (SafeDeleteAccountTxResult, error)
    }

    type SQLStore struct {
        *Queries
        db *sql.DB
    }
    ```
2. Create mock dir
3. generate mock structs with
    ```
    $ mockgen -package mockdb -destination db/mock/store.go github.com/tgfukuda/be-master/db/sqlc Store
    ```
    Note: There'll be the error like `prog.go:12:2: no required module provides package` and solved by https://github.com/golang/mock#debugging-errors.
    It may be better to use https://github.com/vektra/mockery.
4. There will be 2 main structs:
    ```go
    // MockStore is a mock of Store interface.
    type MockStore struct {
        ctrl     *gomock.Controller
        recorder *MockStoreMockRecorder
    }

    // MockStoreMockRecorder is the mock recorder for MockStore.
    type MockStoreMockRecorder struct {
        mock *MockStore
    }
    ```
    The idea is to specify how many times the interface should called and what arguments are.

