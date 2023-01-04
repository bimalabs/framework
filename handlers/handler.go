package handlers

import (
	"context"
	"strings"

	"github.com/bimalabs/framework/v4/events"
	"github.com/bimalabs/framework/v4/loggers"
	"github.com/bimalabs/framework/v4/paginations"
	"github.com/bimalabs/framework/v4/repositories"
)

type (
	handler struct {
		debug      bool
		dispatcher *events.Dispatcher
		repository repositories.Repository
		adapter    paginations.Adapter
	}

	Handler interface {
		Paginate(paginator *paginations.Pagination, result interface{}) paginations.Metadata
		Create(v interface{}) error
		Update(v interface{}, id string) error
		Bind(v interface{}, id string) error
		FindBy(v interface{}, filters ...repositories.Filter) error
		All(v interface{}) error
		Delete(v interface{}, id string) error
		Repository() repositories.Repository
	}
)

func New(debug bool, dispatcher *events.Dispatcher, repository repositories.Repository, adapter paginations.Adapter) Handler {
	return &handler{
		debug:      debug,
		dispatcher: dispatcher,
		repository: repository,
		adapter:    adapter,
	}
}

func (h *handler) Paginate(paginator *paginations.Pagination, result interface{}) paginations.Metadata {
	ctx := context.WithValue(context.Background(), loggers.ScopeKey, "handler")

	adapter := h.adapter.CreateAdapter(ctx, *paginator)
	if adapter == nil {
		loggers.Logger.Error(ctx, "error when creating adapter")

		return paginations.Metadata{}
	}

	var total64 int64
	_ = paginator.Paginate(adapter, result, &total64)

	var total = int(total64)
	next := paginator.Page + 1
	if paginator.Page*paginator.Limit >= int(total) {
		next = -1
	}

	return paginations.Metadata{
		Page:     paginator.Page,
		Previous: paginator.Page - 1,
		Next:     next,
		Limit:    paginator.Limit,
		Total:    total,
	}
}

func (h *handler) Create(v interface{}) error {
	return h.repository.Transaction(func(r repositories.Repository) error {
		var log strings.Builder
		ctx := context.WithValue(context.Background(), loggers.ScopeKey, "handler")
		if h.debug {
			log.WriteString("dispatching ")
			log.WriteString(events.BeforeCreateEvent.String())

			loggers.Logger.Debug(ctx, log.String())
		}

		_ = h.dispatcher.Dispatch(events.BeforeCreateEvent.String(), &events.Model{
			Data:       v,
			Repository: r,
		})

		if err := r.Create(v); err != nil {
			loggers.Logger.Error(ctx, "error when creating resource(s), Rolling back")

			return err
		}

		if h.debug {
			log.Reset()
			log.WriteString("dispatching ")
			log.WriteString(events.AfterCreateEvent.String())

			loggers.Logger.Debug(ctx, log.String())
		}

		_ = h.dispatcher.Dispatch(events.AfterCreateEvent.String(), &events.Model{
			Data:       v,
			Repository: r,
		})

		return nil
	})
}

func (h *handler) Update(v interface{}, id string) error {
	return h.repository.Transaction(func(r repositories.Repository) error {
		var log strings.Builder
		ctx := context.WithValue(context.Background(), loggers.ScopeKey, "handler")
		if h.debug {
			log.WriteString("dispatching ")
			log.WriteString(events.BeforeUpdateEvent.String())

			loggers.Logger.Debug(ctx, log.String())
		}

		_ = h.dispatcher.Dispatch(events.BeforeUpdateEvent.String(), &events.Model{
			Id:         id,
			Data:       v,
			Repository: r,
		})

		if err := r.Update(v); err != nil {
			loggers.Logger.Error(ctx, "error when updating resource(s), Rolling back")

			return err
		}
		if h.debug {
			log.Reset()
			log.WriteString("dispatching ")
			log.WriteString(events.AfterUpdateEvent.String())

			loggers.Logger.Debug(ctx, log.String())
		}

		_ = h.dispatcher.Dispatch(events.AfterUpdateEvent.String(), &events.Model{
			Id:         id,
			Data:       v,
			Repository: r,
		})

		return nil
	})
}

func (h *handler) Bind(v interface{}, id string) error {
	return h.repository.Bind(v, id)
}

func (h *handler) All(v interface{}) error {
	return h.repository.All(v)
}

func (h *handler) FindBy(v interface{}, filters ...repositories.Filter) error {
	return h.repository.FindBy(v, filters...)
}

func (h *handler) Delete(v interface{}, id string) error {
	return h.repository.Transaction(func(r repositories.Repository) error {
		var log strings.Builder
		ctx := context.WithValue(context.Background(), loggers.ScopeKey, "handler")
		if h.debug {
			log.WriteString("dispatching ")
			log.WriteString(events.BeforeDeleteEvent.String())

			loggers.Logger.Debug(ctx, log.String())
		}

		_ = h.dispatcher.Dispatch(events.BeforeDeleteEvent.String(), &events.Model{
			Id:         id,
			Data:       v,
			Repository: r,
		})

		if err := r.Delete(v, id); err != nil {
			loggers.Logger.Error(ctx, "error when deleting resource(s), Rolling back")

			return err
		}

		if h.debug {
			log.Reset()
			log.WriteString("dispatching ")
			log.WriteString(events.AfterDeleteEvent.String())

			loggers.Logger.Debug(ctx, log.String())
		}

		_ = h.dispatcher.Dispatch(events.AfterDeleteEvent.String(), &events.Model{
			Id:         id,
			Data:       v,
			Repository: r,
		})

		return nil
	})
}

func (h *handler) Repository() repositories.Repository {
	return h.repository
}
