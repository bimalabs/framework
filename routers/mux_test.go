package routers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bimalabs/framework/v4/loggers"
	"github.com/bimalabs/framework/v4/middlewares"
	middlewareMocks "github.com/bimalabs/framework/v4/mocks/middlewares"
	routeMocks "github.com/bimalabs/framework/v4/mocks/routes"
	"github.com/bimalabs/framework/v4/routes"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Test_Mux_Router(t *testing.T) {
	loggers.Default("test")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	endpoint := "0.0.0.0:111"
	conn, _ := grpc.DialContext(ctx, endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))

	server := runtime.NewServeMux()

	route := routeMocks.NewRoute(t)
	route.On("Path").Return("/without-middleware").Once()
	route.On("Method").Return(http.MethodGet).Once()
	route.On("SetClient", mock.Anything).Once()
	route.On("Middlewares").Return(nil).Once()
	route.On("Handle", mock.Anything, mock.Anything, mock.Anything).Once()

	router := MuxRouter{ApiPrefix: "/api"}
	router.Register([]routes.Route{route})

	assert.Equal(t, -255, router.Priority())
	assert.Equal(t, 1, len(router.routes))

	router.Handle(context.TODO(), server, conn)

	req := httptest.NewRequest("GET", "http://bima.framework/api/without-middleware", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	route.AssertExpectations(t)

	middleware := middlewareMocks.NewMiddleware(t)
	middleware.On("Attach", mock.Anything, mock.Anything).Return(false).Once()

	route = routeMocks.NewRoute(t)
	route.On("Path").Return("/middleware").Once()
	route.On("Method").Return(http.MethodGet).Once()
	route.On("SetClient", mock.Anything).Once()
	route.On("Middlewares").Return([]middlewares.Middleware{middleware}).Once()
	route.On("Handle", mock.Anything, mock.Anything, mock.Anything).Once()

	router = MuxRouter{ApiPrefix: "/api"}
	router.Register([]routes.Route{route})

	assert.Equal(t, -255, router.Priority())
	assert.Equal(t, 1, len(router.routes))

	router.Handle(context.TODO(), server, conn)

	req = httptest.NewRequest("GET", "http://bima.framework/api/middleware", nil)
	w = httptest.NewRecorder()

	server.ServeHTTP(w, req)

	route.AssertExpectations(t)

	middleware = middlewareMocks.NewMiddleware(t)
	middleware.On("Attach", mock.Anything, mock.Anything).Return(true).Once()

	route = routeMocks.NewRoute(t)
	route.On("Path").Return("/middleware-stop").Once()
	route.On("Method").Return(http.MethodGet).Once()
	route.On("SetClient", mock.Anything).Once()
	route.On("Middlewares").Return([]middlewares.Middleware{middleware}).Once()

	router = MuxRouter{ApiPrefix: "/api"}
	router.Register([]routes.Route{route})

	assert.Equal(t, -255, router.Priority())
	assert.Equal(t, 1, len(router.routes))

	router.Handle(context.TODO(), server, conn)

	req = httptest.NewRequest("GET", "http://bima.framework/api/middleware-stop", nil)
	w = httptest.NewRecorder()

	server.ServeHTTP(w, req)

	route.AssertExpectations(t)

	middleware = middlewareMocks.NewMiddleware(t)
	middleware.On("Attach", mock.Anything, mock.Anything).Return(true).Once()

	route = routeMocks.NewRoute(t)
	route.On("Path").Return("/middleware-stop").Once()
	route.On("Method").Return(http.MethodGet).Once()
	route.On("SetClient", mock.Anything).Once()
	route.On("Middlewares").Return([]middlewares.Middleware{middleware}).Once()

	router = MuxRouter{
		Debug:     true,
		ApiPrefix: "/api",
	}
	router.Register([]routes.Route{route})

	assert.Equal(t, -255, router.Priority())
	assert.Equal(t, 1, len(router.routes))

	router.Handle(context.TODO(), server, conn)

	req = httptest.NewRequest("GET", "http://bima.framework/api/middleware-stop", nil)
	w = httptest.NewRecorder()

	server.ServeHTTP(w, req)

	route.AssertExpectations(t)

	middleware = middlewareMocks.NewMiddleware(t)
	middleware.On("Attach", mock.Anything, mock.Anything).Return(false).Once()

	route = routeMocks.NewRoute(t)
	route.On("Path").Return("/middleware-stop").Once()
	route.On("Method").Return(http.MethodGet).Once()
	route.On("SetClient", mock.Anything).Once()
	route.On("Middlewares").Return([]middlewares.Middleware{middleware}).Once()
	route.On("Handle", mock.Anything, mock.Anything, mock.Anything).Once()

	router = MuxRouter{
		Debug:     true,
		ApiPrefix: "/api",
	}
	router.Register([]routes.Route{route})

	assert.Equal(t, -255, router.Priority())
	assert.Equal(t, 1, len(router.routes))

	router.Handle(context.TODO(), server, conn)

	req = httptest.NewRequest("GET", "http://bima.framework/api/middleware-stop", nil)
	w = httptest.NewRecorder()

	server.ServeHTTP(w, req)

	route.AssertExpectations(t)
}
