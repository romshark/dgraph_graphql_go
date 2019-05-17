package http

import (
	"net/http"
)

// servePlayground serves the playground app if it's enabled
func (t *Server) servePlayground(
	resp http.ResponseWriter,
	req *http.Request,
) {
	if !t.opts.Playground {
		http.Error(
			resp,
			http.StatusText(http.StatusNotFound),
			http.StatusNotFound,
		)
		return
	}

	if _, err := resp.Write(playgroundSrc); err != nil {
		http.Error(
			resp,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
		return
	}
}
