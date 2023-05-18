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

There's a list of [popular frameworks repository](https://github.com/mingrammer/go-web-framework-stars).

## Gin: most popular framework

To have a quick reference, see [Examples](https://github.com/gin-gonic/examples)
and [Docs](https://github.com/gin-gonic/gin/blob/master/docs/doc.md).

There're many features and we use it by need.

### Validator notations

Golang has a [tag](https://zetcode.com/golang/struct-tag/) feature that enables us so many and Gin also uses tag.
Gin uses a for model binding and validation.

That looks like

```go
type CreateAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR JPY"`
}
```

See also https://github.com/gin-gonic/gin/blob/master/docs/doc.md#model-binding-and-validation, https://pkg.go.dev/github.com/go-playground/validator#hdr-Baked_In_Validators_and_Tags and https://github.com/gin-gonic/examples/blob/master/custom-validation/server.go.

It's very useful but how about the hardcoded values of `USD EUR JPY`?

#### Custom Validator

To utilize validator, we need to implement a validator function of type Func exported by `"github.com/go-playground/validator/v10"`.

(See [validator.go](./validator.go))
Register the function with

```go
if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
    v.RegisterValidation("currency", validCurrency)
}
```

and now we can use `currency` tag like `binding:"required,currency"` in Gin.

### Simple Request Handlers with Gin

There're three major request parameter handlers in [account.go](./account.go).
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

Implement mock testing in [account_test.go](./account_test.go).

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

### Gomock matcher

In the stub object,

```go
store.EXPECT().
    CreateUser(gomock.Any(), gomock.Eq(arg)).
    Times(1).
    Return(user, nil)
```

we use `gomock.Eq(...)` to specify an input with

```go
// Eq returns a matcher that matches on equality.
//
// Example usage:
//   Eq(5).Matches(5) // returns true
//   Eq(5).Matches(4) // returns false
func Eq(x interface{}) Matcher { return eqMatcher{x} }
```

There're several builtin matchers, see https://github.com/golang/mock/blob/main/gomock/matchers.go.
However, in the cases like [bcrypt auth](./AUTH.md), bcrypt always return a different output even with the same password.

In such cases, we can't guess a returned object in advance, but

```go
    CreateUser(gomock.Any(), gomock.Any()).
    Times(1).
    Return(user, nil)
```

makes the test very weak.

### Gomock Custom matcher

The matcher interface is

```go
type Matcher interface {
	// Matches returns whether x is a match.
	Matches(x interface{}) bool

	// String describes what the matcher matches.
	String() string
}
```

and we can define as

```go
type eqCreateUserParamsMatcher struct {
	arg      db.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}

	err := util.CheckPassword(e.password, arg.HashedPassword)
	if err != nil {
		return false
	}

	e.arg.HashedPassword = arg.HashedPassword
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("mathers arg %v and password %v", e.arg, e.password)
}

func EqCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg: arg, password: password}
}
```

## References

- https://medium.com/golangspec/tags-in-golang-3e5db0b8ef3e
- https://pkg.go.dev/reflect
- https://github.com/mingrammer/go-web-framework-stars
