# Implement Distributed Tracing (RequestID)

- Add package using `go get github.com/bimalabs/middlewares`

- Add requestid middleware to `dics/container.go`

```go
{
    Name:  "bima:middleware:requestid",
    Scope: bima.Application,
    Build: (*requestid.RequestID)(nil),
    Params: dingo.Params{
        "RequestIDHeader": "X-Request-Id",
    },
},
```

- Add to `configs/middlewares.yaml`

```yaml
middlewares:
    - requestid
```

RequestID added to response header

