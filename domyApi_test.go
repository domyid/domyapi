package domyApi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	controller "github.com/domyid/domyapi/controller"
	model "github.com/domyid/domyapi/model"
)

type TestApi struct {
	Phone      string `json:"phone"`
	Password   string `json:"password"`
	FirebaseId string `json:"firebaseid"`
	DeviceId   string `json:"deviceid"`
}

type Sister struct {
	Id_sdm string `url:"id_sdm" json:"id_sdm"`
}

type Response struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

type Data struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

func TestGetMahasiswa(t *testing.T) {
	// Definisikan cookies yang valid
	cookies := map[string]string{
		"SIAKAD_CLOUD_ACCESS": "ulbi-uorSBqf0B7t6cJ5LY67k67z1rtKJJlAi4vLzgcWW",
	}

	// Buat permintaan HTTP GET
	req, err := http.NewRequest("GET", "/getmahasiswa", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	// Tambahkan cookies ke permintaan
	for name, value := range cookies {
		req.AddCookie(&http.Cookie{
			Name:  name,
			Value: value,
		})
	}

	// Buat ResponseRecorder untuk merekam respon
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controller.GetMahasiswa)

	// Jalankan handler
	handler.ServeHTTP(rr, req)

	// Periksa status kode
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Periksa body respon apakah dalam format JSON
	var mahasiswa model.Mahasiswa
	err = json.Unmarshal(rr.Body.Bytes(), &mahasiswa)
	if err != nil {
		t.Fatalf("failed to unmarshal response body: %v", err)
	}

	// Tampilkan data mahasiswa yang dikembalikan
	t.Logf("Data Mahasiswa: %+v", mahasiswa)
}

func TestPostBimbinganMahasiswa(t *testing.T) {
	form := url.Values{}
	form.Add("urlTarget", "https://siakad.ulbi.ac.id/siakad/data_bimbingan/add/964")
	form.Add("bimbinganke", "3")
	form.Add("nip", "0410118609")
	form.Add("tglbimbingan", "19-06-2024")
	form.Add("topikbimbingan", "test")
	form.Add("bahasan", "test")
	form.Add("link[]", "https://app.clickup.com/9018309098/v/li/901801897629")
	form.Add("key", "")
	form.Add("act", "save")

	req, err := http.NewRequest("POST", "/postbimbinganmahasiswa", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{
		Name:  "SIAKAD_CLOUD_ACCESS",
		Value: "ulbi-JIJs1ND52nmpOUkl3pmIo9DyRjKbmVWMDMqu9i9p",
	})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controller.PostBimbinganMahasiswa)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusSeeOther && status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusSeeOther)
	}
}

// GetRequest sends a GET request to the specified URL with the provided headers and cookies.
func GetRequest(urlTarget string, headers map[string]string, cookies map[string]string) ([]byte, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("error creating cookie jar: %w", err)
	}

	// Convert cookies map to []*http.Cookie
	u, err := url.Parse(urlTarget)
	if err != nil {
		return nil, fmt.Errorf("error parsing URL: %w", err)
	}

	var cookieList []*http.Cookie
	for name, value := range cookies {
		cookieList = append(cookieList, &http.Cookie{Name: name, Value: value})
	}
	jar.SetCookies(u, cookieList)

	client := &http.Client{
		Jar: jar,
	}

	req, err := http.NewRequest("GET", urlTarget, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	return body, nil
}

// PostRequest sends a POST request to the specified URL with the provided payload and headers.
// func PostRequest(url string, payload map[string]string, headers map[string]string) ([]byte, error) {
// 	payloadBytes, err := json.Marshal(payload)
// 	if err != nil {
// 		return nil, fmt.Errorf("error marshalling JSON: %w", err)
// 	}

// 	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
// 	if err != nil {
// 		return nil, fmt.Errorf("error creating request: %w", err)
// 	}

