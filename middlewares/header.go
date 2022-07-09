package middlewares

import (
	"net/http"
)

type Header struct {
}

func (h *Header) Attach(_ *http.Request, response http.ResponseWriter) bool {
	response.Header().Add("X-Server-Id", "Bima")

	return false
}

func (h *Header) Priority() int {
	return -257
}
