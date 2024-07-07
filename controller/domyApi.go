package domyApi

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
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
		at.WriteJSON(respw, http.StatusNotFound, "Token tidak ditemukan! Silahkan Login Kembali")
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
		at.WriteJSON(w, http.StatusNotFound, "Token tidak ditemukan! Silahkan Login Kembali")
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

func GetListBimbinganMahasiswa(w http.ResponseWriter, r *http.Request) {
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

	urlTarget := fmt.Sprintf("https://siakad.ulbi.ac.id/siakad/list_bimbingan/%s", dataID)

	// Mengambil token dari database berdasarkan user_id
	tokenData, err := atdb.GetOneDoc[model.TokenData](config.Mongoconn, "tokens", primitive.M{"user_id": userID})
	if err != nil {
		fmt.Println("Error Fetching Token:", err)
		at.WriteJSON(w, http.StatusNotFound, "Token tidak ditemukan! Silahkan Login Kembali")
		return
	}

	// Buat payload berisi informasi token
	payload := map[string]string{
		"SIAKAD_CLOUD_ACCESS": tokenData.Token,
	}

	// Mengirim permintaan untuk mengambil data list bimbingan
	doc, err := api.GetData(urlTarget, payload, nil)
	if err != nil {
		at.WriteJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Ekstrak informasi dari respon
	var listBimbingan []model.ListBimbingan
	doc.Find("tbody tr").Each(func(i int, s *goquery.Selection) {
		no := strings.TrimSpace(s.Find("td").Eq(0).Text())
		tanggal := strings.TrimSpace(s.Find("td").Eq(1).Text())
		dosenPembimbing := strings.TrimSpace(s.Find("td").Eq(2).Text())
		topik := strings.TrimSpace(s.Find("td").Eq(3).Text())
		disetujui := s.Find("td").Eq(4).Find("i").HasClass("fa-check")
		dataID, _ := s.Find("td").Eq(5).Find("button").Attr("data-id")

		listbimbingan := model.ListBimbingan{
			No:              no,
			Tanggal:         tanggal,
			DosenPembimbing: dosenPembimbing,
			Topik:           topik,
			Disetujui:       disetujui,
			DataID:          dataID,
		}
		listBimbingan = append(listBimbingan, listbimbingan)
	})

	// Kembalikan daftar bimbingan sebagai respon JSON
	at.WriteJSON(w, http.StatusOK, listBimbingan)
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
		at.WriteJSON(respw, http.StatusNotFound, "Token tidak ditemukan! Silahkan Login Kembali")
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

func GetListTugasAkhirAllMahasiswa(respw http.ResponseWriter, req *http.Request) {
	// Mengambil user_id dari header
	userID := req.Header.Get("user_id")
	if userID == "" {
		http.Error(respw, "No valid user ID found", http.StatusForbidden)
		return
	}

	// Memanggil fungsi helper untuk mendapatkan list tugas akhir semua mahasiswa
	listTA, err := api.FetchListTugasAkhirAllMahasiswa(userID)
	if err != nil || len(listTA) == 0 {
		at.WriteJSON(respw, http.StatusNotFound, "Failed to fetch Tugas Akhir or no data found")
		return
	}

	// Kembalikan daftar TA sebagai respon JSON
	at.WriteJSON(respw, http.StatusOK, listTA)
}

func GetListBimbinganMahasiswabyNim(w http.ResponseWriter, r *http.Request) {
	// Mengambil user_id dari header
	userID := r.Header.Get("user_id")
	if userID == "" {
		http.Error(w, "No valid user ID found", http.StatusForbidden)
		return
	}

	// Mengambil NIM dari body request
	var requestData struct {
		NIM string `json:"nim"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if requestData.NIM == "" {
		http.Error(w, "No valid NIM found", http.StatusBadRequest)
		return
	}

	// Memanggil fungsi helper untuk mendapatkan list tugas akhir semua mahasiswa
	listTA, err := api.FetchListTugasAkhirAllMahasiswa(userID)
	if err != nil || len(listTA) == 0 {
		at.WriteJSON(w, http.StatusNotFound, "Failed to fetch Tugas Akhir or no data found")
		return
	}

	// Cari data_id berdasarkan NIM
	var dataID string
	for _, ta := range listTA {
		if ta.NIM == requestData.NIM {
			dataID = ta.DataID
			break
		}
	}
	if dataID == "" {
		http.Error(w, "No valid data ID found for the given NIM", http.StatusNotFound)
		return
	}

	urlTarget := fmt.Sprintf("https://siakad.ulbi.ac.id/siakad/list_bimbingan/%s", dataID)

	// Mengambil token dari database berdasarkan user_id
	tokenData, err := atdb.GetOneDoc[model.TokenData](config.Mongoconn, "tokens", primitive.M{"user_id": userID})
	if err != nil {
		fmt.Println("Error Fetching Token:", err)
		at.WriteJSON(w, http.StatusNotFound, "Token tidak ditemukan! Silahkan Login Kembali")
		return
	}

	// Buat payload berisi informasi token
	payload := map[string]string{
		"SIAKAD_CLOUD_ACCESS": tokenData.Token,
	}

	// Mengirim permintaan untuk mengambil data list bimbingan
	doc, err := api.GetData(urlTarget, payload, nil)
	if err != nil {
		at.WriteJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Ekstrak informasi dari respon
	var listBimbingan []model.ListBimbingan
	doc.Find("tbody tr").Each(func(i int, s *goquery.Selection) {
		no := strings.TrimSpace(s.Find("td").Eq(0).Text())
		tanggal := strings.TrimSpace(s.Find("td").Eq(1).Text())
		dosenPembimbing := strings.TrimSpace(s.Find("td").Eq(2).Text())
		topik := strings.TrimSpace(s.Find("td").Eq(3).Text())
		disetujui := s.Find("td").Eq(4).Find("i").HasClass("fa-check")
		dataID, _ := s.Find("td").Eq(5).Find("button").Attr("data-id")

		listbimbingan := model.ListBimbingan{
			No:              no,
			Tanggal:         tanggal,
			DosenPembimbing: dosenPembimbing,
			Topik:           topik,
			Disetujui:       disetujui,
			DataID:          dataID,
		}
		listBimbingan = append(listBimbingan, listbimbingan)
	})

	// Kembalikan daftar bimbingan sebagai respon JSON
	at.WriteJSON(w, http.StatusOK, listBimbingan)
}

// Fungsi untuk mengambil data bimbingan dari URL dan mengupdate field "disetujui"
func UpdateBimbinganDisetujui(w http.ResponseWriter, r *http.Request) {
	// Mengambil user_id dari header
	userID := r.Header.Get("user_id")
	if userID == "" {
		http.Error(w, "No valid user ID found", http.StatusForbidden)
		return
	}

	// Mengambil data ID dari body request
	var requestData struct {
		Topik string `json:"topik"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil || requestData.Topik == "" {
		http.Error(w, "Invalid request body or topik not provided", http.StatusBadRequest)
		return
	}

	// Mengambil token dari database berdasarkan user_id
	tokenData, err := atdb.GetOneDoc[model.TokenData](config.Mongoconn, "tokens", primitive.M{"user_id": userID})
	if err != nil {
		fmt.Println("Error Fetching Token:", err)
		at.WriteJSON(w, http.StatusNotFound, "Token tidak ditemukan! Silahkan Login Kembali")
		return
	}

	// Buat payload berisi informasi token
	payload := map[string]string{
		"SIAKAD_CLOUD_ACCESS": tokenData.Token,
	}

	// URL untuk mendapatkan data bimbingan
	listURL := "https://siakad.ulbi.ac.id/siakad/list_bimbingan"

	// Mengirim permintaan GET untuk mendapatkan data list bimbingan
	doc, err := api.GetData(listURL, payload, nil)
	if err != nil {
		at.WriteJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Mencari topik bimbingan yang belum disetujui
	var dataID string
	doc.Find("tbody tr").Each(func(i int, s *goquery.Selection) {
		topik := strings.TrimSpace(s.Find("td").Eq(3).Text())
		disetujui := s.Find("td").Eq(4).Find("i").HasClass("fa-check")
		if topik == requestData.Topik && !disetujui {
			dataID, _ = s.Find("td").Eq(5).Find("button").Attr("data-id")
			return
		}
	})

	if dataID == "" {
		http.Error(w, "No valid data ID found for the provided topik", http.StatusForbidden)
		return
	}

	// URL untuk mendapatkan data bimbingan tertentu
	getURL := fmt.Sprintf("https://siakad.ulbi.ac.id/siakad/data_bimbingan/edit/%s", dataID)

	// Mengirim permintaan GET untuk mendapatkan data bimbingan
	doc, err = api.GetData(getURL, payload, nil)
	if err != nil {
		at.WriteJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Ekstrak informasi dari respon
	bimbinganke := doc.Find("input[name='bimbinganke']").AttrOr("value", "")
	nip := doc.Find("input[name='nip']").AttrOr("value", "")
	tglbimbingan := doc.Find("input[name='tglbimbingan']").AttrOr("value", "")
	topikbimbingan := doc.Find("input[name='topikbimbingan']").AttrOr("value", "")
	bahasan := doc.Find("textarea[name='bahasan']").Text()
	link := doc.Find("input[name='link[]']").AttrOr("value", "")
	key := doc.Find("input[name='key']").AttrOr("value", "")
	act := doc.Find("input[name='act']").AttrOr("value", "")

	// Membuat form data untuk dikirimkan
	form := url.Values{}
	form.Add("bimbinganke", bimbinganke)
	form.Add("nip", nip)
	form.Add("tglbimbingan", tglbimbingan)
	form.Add("topikbimbingan", topikbimbingan)
	form.Add("disetujui", "1")
	form.Add("bahasan", bahasan)
	form.Add("link[]", link)
	form.Add("key", key)
	form.Add("act", act)

	// URL untuk mengupdate data bimbingan
	postURL := fmt.Sprintf("https://siakad.ulbi.ac.id/siakad/data_bimbingan/edit/%s", dataID)

	// Mengirim permintaan POST untuk memperbarui data bimbingan
	req, err := http.NewRequest("POST", postURL, strings.NewReader(form.Encode()))
	if err != nil {
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Cookie", fmt.Sprintf("SIAKAD_CLOUD_ACCESS=%s", tokenData.Token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Error sending request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusSeeOther && resp.StatusCode != http.StatusOK {
		at.WriteJSON(w, resp.StatusCode, "unexpected status code")
		return
	}

	// Buat respons sukses berisi data bimbingan yang diperbarui
	responseData := map[string]interface{}{
		"status":  "success",
		"message": "Bimbingan berhasil di approve!",
	}

	at.WriteJSON(w, http.StatusOK, responseData)
}

func NotFound(respw http.ResponseWriter, req *http.Request) {
	var resp model.Response
	resp.Response = "Not Found"
	at.WriteJSON(respw, http.StatusNotFound, resp)
}
