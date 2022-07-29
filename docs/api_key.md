# Implement Api Key Auth

- Add package using `go get github.com/bimalabs/middlewares`

- Add basic auth middleware to `dics/container.go`

```go
{
    Name:  "bima:middleware:api-key",
    Scope: bima.Application,
    Build: func() (middlewares.Middleware, error) {
        return api_key.New(api_key.LocationQueries, "api", "my-api-key"), nil
    },
},
```

- Available locations `api_key.LocationQueries` and `api_key.LocationHeader`

- Add to `configs/middlewares.yaml`

```yaml
middlewares:
    - api-key
```

