package domyApi

import (
	"context"
	"log"
	"net/http"
	"strings"

	config "github.com/domyid/domyapi/config"
	at "github.com/domyid/domyapi/helper/at"
	api "github.com/domyid/domyapi/helper/atapi"
	atdb "github.com/domyid/domyapi/helper/atdb"
	model "github.com/domyid/domyapi/model"
	"go.mongodb.org/mongo-driver/bson"
)

func GetMahasiswa(respw http.ResponseWriter, req *http.Request) {
	urlTarget := "https://siakad.ulbi.ac.id/siakad/data_mahasiswa"

	// Ambil login dari header
	login := at.GetLoginFromHeader(req)
	if login == "" {
		http.Error(respw, "No valid login found", http.StatusForbidden)
		return
	}

	// Buat payload berisi informasi login
	payload := map[string]string{
		"SIAKAD_CLOUD_ACCESS": login,
	}

	doc, err := api.GetData(urlTarget, payload, nil)
	if err != nil {
		at.WriteJSON(respw, http.StatusInternalServerError, err.Error())
		return
	}

	// Ekstrak informasi mahasiswa dan hapus spasi berlebih
	nim := strings.TrimSpace(doc.Find("#block-nim .input-nim").Text())
	nama := strings.TrimSpace(doc.Find("#block-nama .input-nama").Text())
	programStudi := strings.TrimSpace(doc.Find("#block-idunit .input-idunit").Text())
	noHp := strings.TrimSpace(doc.Find("#block-hp .input-hp").Text())

	// Buat instance Mahasiswa
	mahasiswa := model.Mahasiswa{
		NIM:          nim,
		Nama:         nama,
		ProgramStudi: programStudi,
		NomorHp:      noHp,
	}

	// Kembalikan instance Mahasiswa sebagai respon JSON
	at.WriteJSON(respw, http.StatusOK, mahasiswa)
}

func PostBimbinganMahasiswa(w http.ResponseWriter, r *http.Request) {
	urlTarget := "https://siakad.ulbi.ac.id/siakad/data_bimbingan/add/964"

	// Ambil login dari header
	login := at.GetLoginFromHeader(r)
	if login == "" {
		http.Error(w, "No valid login found", http.StatusForbidden)
		return
	}

	// Buat payload berisi informasi login
	payload := map[string]string{
		"SIAKAD_CLOUD_ACCESS": login,
	}

	formData := map[string]string{
		"bimbinganke":    r.FormValue("bimbinganke"),
		"nip":            r.FormValue("nip"),
		"tglbimbingan":   r.FormValue("tglbimbingan"),
		"topikbimbingan": r.FormValue("topikbimbingan"),
		"bahasan":        r.FormValue("bahasan"),
		"link[]":         r.FormValue("link[]"),
		"key":            r.FormValue("key"),
		"act":            r.FormValue("act"),
	}

	fileFieldName := "lampiran[]"
	filePath := "" // Kosongkan path file

	resp, err := api.PostData(urlTarget, payload, formData, fileFieldName, filePath)
	if err != nil {
		log.Printf("Error in PostBimbinganMahasiswa: %v", err)
		at.WriteJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusSeeOther && resp.StatusCode != http.StatusOK {
		at.WriteJSON(w, resp.StatusCode, "unexpected status code")
		return
	}

	at.WriteJSON(w, resp.StatusCode, "success")
}

func GetDosen(respw http.ResponseWriter, req *http.Request) {
	urltarget := req.URL.Query().Get("url")
	if urltarget == "" {
		http.Error(respw, "url query parameter is required", http.StatusBadRequest)
		return

	}

	doc, err := api.GetData(urltarget, nil, nil)
	if err != nil {
		at.WriteJSON(respw, http.StatusInternalServerError, err.Error())
		return
	}

	// Ekstrak informasi dosen dengan trim spasi
	nip := strings.TrimSpace(doc.Find("#block-nip .input-nip").Text())
	nidn := strings.TrimSpace(doc.Find("#block-nidn .input-nidn").Text())
	nama := strings.TrimSpace(doc.Find("#block-nama .input-nama").Text())
	noHp := strings.TrimSpace(doc.Find("#block-nohp .input-nohp").Text())

	// Buat instance Dosen
	dosen := model.Dosen{
		NIP:  nip,
		NIDN: nidn,
		Nama: nama,
		NoHp: noHp,
	}

	// Cek apakah data dosen sudah ada di database
	filter := bson.M{"nip": dosen.NIP}
	var existingDosen model.Dosen
	if err := config.Mongoconn.Collection("dosen").FindOne(context.TODO(), filter).Decode(&existingDosen); err == nil {
		// Data sudah ada, tidak perlu menambah data baru
		at.WriteJSON(respw, http.StatusOK, existingDosen)
		return
	}

	// Simpan ke MongoDB jika data belum ada
	if _, err := atdb.InsertOneDoc(config.Mongoconn, "dosen", dosen); err != nil {
		at.WriteJSON(respw, http.StatusInternalServerError, err.Error())
		return
	}

	// Konversi ke JSON dan kirimkan sebagai respon
	at.WriteJSON(respw, http.StatusOK, dosen)
}

func NotFound(respw http.ResponseWriter, req *http.Request) {
	var resp model.Response
	resp.Response = "Not Found"
	at.WriteJSON(respw, http.StatusNotFound, resp)
}
