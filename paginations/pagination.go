package paginations

import (
	"context"

	"github.com/iancoleman/strcase"
	"github.com/vcraescu/go-paginator/v2"
)

type (
	Adapter interface {
		CreateAdapter(ctx context.Context, paginator Pagination) paginator.Adapter
	}

	Pagination struct {
		Limit   int
		Page    int
		Filters Filter
		Search  string
		Model   interface{}
		Table   string
	}

	Filter map[string]string

	Metadata struct {
		Page     int
		Previous int
		Next     int
		Limit    int
		Total    int
	}

	Request struct {
		Page    int32
		Limit   int32
		Filters Filter
	}
)

func (p *Pagination) Handle(request Request) {
	if request.Page == 0 {
		request.Page = 1
	}

	if request.Limit == 0 {
		request.Limit = 17
	}

	p.Limit = int(request.Limit)
	p.Page = int(request.Page)

	p.Filters = Filter{}
	for k, v := range request.Filters {
		p.Filters[strcase.ToSnake(k)] = v
	}
}

func (p *Pagination) Paginate(adapter paginator.Adapter, results interface{}, total *int64) error {
	pager := paginator.New(adapter, p.Limit)

	pager.SetPage(p.Page)
	p.Page, _ = pager.Page()

	err := pager.Results(results)
	if err != nil {
		return err
	}

	*total, err = pager.Nums()

	return err
}
