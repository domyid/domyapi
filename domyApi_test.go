package domyApi

import (
	"fmt"
	"testing"
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

// Single Request Testing
func TestGetMahasiswa(t *testing.T) {
	// Define the cookies
	cookies := map[string]string{
		"SIAKAD_CLOUD_ACCESS": "ulbi-A7tGVCVKioR96xeuzA5Y0eoNwGRSUlw0GxfYRi6h",
	}

	// Define additional headers
	headers := map[string]string{
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.93 Safari/537.36",
		"Accept":     "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
	}

	// Make the GET request and extract student information in JSON format
	jsonData, err := GetMahasiswa("https://siakad.ulbi.ac.id/siakad/data_mahasiswa", cookies, headers)
	if err != nil {
		t.Fatalf("Failed to get data: %s", err)
	}

	// Print the extracted JSON data for verification
	t.Logf("Extracted JSON Data: %s", jsonData)

	// Define MongoDB connection details
	mongoURI := "mongodb+srv://fahira:aqilah123@domyapi.wit3inx.mongodb.net/"
	dbName := "siakad"
	collectionName := "mahasiswa"

	// Save the extracted information to MongoDB
	err = SaveToMongoDB(jsonData, mongoURI, dbName, collectionName)
	if err != nil {
		t.Fatalf("Error saving to MongoDB: %s", err)
	}

	t.Log("Data berhasil disimpan ke MongoDB")
}

// Single Request Testing
func TestGetDosen(t *testing.T) {
	// Define the cookies
	cookies := map[string]string{
		"SIAKAD_CLOUD_ACCESS": "ulbi-9BD962EETw7IdYM2D67HxbbFOiFxxgeQrcftdWwj",
	}

	// Define additional headers
	headers := map[string]string{
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.93 Safari/537.36",
		"Accept":     "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
	}

	// Make the GET request and extract student information in JSON format
	jsonData, err := GetDosen("https://siakad.ulbi.ac.id/siakad/data_pegawai", cookies, headers)
	if err != nil {
		t.Fatalf("Failed to get data: %s", err)
	}

	// Print the extracted JSON data for verification
	t.Logf("Extracted JSON Data: %s", jsonData)

	// Define MongoDB connection details
	mongoURI := "mongodb+srv://fahira:aqilah123@domyapi.wit3inx.mongodb.net/"
	dbName := "siakad"
	collectionName := "dosen"

	// Save the extracted information to MongoDB
	err = SaveToMongoDB(jsonData, mongoURI, dbName, collectionName)
	if err != nil {
		t.Fatalf("Error saving to MongoDB: %s", err)
	}

	t.Log("Data berhasil disimpan ke MongoDB")
}

// // Single Request Testing
// func TestGet(t *testing.T) {
// 	// Define the cookies
// 	cookies := map[string]string{
// 		"SIAKAD_CLOUD_ACCESS": "ulbi-hflkskFmFT2rgoojscMRaFfKMBvSOW5m4qrDMC9Y",
// 	}

// 	// Define additional headers
// 	headers := map[string]string{
// 		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.93 Safari/537.36",
// 		"Accept":     "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
// 	}

// 	// Make the GET request
// 	body, err := Get("https://siakad.ulbi.ac.id/siakad/data_mahasiswa", cookies, headers)
// 	if err != nil {
// 		t.Fatalf("Failed to get data: %s", err)
// 	}

// 	// Define the file name and create the file
// 	fileName := "response_body.txt"
// 	file, err := os.Create(fileName)
// 	if err != nil {
// 		t.Fatalf("Failed to create file: %s", err)
// 	}
// 	defer file.Close()

// 	// Write the response body to the file
// 	_, err = file.Write(body)
// 	if err != nil {
// 		t.Fatalf("Failed to write to file: %s", err)
// 	}

// 	fmt.Printf("Response body saved to %s\n", fileName)
// }

// Multiple Request Testing
// func loadTestGet(t *testing.T, url string, cookies map[string]string, headers map[string]string, numRequests int) {
// 	var wg sync.WaitGroup

// 	start := time.Now()

// 	for i := 0; i < numRequests; i++ {
// 		wg.Add(1)
// 		go func() {
// 			defer wg.Done()
// 			_, err := Get(url, cookies, headers)
// 			if err != nil {
// 				t.Errorf("Failed to get data: %s", err)
// 			}
// 		}()
// 	}

// 	wg.Wait()

// 	duration := time.Since(start)
// 	fmt.Printf("%d requests completed in %v\n", numRequests, duration)
// }

// func TestLoad(t *testing.T) {
// 	// Define the cookies
// 	cookies := map[string]string{
// 		"PHPSESSID":      "5k34phdb336nuonu5u7j2htjdo",
// 		"PortalMHS[JWT]": "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzUxMiJ9.eyJpYXQiOjE3MTY3OTUxNjgsImp0aSI6InFhTFwvZkoxbUVCS0R3Y1wvT01GVFBpOWE3d1wvNG0rUVJ5amxVbXkyTWJrNmM9IiwiaXNzIjoiYXBwIiwibmJmIjowLCJleHAiOjE3MTY3OTgxNjgsInNlY3VyaXR5Ijp7InVzZXJuYW1lIjoiMTIwNDA0NCIsInVzZXJpZCI6IjEyMDQwNDQiLCJwYXJlbnR1c2VyaWQiOm51bGwsInVzZXJsZXZlbGlkIjotMn19.GAe691m4hfLgfT0UmoHZeK5FOXx9282AGjPGbuEIO3iwG1kA9rUyvJpy2BKSXHRbjUAf6CAydlg4xRnwpK0YPw",
// 	}

// 	// Define additional headers
// 	headers := map[string]string{
// 		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.93 Safari/537.36",
// 		"Accept":     "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
// 	}

// 	url := "https://siapmhs.ulbi.ac.id/Dashboard1"

// 	// Load test with different numbers of requests
// 	loadTestGet(t, url, cookies, headers, 1)
// 	// loadTestGet(t, url, cookies, headers, 5000)
// 	// loadTestGet(t, url, cookies, headers, 10000)
// }

func TestGetStruct(t *testing.T) {
	dt := Sister{
		Id_sdm: "8fe6735c-6e28-43e7-9eb3-3ae092bbcd62",
	}
	url := "https://httpbin.org/get"
	res := GetStruct(dt, url)
	fmt.Println("GetStruct : ", res)
}

func TestPostStruct(t *testing.T) {
	dt := TestApi{
		Phone:      "+6285155476774",
		Password:   "#P@ssw0rd",
		FirebaseId: "123",
		DeviceId:   "6580fb6e714844ca",
	}
	url := "https://httpbin.org/post"
	res, err := PostStruct[Response](dt, url)
	if err != "" {
		t.Fatalf("PostStruct failed: %s", err)
	}
	fmt.Println("PostStruct : ", res)
}

func TestRequestStructWithToken(t *testing.T) {
	dt := Sister{
		Id_sdm: "8fe6735c-6e28-43e7-9eb3-3ae092bbcd62",
	}
	urlGet := "https://httpbin.org/get"
	urlPost := "https://httpbin.org/post"

	var result interface{}
	var err string

	// Test GetStructWithToken
	result, err = GetStructWithToken[interface{}]("token", "dsfdsfdsfdsfdsf", dt, urlGet)
	if err != "" {
		t.Fatalf("GetStructWithToken failed: %s", err)
	}
	fmt.Println("GetStructWithToken result:", result)

	// Test PostStructWithToken
	dta := TestApi{
		Phone:      "+6285155476774",
		Password:   "#P@ssw0rd",
		FirebaseId: "123",
		DeviceId:   "6580fb6e714844ca",
	}
	result, err = PostStructWithToken[interface{}]("Login", "dsfdsfdsfdsfdsf", dta, urlPost)
	if err != "" {
		t.Fatalf("PostStructWithToken failed: %s", err)
	}
	fmt.Println("PostStructWithToken result:", result)

	// Test PostStructWithBearer
	result, err = PostStructWithBearer[interface{}]("dsfdsfdsfdsfdsf", dta, urlPost)
	if err != "" {
		t.Fatalf("PostStructWithBearer failed: %s", err)
	}
	fmt.Println("PostStructWithBearer result:", result)

	// Test GetStructWithBearer
	result, err = GetStructWithBearer[interface{}]("dsfdsfdsfdsfdsf", dt, urlGet)
	if err != "" {
		t.Fatalf("GetStructWithBearer failed: %s", err)
	}
	fmt.Println("GetStructWithBearer result:", result)
}
