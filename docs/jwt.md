# Implement JWT Auth

## Implement JWT Login

- Add package using `go get github.com/bimalabs/middlewares`

- Add jwt login route to `dics/container.go`

```go
{
    Name: "bima:route:jwt:login",
    Scope: bima.Application,
    Build: func(env *configs.Env) (routes.Route, error) {
        return jwt.DefaultJwtLogin("/api/v1/login", env.Secret, jwt.SigningMethodHS512.Name, true, routes.FindUserByUsernameAndPassword(func(username, password string) jwt.MapClaims {
            return jwt.MapClaims{
                "user": "admin",
            }
        })), nil
    },
    Params: dingo.Params{
        "0": dingo.Service("bima:config"),
    },
}
```

You need to implmement `routes.FindUserByUsernameAndPassword` function with your own logic. If you don't implement refresh token, pass `false` to 4th argument.

- Add to `configs/routes.yaml`

```yaml
routes:
    - jwt:login
```

## Implement JWT Validator (Middleware)

- Add jwt middleware to `dics/container.go`

```go
{
    Name: "bima:middleware:jwt",
    Scope: bima.Application,
    Build: func(env *configs.Env) (middlewares.Middleware, error) {
        return jwt.NewJwt(env, jwt.SigningMethodHS512.Name, "/health$"), nil
    },
    Params: dingo.Params{
        "0": dingo.Service("bima:config"),
    },
}
```

- Add to `configs/middlewares.yaml`

```yaml
middlewares:
    - jwt
```

You can access user using `configs.Env.User` or via `request.Header.Get("X-Bima-User")`

## Implement Refresh  JWT

- Add refresh jwt route to `dics/container.go`

```go
{
    Name: "bima:route:jwt:refresh",
    Scope: bima.Application,
    Build: func(env *configs.Env) (routes.Route, error) {
        return jwt.NewJwtRefresh("/api/v1/token-refresh", env.Secret, jwt.SigningMethodHS512.Name, 730), nil
    },
    Params: dingo.Params{
        "0": dingo.Service("bima:config"),
    },
},
```

- Add to `configs/routes.yaml`

```yaml
routes:
    - jwt:login
    - jwt:refresh
```
