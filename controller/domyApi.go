package domyApi

import (
	"context"
	"fmt"
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
	urltarget := req.URL.Query().Get("url")
	if urltarget == "" {
		http.Error(respw, "url query parameter is required", http.StatusBadRequest)
		return
	}

	cookies := make(map[string]string)
	for _, cookie := range req.Cookies() {
		cookies[cookie.Name] = cookie.Value
	}

	log.Printf("Fetching data from URL: %s with cookies: %v", urltarget, cookies)
	doc, err := api.FetchDataFromURL(urltarget, cookies, nil)
	if err != nil {
		http.Error(respw, fmt.Sprintf("failed to fetch data: %v", err), http.StatusInternalServerError)
		return
	}

	// Ekstrak informasi mahasiswa dengan trim spasi
	nim := strings.TrimSpace(doc.Find("#block-nim .input-nim").Text())
	nama := strings.TrimSpace(doc.Find("#block-nama .input-nama").Text())
	programStudi := strings.TrimSpace(doc.Find("#block-idunit .input-idunit").Text())
	noHp := strings.TrimSpace(doc.Find("#block-hp .input-hp").Text())

	log.Printf("Extracted data - NIM: %s, Nama: %s, ProgramStudi: %s, NoHp: %s", nim, nama, programStudi, noHp)

	// Buat instance Mahasiswa
	mahasiswa := model.Mahasiswa{
		NIM:          nim,
		Nama:         nama,
		ProgramStudi: programStudi,
		NomorHp:      noHp,
	}

	// Cek apakah data mahasiswa sudah ada di database
	filter := bson.M{"nim": mahasiswa.NIM}
	var existingMahasiswa model.Mahasiswa
	if err := config.Mongoconn.Collection("mahasiswa").FindOne(context.TODO(), filter).Decode(&existingMahasiswa); err == nil {
		// Data sudah ada, tidak perlu menambah data baru
		at.WriteJSON(respw, http.StatusOK, existingMahasiswa)
		return
	}

	// Simpan ke MongoDB jika data belum ada
	if _, err := atdb.InsertOneDoc(config.Mongoconn, "mahasiswa", mahasiswa); err != nil {
		http.Error(respw, fmt.Sprintf("failed to insert document into MongoDB: %v", err), http.StatusInternalServerError)
		return
	}

	// Konversi ke JSON dan kirimkan sebagai respon
	at.WriteJSON(respw, http.StatusOK, mahasiswa)
}

// PostMahasiswa handles the POST request to add mahasiswa data.
func PostBimbinganMahasiswa(w http.ResponseWriter, r *http.Request) {
	urlTarget := r.FormValue("urlTarget")
	if urlTarget == "" {
		at.WriteJSON(w, http.StatusBadRequest, "urlTarget parameter is required")
		return
	}

	cookies := make(map[string]string)
	for _, cookie := range r.Cookies() {
		cookies[cookie.Name] = cookie.Value
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

	resp, err := api.PostDataToURL(urlTarget, cookies, formData, fileFieldName, filePath)
	if err != nil {
		log.Printf("Error in PostBimbinganMahasiswa: %v", err)
		at.WriteJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer resp.Body.Close()

	// log.Printf("Response Status: %v", resp.Status)
	// log.Printf("Response Headers: %v", resp.Header)

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

	doc, err := api.FetchDataFromURL(urltarget, nil, nil)
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
