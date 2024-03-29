package routes

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Test_Api_Doc(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	endpoint := "0.0.0.0:111"
	conn, _ := grpc.DialContext(ctx, endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))

	req := httptest.NewRequest("GET", "http://bima.framework/api/docs/", nil)
	w := httptest.NewRecorder()

	route := ApiDoc{
		Debug: true,
	}
	route.SetClient(conn)

	route.Handle(w, req, map[string]string{})

	resp := w.Result()

	assert.Equal(t, http.MethodGet, route.Method())
	assert.Equal(t, fmt.Sprintf("%s/{path}", ApiDocPath), route.Path())
	assert.Nil(t, route.Middlewares())

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	route = ApiDoc{
		Debug: false,
	}
	route.SetClient(conn)

	route.Handle(w, req, map[string]string{})

	resp = w.Result()

	assert.Equal(t, http.MethodGet, route.Method())
	assert.Equal(t, fmt.Sprintf("%s/{path}", ApiDocPath), route.Path())
	assert.Nil(t, route.Middlewares())

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}
