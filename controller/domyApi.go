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
	urltarget := "https://siakad.ulbi.ac.id/siakad/data_mahasiswa"
	log.Println("Request URL:", urltarget)

	// Get the cookies from the request
	cookies := make(map[string]string)
	for _, cookie := range req.Cookies() {
		log.Printf("Received cookie: %s = %s", cookie.Name, cookie.Value)
		cookies[cookie.Name] = cookie.Value
	}

	// Fetch the data from the URL
	doc, err := api.GetData(urltarget, cookies, nil)
	if err != nil {
		at.WriteJSON(respw, http.StatusInternalServerError, fmt.Sprintf("failed to fetch data: %v", err))
		return
	}

	// Extract information
	nim := strings.TrimSpace(doc.Find("#block-nim .input-nim").Text())
	nama := strings.TrimSpace(doc.Find("#block-nama .input-nama").Text())
	programStudi := strings.TrimSpace(doc.Find("#block-idunit .input-idunit").Text())
	noHp := strings.TrimSpace(doc.Find("#block-hp .input-hp").Text())

	if nim == "" || nama == "" || programStudi == "" || noHp == "" {
		at.WriteJSON(respw, http.StatusNotFound, "data not found")
		return
	}

	mahasiswa := model.Mahasiswa{
		NIM:          nim,
		Nama:         nama,
		ProgramStudi: programStudi,
		NomorHp:      noHp,
	}

	at.WriteJSON(respw, http.StatusOK, mahasiswa)
}

// PostMahasiswa handles the POST request to add mahasiswa data.
func PostBimbinganMahasiswa(w http.ResponseWriter, r *http.Request) {
	urlTarget := "https://siakad.ulbi.ac.id/siakad/data_bimbingan/add/964"

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

	resp, err := api.PostData(urlTarget, cookies, formData, fileFieldName, filePath)
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
