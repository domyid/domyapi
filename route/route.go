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
	case method == "GET" && path == "/data/bimbingan/mahasiswa":
		controller.GetListBimbinganMahasiswa(w, r)
	case method == "POST" && path == "/data/bimbingan/mahasiswa":
		controller.PostBimbinganMahasiswa(w, r)
	case method == "GET" && path == "/data/dosen":
		controller.GetDosen(w, r)
	case method == "POST" && path == "/jadwalmengajar":
		controller.GetJadwalMengajar(w, r)
	case method == "POST" && path == "/riwayatmengajar":
		controller.GetRiwayatPerkuliahan(w, r)
	case method == "POST" && path == "/absensi":
		controller.GetAbsensiKelas(w, r)
	case method == "POST" && path == "/nilai":
		controller.GetNilaiMahasiswa(w, r)
	case method == "POST" && path == "/BAP":
		controller.GetBAP(w, r)
	case method == "POST" && path == "/ApproveBAP":
		controller.ApproveBAP(w, r)
	case method == "POST" && path == "/StatusApproval":
		controller.CekStatusApproval(w, r)
	case method == "GET" && path == "/data/list/ta":
		controller.GetListTugasAkhirMahasiswa(w, r)
	case method == "POST" && path == "/data/list/bimbingan":
		controller.GetListBimbinganMahasiswa(w, r)
	case method == "POST" && path == "/approve/bimbingan":
		controller.ApproveBimbingan(w, r)
	case method == "POST" && path == "/BKD":
		controller.GenerateBAPBimbingan(w, r)
	default:
		controller.NotFound(w, r)
	}
}
