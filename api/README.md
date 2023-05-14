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
