# Enable CORS 

- Add package using `go get github.com/bimalabs/middlewares`

- Add CORS middleware to `dics/container.go`

```go
{
    Name:  "bima:middleware:cors",
    Scope: bima.Application,
    Build: (*cors.Cors)(nil),
    Params: dingo.Params{
        "Options": cors.Options{},
    },
},
```

You can refer to [cors](github.com/rs/cors) for more options

- Add to `configs/middlewares.yaml`

```yaml
middlewares:
    - cors
```
