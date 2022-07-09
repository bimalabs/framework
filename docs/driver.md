# Create Own Database Driver

- We use `sqlite` for this example `https://gorm.io/docs/connecting_to_the_database.html#SQLite` 

- Create `sqlite.go` and here the code

```go
package drivers

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Sqlite string

func (s Sqlite) Name() string {
    return string(s)
}

func (_ Sqlite) Connect(_ string, _ int, _ string, _ string, dbname string, _ bool) *gorm.DB {
    db, err := gorm.Open(sqlite.Open(dbname), &gorm.Config{})
    if err != nil {
        panic(err)
    }

	return db
}

```

- Add definition to `dics/container.go`

```go
{
    Name:  "bima:driver:sqlite",
    Scope: bima.Application,
    Build: func() (*drivers.Sqlite, error) {
        return Sqlite("sqlite")
    },
},
```

- Add to `configs/drivers.yaml`

```yaml
drivers:
    - sqlite
```
