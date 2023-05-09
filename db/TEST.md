# part 5

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
