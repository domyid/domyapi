package domyApi

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"
	"sync"

	config "github.com/domyid/domyapi/config"
	at "github.com/domyid/domyapi/helper/at"
	api "github.com/domyid/domyapi/helper/atapi"
	atdb "github.com/domyid/domyapi/helper/atdb"
	ghp "github.com/domyid/domyapi/helper/ghupload"
	pdf "github.com/domyid/domyapi/helper/pdf"
	model "github.com/domyid/domyapi/model"
	"github.com/google/go-github/v32/github"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/oauth2"
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

var (
	poolOnce          sync.Once
	PoolStringBuilder *sync.Pool
)

func initStringBuilderPool() {
	poolOnce.Do(func() {
		PoolStringBuilder = &sync.Pool{
			New: func() interface{} {
				return new(strings.Builder)
			},
		}
	})
}

func GetBAP(w http.ResponseWriter, r *http.Request) {
	initStringBuilderPool()

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

	// Generate PDF
	buf, fileName, err := pdf.GenerateBAPPDF(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Fetch GitHub credentials from database
	gh, err := atdb.GetOneDoc[model.Ghcreates](config.Mongoconn, "github", bson.M{})
	if err != nil {
		http.Error(w, "Failed to fetch GitHub credentials from database", http.StatusInternalServerError)
		return
	}

	// Define GitHub path
	gitHubPath := "2023-2/" + fileName

	// Check if file exists in GitHub
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: gh.GitHubAccessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// Check if file exists
	_, _, _, err = client.Repositories.GetContents(ctx, "repoulbi", "buktiajar", gitHubPath, nil)
	if err == nil {
		// File already exists, build URL from repository page
		strPol := PoolStringBuilder.Get().(*strings.Builder)
		defer func() {
			strPol.Reset()
			PoolStringBuilder.Put(strPol)
		}()

		filePath := "buktiajar/2023-2/" + fileName
		filePathEncoded := base64.StdEncoding.EncodeToString([]byte(filePath))
		strPol.WriteString("https://repo.ulbi.ac.id/view/#" + filePathEncoded)
		at.WriteJSON(w, http.StatusOK, map[string]string{"url": strPol.String()})
		return
	}

	// File tidak ada di GitHub, upload PDF
	_, _, err = ghp.GithubUpload(
		gh.GitHubAccessToken,
		gh.GitHubAuthorName,
		gh.GitHubAuthorEmail,
		&multipart.FileHeader{
			Filename: fileName,
			Header:   textproto.MIMEHeader{},
			Size:     int64(buf.Len()),
		},
		"repoulbi",
		"buktiajar",
		gitHubPath,
		false,
		buf,
	)
	if err != nil {
		http.Error(w, "Failed to upload PDF to GitHub: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Build the URL from repository page again after upload
	strPol := PoolStringBuilder.Get().(*strings.Builder)
	defer func() {
		strPol.Reset()
		PoolStringBuilder.Put(strPol)
	}()

	filePath := "buktiajar/2023-2/" + fileName
	filePathEncoded := base64.StdEncoding.EncodeToString([]byte(filePath))
	strPol.WriteString("https://repo.ulbi.ac.id/view/#" + filePathEncoded)
	at.WriteJSON(w, http.StatusOK, map[string]string{"url": strPol.String()})

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

	// Fetch list bimbingan using the helper function
	listBimbingan, err := api.FetchListBimbingan(dataID, tokenData.Token)
	if err != nil {
		at.WriteJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	var bimbinganID string
	for _, bimbingan := range listBimbingan {
		if bimbingan.Topik == requestData.Topik && !bimbingan.Disetujui {
			bimbinganID = bimbingan.DataID
			break
		}
	}

	if bimbinganID == "" {
		http.Error(w, "No valid data ID found for the provided topik", http.StatusForbidden)
		return
	}

	// Mendapatkan data detail bimbingan
	editData, err := api.GetDetailBimbingan(bimbinganID, tokenData.Token)
	if err != nil {
		at.WriteJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Approve the bimbingan
	err = api.ApproveBimbingan(bimbinganID, tokenData.Token, editData)
	if err != nil {
		at.WriteJSON(w, http.StatusInternalServerError, err.Error())
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
