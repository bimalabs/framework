package routes

import (
	"fmt"
	"net/http"

	"github.com/KejawenLab/bima/v2/configs"
	"google.golang.org/grpc"
)

const API_DOC_REDIRECT_PATH = "/api/docs"

type ApiDocRedirect struct {
}

func (a *ApiDocRedirect) Path() string {
	return API_DOC_REDIRECT_PATH
}

func (a *ApiDocRedirect) Method() string {
	return http.MethodGet
}

func (a *ApiDocRedirect) SetClient(client *grpc.ClientConn) {
}

func (a *ApiDocRedirect) Middlewares() []configs.Middleware {
	return nil
}

func (a *ApiDocRedirect) Handle(w http.ResponseWriter, r *http.Request, params map[string]string) {
	http.Redirect(w, r, fmt.Sprintf("%s/", r.URL.RequestURI()), http.StatusPermanentRedirect)
}