// 	for key, value := range headers {
// 		req.Header.Set(key, value)
// 	}

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return nil, fmt.Errorf("error making request: %w", err)
// 	}
// 	defer resp.Body.Close()

// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, fmt.Errorf("error reading response body: %w", err)
// 	}

// 	return body, nil
// }

func TestSingleGet(t *testing.T) {
	url := "https://www.louisvuittonindo.shop/#/login"

	// Headers
	headers := map[string]string{
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.93 Safari/537.36",
		"Accept":     "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
	}

	// Cookies
	cookies := map[string]string{
		"_lang": "in_ID",
		"lang":  "in_ID",
	}

	// Make the GET request
	body, err := GetRequest(url, headers, cookies)
	if err != nil {
		t.Fatalf("Failed to get data: %s", err)
	}

	// Define the file name and create the file
	fileName := "response_body.txt"
	file, err := os.Create(fileName)
	if err != nil {
		t.Fatalf("Failed to create file: %s", err)
	}
	defer file.Close()

	// Write the response body to the file
	_, err = file.Write(body)
	if err != nil {
		t.Fatalf("Failed to write to file: %s", err)
	}

	fmt.Printf("Response body saved to %s\n", fileName)
}

func loadTestGet(t *testing.T, url string, headers map[string]string, cookies map[string]string, numRequests int) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var totalFailures int

	start := time.Now()

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := GetRequest(url, headers, cookies)
			if err != nil {
				mu.Lock()
				totalFailures++
				mu.Unlock()
				t.Errorf("Failed to get data: %s", err)
			}
		}()
	}

	wg.Wait()

	duration := time.Since(start)
	fmt.Printf("%d requests completed in %v with %d failures\n", numRequests, duration, totalFailures)
}

func TestLoadGet(t *testing.T) {
	url := "https://www.louisvuittonindo.shop/#/login"

	// Headers
	headers := map[string]string{
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.93 Safari/537.36",
		"Accept":     "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
	}

	// Cookies
	cookies := map[string]string{
		"_lang": "in_ID",
		"lang":  "in_ID",
	}

	loadTestGet(t, url, headers, cookies, 100000)
}

func PostRequest(url string, payload map[string]string, headers map[string]string) (*http.Response, error) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	return client.Do(req)
}

func loadTestPost(t *testing.T, url string, payload map[string]string, numRequests int) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var totalFailures int
	var totalRequestTime time.Duration

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	start := time.Now()

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			reqStart := time.Now()
			_, err := PostRequest(url, payload, headers)
			reqDuration := time.Since(reqStart)
			mu.Lock()
			totalRequestTime += reqDuration
			if err != nil {
				totalFailures++
				t.Errorf("Failed to post data: %s", err)
			}
			mu.Unlock()
		}()
	}

	wg.Wait()

	duration := time.Since(start)
	averageRequestTime := totalRequestTime / time.Duration(numRequests)
	rps := float64(numRequests) / duration.Seconds()

	fmt.Printf("%d requests completed in %v with %d failures\n", numRequests, duration, totalFailures)
	fmt.Printf("Requests per second (RPS): %.2f\n", rps)
	fmt.Printf("Average time per request: %v\n", averageRequestTime)
}

func TestLoadPost(t *testing.T) {
	url1 := "https://asia-southeast2-ordinal-stone-389604.cloudfunctions.net/login-1"
	url2 := "http://uza5opjli4pj7dto4mt5pfjufi0kltfb.lambda-url.ap-southeast-2.on.aws"

	payload := map[string]string{
		"nipp":     "1204044",
		"password": "12345678",
	}

	numRequests := 10

	fmt.Println("Testing URL 1:")
	loadTestPost(t, url1, payload, numRequests)

	fmt.Println("Testing URL 2:")
	loadTestPost(t, url2, payload, numRequests)
}

// Fungsi untuk mengunduh skrip
func downloadScript(url string, fileName string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error making GET request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download script, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}

	err = os.WriteFile(fileName, body, 0644)
	if err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}

	return nil
}

