package routers

import (
	"net/http"
	"testing"

	mocks "github.com/bimalabs/framework/v4/mocks/routers"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Test_Router(t *testing.T) {
	ctx := t.Context()

	router1 := mocks.NewRouter(t)
	router1.On("Handle", ctx, mock.Anything, mock.Anything).Once()
	router1.On("Priority").Return(1).Once()

	router2 := mocks.NewRouter(t)
	router2.On("Handle", ctx, mock.Anything, mock.Anything).Once()
	router2.On("Priority").Return(2).Once()

	factory := Factory{
		Routers: []Router{
			router1,
			router2,
		},
	}

	endpoint := "0.0.0.0:111"
	conn, err := grpc.NewClient(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() { _ = conn.Close() })

	factory.Sort()
	factory.Handle(ctx, http.NewServeMux(), conn)

	router1.AssertExpectations(t)
	router2.AssertExpectations(t)
}
