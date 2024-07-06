package domyApi

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	config "github.com/domyid/domyapi/config"
	at "github.com/domyid/domyapi/helper/at"
	api "github.com/domyid/domyapi/helper/atapi"
	atdb "github.com/domyid/domyapi/helper/atdb"
	model "github.com/domyid/domyapi/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetMahasiswa(respw http.ResponseWriter, req *http.Request) {
	urlTarget := "https://siakad.ulbi.ac.id/siakad/data_mahasiswa"

	// Mengambil user_id dari header
	userID := req.Header.Get("user_id")
	if userID == "" {
		http.Error(respw, "No valid user ID found", http.StatusForbidden)
		return
	}

	// Mengambil token dari database berdasarkan user_id
	tokenData, err := atdb.GetOneDoc[model.TokenData](config.Mongoconn, "tokens", primitive.M{"user_id": userID})
	if err != nil {
		fmt.Println("Error Fetching Token:", err)
		at.WriteJSON(respw, http.StatusNotFound, "Token not found for user")
		return
	}

	// Buat payload berisi informasi token
	payload := map[string]string{
		"SIAKAD_CLOUD_ACCESS": tokenData.Token,
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
	// Mengambil user_id dari header
	userID := r.Header.Get("user_id")
	if userID == "" {
		http.Error(w, "No valid user ID found", http.StatusForbidden)
		return
	}

	// Memanggil fungsi helper untuk mendapatkan list tugas akhir
	listTA, err := api.FetchListTugasAkhirMahasiswa(userID)
	if err != nil || len(listTA) == 0 {
		at.WriteJSON(w, http.StatusNotFound, "Failed to fetch Tugas Akhir or no data found")
		return
	}

	// Ambil data_id dari list tugas akhir pertama (atau sesuai logika yang Anda inginkan)
	dataID := listTA[0].DataID
	if dataID == "" {
		http.Error(w, "No valid data ID found", http.StatusForbidden)
		return
	}

	urlTarget := fmt.Sprintf("https://siakad.ulbi.ac.id/siakad/data_bimbingan/add/%s", dataID)

	// Mengambil token dari database berdasarkan user_id
	tokenData, err := atdb.GetOneDoc[model.TokenData](config.Mongoconn, "tokens", primitive.M{"user_id": userID})
	if err != nil {
		fmt.Println("Error Fetching Token:", err)
		at.WriteJSON(w, http.StatusNotFound, "Token not found for user")
		return
	}

	// Buat payload berisi informasi token
	payload := map[string]string{
		"SIAKAD_CLOUD_ACCESS": tokenData.Token,
	}

	// Coba ambil file dari form
	file, handler, err := r.FormFile("lampiran")
	var tempFilePath string
	if err == nil {
		defer file.Close()

		// Simpan file ke direktori sementara
		tempFilePath = filepath.Join(os.TempDir(), handler.Filename)
		tempFile, err := os.Create(tempFilePath)
		if err != nil {
			http.Error(w, "Error creating temp file", http.StatusInternalServerError)
			return
		}
		defer tempFile.Close()

		_, err = io.Copy(tempFile, file)
		if err != nil {
			http.Error(w, "Error saving temp file", http.StatusInternalServerError)
			return
		}
	}

	// Ambil data dari form dan masukkan ke struct Bimbingan
	bimbingan := model.Bimbingan{
		Bimbinganke:    r.FormValue("bimbinganke"),
		NIP:            r.FormValue("nip"),
		TglBimbingan:   r.FormValue("tglbimbingan"),
		TopikBimbingan: r.FormValue("topikbimbingan"),
		Bahasan:        r.FormValue("bahasan"),
		Link:           r.FormValue("link[]"),
		Key:            r.FormValue("key"),
		Act:            r.FormValue("act"),
		Lampiran:       tempFilePath,
	}

	formData := map[string]string{
		"bimbinganke":    bimbingan.Bimbinganke,
		"nip":            bimbingan.NIP,
		"tglbimbingan":   bimbingan.TglBimbingan,
		"topikbimbingan": bimbingan.TopikBimbingan,
		"bahasan":        bimbingan.Bahasan,
		"link[]":         bimbingan.Link,
		"key":            bimbingan.Key,
		"act":            bimbingan.Act,
	}

	fileFieldName := "lampiran"
	filePath := bimbingan.Lampiran

	// Jika tidak ada file yang diunggah, kosongkan filePath
	if tempFilePath == "" {
		fileFieldName = ""
		filePath = ""
	}

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

	// Buat respons sukses berisi data bimbingan yang ditambahkan
	responseData := map[string]interface{}{
		"status":  "success",
		"message": "Data berhasil ditambahkan",
		"data":    bimbingan,
	}

	at.WriteJSON(w, http.StatusOK, responseData)
}

func GetDosen(respw http.ResponseWriter, req *http.Request) {
	urlTarget := "https://siakad.ulbi.ac.id/siakad/data_pegawai"

	// Mengambil user_id dari header
	userID := req.Header.Get("user_id")
	if userID == "" {
		http.Error(respw, "No valid user ID found", http.StatusForbidden)
		return
	}

	// Mengambil token dari database berdasarkan user_id
	tokenData, err := atdb.GetOneDoc[model.TokenData](config.Mongoconn, "tokens", primitive.M{"user_id": userID})
	if err != nil {
		fmt.Println("Error Fetching Token:", err)
		at.WriteJSON(respw, http.StatusNotFound, "Token not found for user")
		return
	}

	// Buat payload berisi informasi token
	payload := map[string]string{
		"SIAKAD_CLOUD_ACCESS": tokenData.Token,
	}

	doc, err := api.GetData(urlTarget, payload, nil)
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

	// Konversi ke JSON dan kirimkan sebagai respon
	at.WriteJSON(respw, http.StatusOK, dosen)
}

func NotFound(respw http.ResponseWriter, req *http.Request) {
	var resp model.Response
	resp.Response = "Not Found"
	at.WriteJSON(respw, http.StatusNotFound, resp)
}
