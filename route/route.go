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

	var method, path string = r.Method, r.URL.Path

	switch {
	case method == "POST" && path == "/login":
		controller.LoginSiakad(w, r)
	case method == "POST" && path == "/refresh-token":
		controller.RefreshTokens(w, r)
	case method == "GET" && path == "/data/mahasiswa":
		controller.GetMahasiswa(w, r)
	case method == "GET" && path == "/data/dosen":
		controller.GetDosen(w, r)
	case method == "POST" && path == "/data/bimbingan/mahasiswa":
		controller.PostBimbinganMahasiswa(w, r)
	default:
		controller.NotFound(w, r)
	}
}
