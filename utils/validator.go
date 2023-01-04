package utils

import (
	"context"
	"errors"
	"strings"

	"github.com/bimalabs/framework/v4/events"
	"github.com/bimalabs/framework/v4/loggers"
	engine "github.com/go-playground/validator/v10"
	"github.com/iancoleman/strcase"
)

type (
	validator struct {
		debug      bool
		dispatcher *events.Dispatcher
	}
)

func Validator(debug bool, dispatcher *events.Dispatcher) *validator {
	return &validator{debug: debug, dispatcher: dispatcher}
}

func (v *validator) Validate(object interface{}) (string, error) {
	ctx := context.WithValue(context.Background(), loggers.ScopeKey, "validator")

	var message strings.Builder
	if v.debug {
		message.WriteString("dispatching ")
		message.WriteString(events.BeforeCreateEvent.String())

		loggers.Logger.Debug(ctx, message.String())
	}

	event := events.Validation{Data: object}
	_ = v.dispatcher.Dispatch(events.BeforeValidation.String(), &event)

	if event.IsError {
		return event.Message, errors.New(event.Message)
	}

	err := engine.New().Struct(object)
	if err == nil {
		return "", nil
	}

	message.Reset()
	for _, ve := range err.(engine.ValidationErrors) {
		message.WriteString(strcase.ToDelimited(ve.Field(), '_'))
		message.WriteString(" is ")
		message.WriteString(ve.Tag())

		break
	}

	return message.String(), err
}
