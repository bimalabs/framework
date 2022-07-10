# Enable CORS 

- Add package using `go get github.com/bimalabs/middlewares`

- Add CORS middleware to `dics/container.go`

```go
{
    Name:  "bima:middleware:cors",
    Scope: bima.Application,
    Build: func (options cors.Options) (middlewares.Middleware, error) {
        return cors.New(options), nil
    },
    Params: dingo.Params{
        "0": cors.Options{},
    },
},
```

You can refer to [cors](github.com/rs/cors) for more options

- Add to `configs/middlewares.yaml`

```yaml
middlewares:
    - cors
```
