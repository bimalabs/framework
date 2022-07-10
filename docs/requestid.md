# Implement Distributed Tracing (RequestID)

- Add package using `go get github.com/bimalabs/middlewares`

- Add requestid middleware to `dics/container.go`

```go
{
    Name:  "bima:middleware:requestid",
    Scope: bima.Application,
    Build: func(header string)(middlewares.Middleware, nil) {
        return requestid.New(header), nil
    },
    Params: dingo.Params{
        "0": "X-Request-Id",
    },
},
```

- Add to `configs/middlewares.yaml`

```yaml
middlewares:
    - requestid
```

RequestID added to response header

