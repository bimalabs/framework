package middlewares

import (
	"net/http"

	bima "github.com/bimalabs/framework/v4"
)

type Header struct {
}

func (h *Header) Attach(_ *http.Request, response http.ResponseWriter) bool {
	response.Header().Add("X-Bima-Version", bima.Version)

	return false
}

func (h *Header) Priority() int {
	return -257
}
