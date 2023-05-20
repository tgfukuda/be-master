# Middleware

## Overview

middleware plays a crucial role in handling the intermediate processing of requests and responses between the client and the server.
It acts as a bridge or a layer between the application's backend logic and the incoming requests or outgoing responses.
Middleware functions are designed to perform specific tasks such as request parsing, authentication, authorization, logging, error handling, and more.

```mermaid
graph TD
  A(Client) --> B(Middleware 1)
  B --> C(Application)
  C --> D(Middleware 2)
  D --> E(Server)

  subgraph Web Application Stack
    B --> F(Request)
    C --> G(Response)
  end

  style A fill:#98FB98, stroke:#000000, stroke-width:2px;
  style E fill:#98FB98, stroke:#000000, stroke-width:2px;
```

Note: this is a simplified representation, and in real-world scenarios, there might be multiple middleware components involved in the request-response flow.

## implement with gin

See [middleware](./middleware.go).

Gin provides 

```go
// HandlerFunc defines the handler used by gin middleware as return value.
type HandlerFunc func(*Context)
```

to declare middleware and we use it.

To get a specific key and abort the request if it shouldn't process,

```go
authHeader := ctx.GetHeader(authorizationHeaderKey) // get the header
if len(authHeader) == 0 {
    err := errors.New("authorization header is not provided")
    ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))    // abort the response and won't reach app logic
    return
}

//     :
// do something
//     :
ctx.Set(authorizationPayloadKey, payload) // provide the payload via context
ctx.Next()  // do next job
```

We need to call `ctx.Next` in the end of middleware.

```go
// Next should be used only inside middleware. It executes the pending handlers in the chain inside the calling handler. See example in GitHub.
func (*gin.Context).Next()
```

Note:
Authentication is only responsible for authenticating the request with a given access token.
It doesn't care how does each API do with it (who owns the token, what are operations allowed ...).

## Authorization type

Basically, authorization type of

- No Auth
- API Key
- Bearer Token
- Basic Auth
- Digest Auth
- OAUTH 1.0/2.0
- Hawk Authentication
- AWS signature
- NTLM Authentication
- Akamai EdgeGrid

and so on will compose the app authentication header like

The format of authentication in a request header typically follows a standardized format known as the "Authorization" header.
The Authorization header is used to include authentication credentials or tokens in an HTTP request.
There are different authentication schemes or protocols that define the specific format and usage of the Authorization header.
The two common authentication schemes used in request headers are:

Basic Authentication:

```json
{ "authorization": "Basic base64encodedCredentials" }
```

Bearer Token Authentication:

```json
{ "authorization": "Bearer tokenValue" }
```

## Add Middleware with gin

```go
authRoutes := router.
    Group("/"). // prepare router group
    Use(authMiddleWare(server.tokenMaker))  // common middleware

// add routes to the router group
authRoutes.POST("/accounts", server.CreateAccount)
authRoutes.GET("/accounts/:id", server.GetAccount)
authRoutes.GET("/accounts", server.ListAccounts)
authRoutes.POST("/delete_account", server.DeleteAccount)
authRoutes.POST("/transfers", server.CreateTransfer)
```
