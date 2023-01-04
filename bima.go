package bima

import (
	"github.com/bimalabs/framework/v4/handlers"
	"github.com/bimalabs/framework/v4/messengers"
	"github.com/bimalabs/framework/v4/models"
	"github.com/bimalabs/framework/v4/paginations"
	"github.com/bimalabs/framework/v4/utils"
	"github.com/olivere/elastic/v7"
	"gorm.io/gorm"
)

const (
	Version = "v4.3.1"

	HighestPriority = 255
	LowestPriority  = -255

	Application = "application"
	Generator   = "generator"
)

type (
	Module interface {
		Debug() bool
		Handler() handlers.Handler
		Cache() *utils.Cache
		Paginator() *paginations.Pagination
		Validate(v interface{}) (string, error)
	}

	module struct {
		debug     bool
		handler   handlers.Handler
		cache     *utils.Cache
		paginator *paginations.Pagination
		validator utils.Validator
	}

	GormModel struct {
		models.GormBase
	}

	Server struct {
		Debug bool
	}
)

func NewModule(debug bool, handler handlers.Handler, cache *utils.Cache, validator utils.Validator, paginator *paginations.Pagination) Module {
	return &module{
		debug:     debug,
		handler:   handler,
		cache:     cache,
		paginator: paginator,
		validator: validator,
	}
}

func (m *module) Debug() bool {
	return m.debug
}

func (m *module) Handler() handlers.Handler {
	return m.handler
}

func (m *module) Cache() *utils.Cache {
	return m.cache
}

func (m *module) Paginator() *paginations.Pagination {
	return m.paginator
}

func (m *module) Validate(v interface{}) (string, error) {
	return m.validator.Validate(v)
}

func NewModel() GormModel {
	return GormModel{
		GormBase: models.GormBase{},
	}
}

func (s *Server) Consume(messenger *messengers.Messenger) {
}

func (s *Server) Sync(client *elastic.Client) {
}

func (s *Server) Migrate(db *gorm.DB) {
}
