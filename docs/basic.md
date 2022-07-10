# Implement Basic Auth

- Add package using `go get github.com/bimalabs/middlewares`

- Add basic auth middleware to `dics/container.go`

```go
{
    Name:  "bima:middleware:basic-auth",
    Scope: bima.Application,
    Build: func(validator basic_auth.ValidateUsernameAndPassword) (middlewares.Middleware, error) {
        return basic_auth.New(validator), nil
    },
    Params: dingo.Params{
        "0": func(username, password string) bool {
			return true
		},
    },
},
```

You need to implement `Validator` with your own logic

- Add to `configs/middlewares.yaml`

```yaml
middlewares:
    - basic-auth
```

