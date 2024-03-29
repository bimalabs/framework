package paginations

import (
	"context"

	"github.com/vcraescu/go-paginator/v2"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type (
	Adapter interface {
		CreateAdapter(ctx context.Context, paginator Pagination) paginator.Adapter
	}

	Pagination struct {
		Limit   int
		Page    int
		Filters []Filter
		Search  string
		Model   interface{}
		Table   string
	}

	Filter struct {
		Field string
		Value string
	}

	Metadata struct {
		Page     int
		Previous int
		Next     int
		Limit    int
		Total    int
	}

	Request struct {
		Page   int32
		Limit  int32
		Fields []string
		Values []string
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

	n := len(request.Fields)
	if n == 0 || n != len(request.Values) {
		return
	}

	p.Filters = make([]Filter, 0, len(request.Fields))
	for k, v := range request.Fields {
		if v == "" || request.Values[k] == "" {
			continue
		}

		p.Filters = append(p.Filters, Filter{Field: cases.Title(language.English).String(v), Value: request.Values[k]})
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
