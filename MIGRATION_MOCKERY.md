# Mockery

Gomock is no longer maintained by [google](https://github.com/golang/mock) and now by [uber](https://github.com/uber-go/mock).

Another choice is [mockery](https://vektra.github.io/mockery/latest/).

## Installation

Follow the official [docs](https://vektra.github.io/mockery/latest/installation/).

## Generate mock

```bash
$ mockery --all -r --with-expecter
```

## Supported Syntax

Refer to https://vektra.github.io/mockery/latest/examples/ and https://vektra.github.io/mockery/latest/features/ for usage.

Basically to build stub,

```golang
client.On("<Function name>", ...<Arguments>).Return(...Returned value).Times(1)
```

Expect feature like gomock need to configure with `--with-expecter`.

```golang
store.EXPECT().
    GetAccount(mock.Anything, account.ID).
    Times(1).
    Return(account, nil)
```
