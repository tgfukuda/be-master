# Restful HTTP API

## Very bothered to implement it with `net/http` ...

We can implement it based on `net/http`, but there're many web frameworks.

- [Gin](https://github.com/gin-gonic/gin)
- [Beego](https://github.com/beego/beego)
- [Echo](https://echo.labstack.com/)
- [Revel](https://revel.github.io/)
- [Martini](https://github.com/go-martini/martini)
- [Fiber](https://github.com/gofiber/fiber)
- [Buffalo](https://gobuffalo.io/)

They'll offer the features like ...

- Routing
- Parameter Binding
- Validation
- Middleware
- Builtin ORM

If we need a routing part of it,

- [FastHttp](https://github.com/valyala/fasthttp)
- [GorillaMux](https://github.com/gorilla#gorilla-toolkit) - already archived and not actively maintaned
- [HttpRouter](https://github.com/julienschmidt/httprouter) - seems not to be maintained well
- [Chi](https://github.com/go-chi/chi)

## Gin: most popular framework

To have a quick reference, see [Examples](https://github.com/gin-gonic/examples)
and [Docs](https://github.com/gin-gonic/gin/blob/master/docs/doc.md).

There're many features and we use it by need.

### Validator notations

see https://github.com/gin-gonic/gin/blob/master/docs/doc.md#model-binding-and-validation, https://pkg.go.dev/github.com/go-playground/validator#hdr-Baked_In_Validators_and_Tags and https://github.com/gin-gonic/examples/blob/master/custom-validation/server.go

### Simple Request Handlers with Gin

There're three major request parameter handlers in `account.go`.
For more details, https://github.com/gin-gonic/gin/blob/master/docs/doc.md#api-examples 

1. Parameters with Body (JSON): We retrieve given parameters from request body. refer to `CreateAccount`.
    ```
    $ curl -X POST "http://localhost:8080/accounts" -d '{"owner":"tgfkd","currency":"USD"}'
    ```
2. Parameters in the path: We get them using path variables. refer to `GetAccount`.
    ```
    $ curl -X GET "http://localhost:8080/accounts/20"
    ```
3. Parameters with query string: We get them using query string (`/some/path?query=strings`). refer to `ListAccounts`.
    ```
    $ curl -X GET "http://localhost:8080/accounts?page_size=10&page_id=2"
    ```

## Configuration with Viper

[Viper](https://github.com/spf13/viper) provides some useful features for us.

1. Resolve config file in many formats: JSON, TOML, YAML, INI and so on.
2. Read config from environment variables and flags.
3. Read config from remote system: Etcd, Consul.
4. Live watching the file: like hot reload of webpack.

We can manage a lot of configuration in development and change it depending on the environment i.e. production, staging, dev, local...

### Basic Usage

See, `util/config.go` and the official docs.

```go
    viper.AddConfigPath(path)   // set the directory path to read values from
	viper.SetConfigName("app")  // we have app.env
	viper.SetConfigType("env")	// json, xml, ... it makes sure that the file follows the correct format and has the correct extension.

	viper.AutomaticEnv()    // override the values if there's any corresponding named env var.
```

## Test with GoMock

See also: ../db/TEST.md.

Implement mock testing in `account_test.go`.

To use the generated mock,
```go
// NewMockStore creates a new mock instance.
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}
```

receives `*gomock.Controller`, we need to setup it.

`*gomock.Controller` asserts tests with `Finish`, keep in mind to call `defer ctrl.Finish()`.

```go
// Finish checks to see if all the methods that were expected to be called
// were called. It should be invoked for each Controller. It is not idempotent
// and therefore can only be invoked once.
func (ctrl *Controller) Finish() { ... }
```

Make a Stub with the object

```go
// build stubs
store.EXPECT().
    GetAccount(gomock.Any(), gomock.Eq(account.ID)).	// expected call of GetAccount
    Times(1).	// How exactly many times?
    Return(account, nil)	// What should be returned?
```

Check API response with recorder instead of a real server.

```go
server := NewServer(store)
recorder := httptest.NewRecorder()

url := fmt.Sprintf("/accounts/%d", account.ID)
request, err := http.NewRequest(http.MethodGet, url, nil)
```

