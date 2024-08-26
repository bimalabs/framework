package paginations

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestAdapter struct {
	err error
}

func (p TestAdapter) Nums() (int64, error) {
	return 10, p.err
}

func (p TestAdapter) Slice(offset, length int, data interface{}) error {
	s := data.(*[]interface{})
	for n := offset + 1; n < offset+length+1; n++ {
		*s = append(*s, n)
	}

	return nil
}

func Test_Pagination_Handle_Request(t *testing.T) {
	pagination := Pagination{}

	request := Request{
		Page:  1,
		Limit: 17,
	}

	pagination.Handle(request)

	assert.Equal(t, pagination.Limit, int(request.Limit))
	assert.Equal(t, pagination.Page, int(request.Page))

	request = Request{
		Page:  0,
		Limit: 0,
	}

	pagination.Handle(request)

	assert.Equal(t, pagination.Limit, 17)
	assert.Equal(t, pagination.Page, 1)

	request = Request{
		Filters: Filter{"a": "b"},
	}

	assert.Equal(t, len(pagination.Filters), 1)
}

func Test_Pagination_Paginate(t *testing.T) {
	var total int64
	result := []interface{}{}
	pagination := Pagination{}
	pagination.Handle(Request{})
	_ = pagination.Paginate(TestAdapter{}, &result, &total)

	assert.Equal(t, int64(10), total)
	assert.Equal(t, len(result), 17)

	_ = pagination.Paginate(TestAdapter{err: errors.New("test")}, &result, &total)

	assert.Equal(t, int64(10), total)
	assert.Equal(t, len(result), 17)
}
