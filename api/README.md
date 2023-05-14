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
    