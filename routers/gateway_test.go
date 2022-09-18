package routers

import (
	"context"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/bimalabs/framework/v4/configs"
	"github.com/bimalabs/framework/v4/loggers"
	mocks "github.com/bimalabs/framework/v4/mocks/configs"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Test_Gateway_Router_Success(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	endpoint := "0.0.0.0:111"
	conn, _ := grpc.DialContext(ctx, endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))

	server := runtime.NewServeMux()

	grpc := mocks.NewServer(t)
	grpc.On("Handle", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

	router := GRpcGateway{}
	router.Register([]configs.Server{grpc})

	assert.Equal(t, 255, router.Priority())
	assert.Equal(t, 1, len(router.servers))

	router.Handle(context.TODO(), server, conn)

	req := httptest.NewRequest("GET", "http://bima.framework/handle", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	grpc.AssertExpectations(t)
}

func Test_Gateway_Router_Error(t *testing.T) {
	loggers.Default("test")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	endpoint := "0.0.0.0:111"
	conn, _ := grpc.DialContext(ctx, endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))

	server := runtime.NewServeMux()

	grpc := mocks.NewServer(t)
	grpc.On("Handle", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("gateway_error")).Once()

	router := GRpcGateway{}
	router.Register([]configs.Server{grpc})

	assert.Equal(t, 255, router.Priority())
	assert.Equal(t, 1, len(router.servers))

	router.Handle(context.TODO(), server, conn)

	req := httptest.NewRequest("GET", "http://bima.framework/handle", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	grpc.AssertExpectations(t)
}
