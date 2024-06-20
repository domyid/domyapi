package domyApi

import (
	"net/http"

	config "github.com/domyid/domyapi/config"
	controller "github.com/domyid/domyapi/controller"
)

func URL(w http.ResponseWriter, r *http.Request) {
	if config.SetAccessControlHeaders(w, r) {
		return // If it's a preflight request, return early.
	}
	config.SetEnv()

	// var method, path string = r.Method, r.URL.Path
	switch {
	// case method == "GET" && path == "/":
	// 	controller.getMahasiswa(w, r)

	default:
		controller.NotFound(w, r)
	}
}
