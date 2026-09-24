package events

import (
	"github.com/bimalabs/framework/v4/paginations"
	"github.com/bimalabs/framework/v4/repositories"
	"github.com/kamva/mgm/v3"
	"github.com/olivere/elastic/v7"
	"go.mongodb.org/mongo-driver/bson"
	"gorm.io/gorm"
)

type (
	Event string

	Model struct {
		Data       any
		Repository repositories.Repository
		Id         string
	}

	Validation struct {
		Data    any
		Message string
		IsError bool
	}

	ElasticsearchPagination struct {
		Model   any
		Query   *elastic.BoolQuery
		Filters paginations.Filter
	}

	MongodbPagination struct {
		Model         any
		Query         *mgm.Collection
		Filters       paginations.Filter
		MongoDbFilter bson.M
	}

	GormPagination struct {
		Model   any
		Query   *gorm.DB
		Filters paginations.Filter
	}
)

const (
	PaginationEvent   = Event("pagination")
	BeforeValidation  = Event("before_validation")
	BeforeCreateEvent = Event("before_create")
	BeforeUpdateEvent = Event("before_update")
	BeforeDeleteEvent = Event("before_delete")
	AfterValidation   = Event("after_validation")
	AfterCreateEvent  = Event("after_create")
	AfterUpdateEvent  = Event("after_update")
	AfterDeleteEvent  = Event("after_delete")
)

func (e Event) String() string {
	return string(e)
}
