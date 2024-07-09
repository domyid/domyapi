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
	"time"

	"github.com/PuerkitoBio/goquery"
	config "github.com/domyid/domyapi/config"
	at "github.com/domyid/domyapi/helper/at"
	api "github.com/domyid/domyapi/helper/atapi"
	atdb "github.com/domyid/domyapi/helper/atdb"
	model "github.com/domyid/domyapi/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetMahasiswa handles the request to get Mahasiswa data
func GetMahasiswa(respw http.ResponseWriter, req *http.Request) {
	// Mengambil nohp dari header
	noHp := req.Header.Get("nohp")
	if noHp == "" {
		http.Error(respw, "No valid phone number found", http.StatusForbidden)
		return
	}

	// Mengambil token dari database berdasarkan nohp
	tokenData, err := atdb.GetOneDoc[model.TokenData](config.Mongoconn, "tokens", primitive.M{"nohp": noHp})
	if err != nil {
		fmt.Println("Error Fetching Token:", err)
		at.WriteJSON(respw, http.StatusNotFound, "Token tidak ditemukan! Silahkan Login Kembali")
		return
	}

	// Buat payload berisi informasi token
	cookies := map[string]string{
		"SIAKAD_CLOUD_ACCESS": tokenData.Token,
	}

	// Gunakan helper untuk mengekstrak data mahasiswa
	mahasiswa, err := api.ExtractMahasiswaData(cookies)
	if err != nil {
		at.WriteJSON(respw, http.StatusInternalServerError, err.Error())
		return
	}

	// Kembalikan instance Mahasiswa sebagai respon JSON
	at.WriteJSON(respw, http.StatusOK, mahasiswa)
}

func PostBimbinganMahasiswa(w http.ResponseWriter, r *http.Request) {
	// Mengambil nohp dari header
	noHp := r.Header.Get("nohp")
	if noHp == "" {
		http.Error(w, "No valid phone number found", http.StatusForbidden)
		return
	}

	// Mengambil token dari database berdasarkan nohp
	tokenData, err := atdb.GetOneDoc[model.TokenData](config.Mongoconn, "tokens", primitive.M{"nohp": noHp})
	if err != nil {
		fmt.Println("Error Fetching Token:", err)
		at.WriteJSON(w, http.StatusNotFound, "Token tidak ditemukan! Silahkan Login Kembali")
		return
	}

	// Memanggil fungsi helper untuk mendapatkan list tugas akhir
	listTA, err := api.FetchListTugasAkhirMahasiswa(tokenData.UserID)
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

	// Mengambil nohp dari header
	noHp := req.Header.Get("nohp")
	if noHp == "" {
		http.Error(respw, "No valid phone number found", http.StatusForbidden)
		return
	}

	// Mengambil token dari database berdasarkan nohp
	tokenData, err := atdb.GetOneDoc[model.TokenData](config.Mongoconn, "tokens", primitive.M{"nohp": noHp})
	if err != nil {
		fmt.Println("Error Fetching Token:", err)
		at.WriteJSON(respw, http.StatusNotFound, "Token tidak ditemukan! Silahkan Login Kembali")
		return
	}

	// Buat payload berisi informasi token
	cookies := map[string]string{
		"SIAKAD_CLOUD_ACCESS": tokenData.Token,
	}

	// Gunakan helper untuk mengekstrak data mahasiswa
	dosen, err := api.ExtractDosenData(cookies)
	if err != nil {
		at.WriteJSON(respw, http.StatusInternalServerError, err.Error())
		return
	}
	// Konversi ke JSON dan kirimkan sebagai respon
	at.WriteJSON(respw, http.StatusOK, dosen)
}

