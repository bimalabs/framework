# Unit Test

To testing module that generated by Bima Cli is easy, for example we want to test `GetPaginated()` method in `todos` module from [ReadMe](../README.md#create-new-module)

## Create Test File

Create `module_test.go` in folder `todos`

## Create `Module` Dependencies

Here `Module` anatomy 

```go
type Module struct {
    *bima.Module
	Model     *Todo
    grpcs.UnimplementedTodosServer
}
```

For `Model` and `grpcs.UnimplementedTodosServer` are easy to pass, for `bima.Module` we need to mock `handlers.Handler` struct like below


```go
import (
    "testing"

    mocks "github.com/bimalabs/framework/v4/mocks/handlers"
)

func Test_Todos_GetPaginated(t *testing.T) {
	handler := mocks.NewHandler(t)

	handler.AssertExpectations(t)
}
```

## Writing Test Scenario

You need to mock `Paginate()` call because in `GetPaginated()` that method is called. Here the full code for testing `GetPaginated()`

```go
package todos

import (
	grpcs "app/protos/builds"
	"context"
	"testing"
	"time"

	bima "github.com/bimalabs/framework/v4"
	mocks "github.com/bimalabs/framework/v4/mocks/handlers"
	"github.com/bimalabs/framework/v4/paginations"
	"github.com/bimalabs/framework/v4/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_Todos_GetPaginated(t *testing.T) {
	paginator := paginations.Pagination{}
	handler := mocks.NewHandler(t)
	handler.On("Paginate", mock.Anything, mock.Anything).Return(paginations.Metadata{
		Page:     paginator.Page,
		Limit:    paginator.Limit,
		Previous: -1,
		Next:     -1,
		Total:    10,
	})

	module := Module{
		Module: &bima.Module{
			Debug:     true,
			Handler:   handler,
			Cache:     utils.NewCache(1 * time.Second),
			Paginator: &paginator,
		},
		Model:                    &Todo{},
		UnimplementedTodosServer: grpcs.UnimplementedTodosServer{},
	}

	_, err := module.GetPaginated(context.TODO(), &grpcs.Pagination{})

	assert.Nil(t, err)

	handler.AssertExpectations(t)
}

```
