package drivers

import (
	"testing"

	mocks "github.com/bimalabs/framework/v4/mocks/drivers"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func Test_Driver(t *testing.T) {
	factory := New(true)

	driver := mocks.NewDriver(t)
	driver.On("Name").Return("test").Once()
	driver.On("Connect", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&gorm.DB{}).Once()

	factory.Register([]Driver{driver})
	factory.Connect("test", "", 0, "", "", "")

	driver.AssertExpectations(t)

	factory.Connect("", "", 0, "", "", "")
}