// Fungsi untuk menangani permintaan HTTP untuk mendapatkan data jadwal mengajar
func GetJadwalMengajar(w http.ResponseWriter, r *http.Request) {
	noHp := r.Header.Get("nohp")
	if noHp == "" {
		http.Error(w, "No valid phone number found", http.StatusForbidden)
		return
	}

	var requestData struct {
		Periode string `json:"periode"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil || requestData.Periode == "" {
		http.Error(w, "Invalid request body or periode not provided", http.StatusBadRequest)
		return
	}

	listJadwal, err := api.FetchJadwalMengajar(noHp, requestData.Periode)
	if err != nil {
		at.WriteJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	at.WriteJSON(w, http.StatusOK, listJadwal)
}

func GetRiwayatPerkuliahan(w http.ResponseWriter, r *http.Request) {
	// Mengambil nohp dari header
	noHp := r.Header.Get("nohp")
	if noHp == "" {
		http.Error(w, "No valid phone number found", http.StatusForbidden)
		return
	}

	// Mengambil token dari database berdasarkan nohp
	tokenData, err := atdb.GetOneDoc[model.TokenData](config.Mongoconn, "tokens", primitive.M{"nohp": noHp})
	if err != nil {
		fmt.Println("Error Fetching Token:", err)
		at.WriteJSON(w, http.StatusNotFound, "Token tidak ditemukan! Silahkan Login Kembali")
		return
	}

	var requestData struct {
		Periode string `json:"periode"`
		Kelas   string `json:"kelas"`
	}
	err = json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil || requestData.Periode == "" || requestData.Kelas == "" {
		http.Error(w, "Invalid request body or periode/kelas not provided", http.StatusBadRequest)
		return
	}

	// Fetch jadwal mengajar
	listJadwal, err := api.FetchJadwalMengajar(noHp, requestData.Periode)
	if err != nil || len(listJadwal) == 0 {
		at.WriteJSON(w, http.StatusNotFound, "Failed to fetch jadwal mengajar or no data found")
		return
	}

	// Cari data_id berdasarkan kelas
	var dataID string
	for _, jadwal := range listJadwal {
		if jadwal.Kelas == requestData.Kelas {
			dataID = jadwal.DataID
			break
		}
	}

	if dataID == "" {
		http.Error(w, "No valid data ID found for the given class", http.StatusNotFound)
		return
	}

	// Fetch list absensi using the data ID
	riwayatMengajar, err := api.FetchRiwayatPerkuliahan(dataID, tokenData.Token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return list absensi as JSON response
	at.WriteJSON(w, http.StatusOK, riwayatMengajar)
}

// Fungsi untuk menangani permintaan HTTP untuk mendapatkan data absensi kelas
func GetAbsensiKelas(w http.ResponseWriter, r *http.Request) {
	noHp := r.Header.Get("nohp")
	if noHp == "" {
		http.Error(w, "No valid phone number found", http.StatusForbidden)
		return
	}

	var requestData struct {
		Periode string `json:"periode"`
		Kelas   string `json:"kelas"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil || requestData.Periode == "" || requestData.Kelas == "" {
		http.Error(w, "Invalid request body or periode/kelas not provided", http.StatusBadRequest)
		return
	}

	// Fetch absensi kelas
	absensiKelas, err := api.FetchAbsensiKelas(noHp, requestData.Kelas, requestData.Periode)
	if err != nil {
		at.WriteJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	at.WriteJSON(w, http.StatusOK, absensiKelas)
}

func GetNilaiMahasiswa(w http.ResponseWriter, r *http.Request) {
	// Mengambil nohp dari header
	noHp := r.Header.Get("nohp")
	if noHp == "" {
		http.Error(w, "No valid phone number found", http.StatusForbidden)
		return
	}

	// Mengambil token dari database berdasarkan nohp
	tokenData, err := atdb.GetOneDoc[model.TokenData](config.Mongoconn, "tokens", primitive.M{"nohp": noHp})
	if err != nil {
		fmt.Println("Error Fetching Token:", err)
		at.WriteJSON(w, http.StatusNotFound, "Token tidak ditemukan! Silahkan Login Kembali")
		return
	}

	var requestData struct {
		Periode string `json:"periode"`
		Kelas   string `json:"kelas"`
	}
	err = json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil || requestData.Periode == "" || requestData.Kelas == "" {
		http.Error(w, "Invalid request body or periode/kelas not provided", http.StatusBadRequest)
		return
	}

	// Fetch jadwal mengajar
	listJadwal, err := api.FetchJadwalMengajar(noHp, requestData.Periode)
	if err != nil || len(listJadwal) == 0 {
		at.WriteJSON(w, http.StatusNotFound, "Failed to fetch jadwal mengajar or no data found")
		return
	}

	// Cari data_id berdasarkan kelas
	var dataID string
	for _, jadwal := range listJadwal {
		if jadwal.Kelas == requestData.Kelas {
			dataID = jadwal.DataID
			break
		}
	}

	if dataID == "" {
		http.Error(w, "No valid data ID found for the given class", http.StatusNotFound)
		return
	}

	// Fetch list nilai using the data ID
	listNilai, err := api.FetchNilai(dataID, tokenData.Token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return list nilai as JSON response
	at.WriteJSON(w, http.StatusOK, listNilai)
}

// Fungsi untuk menangani permintaan HTTP untuk mendapatkan data riwayat perkuliahan, absensi kelas, dan nilai mahasiswa
func GetBAP(w http.ResponseWriter, r *http.Request) {
	noHp := r.Header.Get("nohp")
	if noHp == "" {
		http.Error(w, "No valid phone number found", http.StatusForbidden)
		return
	}

	// Mengambil token dari database berdasarkan nohp
	tokenData, err := atdb.GetOneDoc[model.TokenData](config.Mongoconn, "tokens", primitive.M{"nohp": noHp})
	if err != nil {
		fmt.Println("Error Fetching Token:", err)
		at.WriteJSON(w, http.StatusNotFound, "Token tidak ditemukan! Silahkan Login Kembali")
		return
	}

	var requestData struct {
		Periode string `json:"periode"`
		Kelas   string `json:"kelas"`
	}
	err = json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil || requestData.Periode == "" || requestData.Kelas == "" {
		http.Error(w, "Invalid request body or periode/kelas not provided", http.StatusBadRequest)
		return
	}

	// Fetch jadwal mengajar
	listJadwal, err := api.FetchJadwalMengajar(noHp, requestData.Periode)
	if err != nil || len(listJadwal) == 0 {
		at.WriteJSON(w, http.StatusNotFound, "Failed to fetch jadwal mengajar or no data found")
		return
	}

	// Cari data_id berdasarkan kelas dan tambahkan informasi jadwal ke dalam hasil
	var dataID, kode, mataKuliah, sks, smt, kelas string
	for _, jadwal := range listJadwal {
		if jadwal.Kelas == requestData.Kelas {
			dataID = jadwal.DataID
			kode = jadwal.Kode
			mataKuliah = jadwal.MataKuliah
			sks = jadwal.SKS
			smt = jadwal.Smt
			kelas = jadwal.Kelas
			break
		}
	}

	if dataID == "" {
		http.Error(w, "No valid data ID found for the given class", http.StatusNotFound)
		return
	}

	// Fetch list absensi using the data ID
	riwayatMengajar, err := api.FetchRiwayatPerkuliahan(dataID, tokenData.Token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Fetch absensi kelas
	absensiKelas, err := api.FetchAbsensiKelas(noHp, requestData.Kelas, requestData.Periode)
	if err != nil {
		at.WriteJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Fetch list nilai using the data ID
	listNilai, err := api.FetchNilai(dataID, tokenData.Token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Gabungkan hasil riwayat mengajar, absensi kelas, nilai mahasiswa, dan informasi jadwal
	result := model.BAP{
		Kode:            kode,
		MataKuliah:      mataKuliah,
		SKS:             sks,
		SMT:             smt,
		Kelas:           kelas,
		RiwayatMengajar: riwayatMengajar,
		AbsensiKelas:    absensiKelas,
		ListNilai:       listNilai,
	}

	// Return combined results as JSON response
	at.WriteJSON(w, http.StatusOK, result)
}

func GetListTugasAkhirMahasiswa(respw http.ResponseWriter, req *http.Request) {
	// Mengambil no_hp dari header
	noHp := req.Header.Get("nohp")
	if noHp == "" {
		http.Error(respw, "No valid no_hp found", http.StatusForbidden)
		return
	}

	// Mengambil token dari database berdasarkan no_hp
	tokenData, err := atdb.GetOneDoc[model.TokenData](config.Mongoconn, "tokens", primitive.M{"nohp": noHp})
	if err != nil {
		fmt.Println("Error Fetching Token:", err)
		at.WriteJSON(respw, http.StatusNotFound, "Token tidak ditemukan! Silahkan Login Kembali")
		return
	}

	// Memanggil fungsi helper untuk mendapatkan list tugas akhir semua mahasiswa
	listTA, err := api.FetchListTugasAkhirMahasiswa(tokenData.NoHp)
	if err != nil || len(listTA) == 0 {
		at.WriteJSON(respw, http.StatusNotFound, "Token tidak ditemukan! Silahkan Login Kembali")
		return
	}

	// Kembalikan daftar TA sebagai respon JSON
	at.WriteJSON(respw, http.StatusOK, listTA)
}

func GetListBimbinganMahasiswa(w http.ResponseWriter, r *http.Request) {
	// Mengambil nohp dari header
	noHp := r.Header.Get("nohp")
	if noHp == "" {
		http.Error(w, "No valid phone number found", http.StatusForbidden)
		return
	}

	var requestData struct {
		NIM string `json:"nim"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil || requestData.NIM == "" {
		http.Error(w, "Invalid request body or NIM not provided", http.StatusBadRequest)
		return
	}

	// Mengambil token dari database berdasarkan nohp
	tokenData, err := atdb.GetOneDoc[model.TokenData](config.Mongoconn, "tokens", primitive.M{"nohp": noHp})
	if err != nil {
		fmt.Println("Error Fetching Token:", err)
		at.WriteJSON(w, http.StatusNotFound, "Token tidak ditemukan! Silahkan Login Kembali")
		return
	}

	// Get data ID from Tugas Akhir list
	dataID, err := api.GetDataIDFromTugasAkhir(noHp, requestData.NIM)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Fetch list bimbingan using the helper function
	listBimbingan, err := api.FetchListBimbingan(dataID, tokenData.Token)
	if err != nil {
		at.WriteJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Kembalikan daftar bimbingan sebagai respon JSON
	at.WriteJSON(w, http.StatusOK, listBimbingan)
}

// ApproveBimbingan handles the request to approve the Bimbingan
func ApproveBimbingan(w http.ResponseWriter, r *http.Request) {
	noHp := r.Header.Get("nohp")
	if noHp == "" {
		http.Error(w, "No valid phone number found", http.StatusForbidden)
		return
	}

	var requestData struct {
		NIM   string `json:"nim"`
		Topik string `json:"topik"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil || requestData.NIM == "" || requestData.Topik == "" {
		http.Error(w, "Invalid request body or NIM/Topik not provided", http.StatusBadRequest)
		return
	}

	tokenData, err := atdb.GetOneDoc[model.TokenData](config.Mongoconn, "tokens", primitive.M{"nohp": noHp})
	if err != nil {
		fmt.Println("Error Fetching Token:", err)
		at.WriteJSON(w, http.StatusNotFound, "Token tidak ditemukan! Silahkan Login Kembali")
		return
	}

	payload := map[string]string{
		"SIAKAD_CLOUD_ACCESS": tokenData.Token,
	}

	listTA, err := api.FetchListTugasAkhirMahasiswa(tokenData.NoHp)
	if err != nil || len(listTA) == 0 {
		at.WriteJSON(w, http.StatusNotFound, "Failed to fetch Tugas Akhir or no data found")
		return
	}

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

	urlTarget := fmt.Sprintf("https://siakad.ulbi.ac.id/siakad/list_bimbingan/%s?t=%d", dataID, time.Now().Unix())

	doc, err := api.GetData(urlTarget, payload, nil)
	if err != nil {
		at.WriteJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	var bimbinganID string
	doc.Find("tbody tr").Each(func(i int, s *goquery.Selection) {
		topik := strings.TrimSpace(s.Find("td").Eq(3).Text())
		disetujui := s.Find("td").Eq(4).Text()
		if topik == requestData.Topik && disetujui == "" {
			bimbinganID, _ = s.Find("td").Eq(5).Find("button").Attr("data-id")
			return
		}
	})

	if bimbinganID == "" {
		http.Error(w, "No valid data ID found for the provided topik", http.StatusForbidden)
		return
	}

	getURL := fmt.Sprintf("https://siakad.ulbi.ac.id/siakad/data_bimbingan/edit/%s", bimbinganID)

	doc, err = api.GetData(getURL, payload, nil)
	if err != nil {
		at.WriteJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	bimbinganke := doc.Find("input[name='bimbinganke']").AttrOr("value", "")
	nip := doc.Find("select[name='nip'] option[selected]").AttrOr("value", "")
	tglbimbingan := doc.Find("input[name='tglbimbingan']").AttrOr("value", "")
	topikbimbingan := doc.Find("input[name='topikbimbingan']").AttrOr("value", "")

	form := url.Values{}
	form.Add("bimbinganke", bimbinganke)
	form.Add("nip", nip)
	form.Add("tglbimbingan", tglbimbingan)
	form.Add("topikbimbingan", topikbimbingan)
	form.Add("disetujui", "1")

	postURL := fmt.Sprintf("https://siakad.ulbi.ac.id/siakad/data_bimbingan/edit/%s", bimbinganID)

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
