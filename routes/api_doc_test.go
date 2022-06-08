package routes

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func Test_Api_Doc(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	endpoint := "0.0.0.0:111"
	conn, _ := grpc.DialContext(ctx, endpoint, grpc.WithInsecure())

	req := httptest.NewRequest("GET", "http://example.com/api/docs/", nil)
	w := httptest.NewRecorder()

	route := ApiDoc{
		Debug: true,
	}
	route.SetClient(conn)

	route.Handle(w, req, map[string]string{})

	resp := w.Result()

	assert.Equal(t, http.MethodGet, route.Method())
	assert.Equal(t, fmt.Sprintf("%s/{path}", API_DOC_PATH), route.Path())
	assert.Nil(t, route.Middlewares())

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}
