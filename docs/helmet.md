# Harding security using Helmet 

- Add package using `go get github.com/bimalabs/middlewares`

- Add Helmet middleware to `dics/container.go`

```go
{
    Name:  "bima:middleware:helmet",
    Scope: bima.Application,
    Build: (*helmet.Helmet)(nil),
},
```

- Add to `configs/middlewares.yaml`

```yaml
middlewares:
    - helmet
```
