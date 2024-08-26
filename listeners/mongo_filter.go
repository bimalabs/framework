package listeners

import (
	"strings"

	"github.com/bimalabs/framework/v4/events"
	"github.com/kamva/mgm/v3/operator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MongoDbFilter struct {
}

func (p *MongoDbFilter) Handle(event interface{}) interface{} {
	e, ok := event.(*events.MongodbPagination)
	if !ok {
		return event
	}

	bFilters := bson.M{}
	for k, v := range e.Filters {
		bFilters[strings.ToLower(k)] = bson.M{
			operator.Regex: primitive.Regex{
				Pattern: v,
				Options: "im",
			},
		}
	}

	e.MongoDbFilter = bFilters

	return e
}

func (p *MongoDbFilter) Listen() string {
	return events.PaginationEvent.String()
}

func (p *MongoDbFilter) Priority() int {
	return 255
}
