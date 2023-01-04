package utils

import (
	"testing"

	"github.com/bimalabs/framework/v4/events"
	"github.com/bimalabs/framework/v4/loggers"
	"github.com/stretchr/testify/assert"
)

type Data struct {
	ID   string `validate:"required"`
	Name string `validate:"required"`
}

func Test_Validator(t *testing.T) {
	loggers.Default("test")

	dispatcher := events.Dispatcher{}
	validator := Validator(true, &dispatcher)

	data1 := Data{
		ID: "test",
	}

	msg, err := validator.Validate(&data1)

	assert.NotNil(t, err)
	assert.NotEmpty(t, msg)

	data2 := Data{
		ID:   "test",
		Name: "test",
	}

	msg, err = validator.Validate(&data2)

	assert.Nil(t, err)
	assert.Empty(t, msg)
}
