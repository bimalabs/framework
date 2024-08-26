package listeners

import (
	"strings"

	"github.com/bimalabs/framework/v4/events"
)

type GormFilter struct {
}

func (p *GormFilter) Handle(event interface{}) interface{} {
	e, ok := event.(*events.GormPagination)
	if !ok {
		return event
	}

	var likeClausal strings.Builder
	var likeValue strings.Builder

	for k, v := range e.Filters {
		likeClausal.Reset()
		likeClausal.WriteString(k)
		likeClausal.WriteString(" LIKE ?")

		likeValue.Reset()
		likeValue.WriteString("%")
		likeValue.WriteString(v)
		likeValue.WriteString("%")

		e.Query.Where(likeClausal.String(), likeValue.String())
	}

	return e
}

func (p *GormFilter) Listen() string {
	return events.PaginationEvent.String()
}

func (p *GormFilter) Priority() int {
	return 255
}
