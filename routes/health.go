package routes

import (
	"net/http"

	"github.com/goccy/go-json"

	bima "github.com/bimalabs/framework/v4"
	"github.com/bimalabs/framework/v4/middlewares"
	"google.golang.org/grpc"
)

const HelthPath = "/health"

type Health struct {
}

func (h *Health) Path() string {
	return HelthPath
}

func (h *Health) Method() string {
	return http.MethodGet
}

func (h *Health) Middlewares() []middlewares.Middleware {
	return nil
}

func (h *Health) SetClient(client *grpc.ClientConn) {
}

func (h *Health) Handle(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	w.Header().Set("Content-Type", "application/json")

	_ = json.NewEncoder(w).Encode(map[string]string{
		"version": bima.Version,
		"name":    "Bima Framework",
		"author":  "Muhamad Surya Iksanudin<surya.iksanudin@gmail.com>",
		"link":    "https://github.com/bimalabs/framework",
	})
}