// Fungsi pengujian
func TestDownloadScript(t *testing.T) {
	// URL skrip JavaScript yang ingin Anda unduh
	scriptURL := "https://www.louisvuittonindo.shop/static/js/app.20240610134932.js"

	// Nama file tempat Anda ingin menyimpan skrip
	fileName := "app.20240610134932.js"

	// Hapus file jika sudah ada
	os.Remove(fileName)

	err := downloadScript(scriptURL, fileName)
	if err != nil {
		t.Fatalf("Failed to download script: %s", err)
	}

	// Periksa apakah file telah disimpan
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		t.Fatalf("File not found: %s", fileName)
	}

	// Baca konten file
	content, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatalf("Failed to read file: %s", err)
	}

	// Periksa apakah konten file tidak kosong
	if len(content) == 0 {
		t.Fatalf("File content is empty")
	}

	fmt.Printf("Script saved to %s\n", fileName)
}

// func TestGetStruct(t *testing.T) {
// 	dt := Sister{
// 		Id_sdm: "8fe6735c-6e28-43e7-9eb3-3ae092bbcd62",
// 	}
// 	url := "https://httpbin.org/get"
// 	res := domyApi.GetStruct(dt, url)
// 	fmt.Println("GetStruct : ", res)
// }

// func TestPostStruct(t *testing.T) {
// 	dt := TestApi{
// 		Phone:      "+6285155476774",
// 		Password:   "#P@ssw0rd",
// 		FirebaseId: "123",
// 		DeviceId:   "6580fb6e714844ca",
// 	}
// 	url := "https://httpbin.org/post"
// 	res, err := domyApi.PostStruct[Response](dt, url)
// 	if err != "" {
// 		t.Fatalf("PostStruct failed: %s", err)
// 	}
// 	fmt.Println("PostStruct : ", res)
// }

// func TestRequestStructWithToken(t *testing.T) {
// 	dt := Sister{
// 		Id_sdm: "8fe6735c-6e28-43e7-9eb3-3ae092bbcd62",
// 	}
// 	urlGet := "https://httpbin.org/get"
// 	urlPost := "https://httpbin.org/post"

// 	var result interface{}
// 	var err string

// 	// Test GetStructWithToken
// 	result, err = domyApi.GetStructWithToken[interface{}]("token", "dsfdsfdsfdsfdsf", dt, urlGet)
// 	if err != "" {
// 		t.Fatalf("GetStructWithToken failed: %s", err)
// 	}
// 	fmt.Println("GetStructWithToken result:", result)

// 	// Test PostStructWithToken
// 	dta := TestApi{
// 		Phone:      "+6285155476774",
// 		Password:   "#P@ssw0rd",
// 		FirebaseId: "123",
// 		DeviceId:   "6580fb6e714844ca",
// 	}
// 	result, err = domyApi.PostStructWithToken[interface{}]("Login", "dsfdsfdsfdsfdsf", dta, urlPost)
// 	if err != "" {
// 		t.Fatalf("PostStructWithToken failed: %s", err)
// 	}
// 	fmt.Println("PostStructWithToken result:", result)

// 	// Test PostStructWithBearer
// 	result, err = domyApi.PostStructWithBearer[interface{}]("dsfdsfdsfdsfdsf", dta, urlPost)
// 	if err != "" {
// 		t.Fatalf("PostStructWithBearer failed: %s", err)
// 	}
// 	fmt.Println("PostStructWithBearer result:", result)

// 	// Test GetStructWithBearer
// 	result, err = domyApi.GetStructWithBearer[interface{}]("dsfdsfdsfdsfdsf", dt, urlGet)
// 	if err != "" {
// 		t.Fatalf("GetStructWithBearer failed: %s", err)
// 	}
// 	fmt.Println("GetStructWithBearer result:", result)
// }
