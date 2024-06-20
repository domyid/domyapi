package domyApi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"strings"

	"github.com/PuerkitoBio/goquery"
	config "github.com/domyid/domyapi/config"
	at "github.com/domyid/domyapi/helper/at"
	atdb "github.com/domyid/domyapi/helper/atdb"
	model "github.com/domyid/domyapi/model"
	"github.com/google/go-querystring/query"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Get(urltarget string, cookies map[string]string, headers map[string]string) (result []byte, err error) {
	// Create a cookie jar
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create cookie jar: %w", err)
	}

	// Create an HTTP client with the cookie jar
	client := &http.Client{
		Jar: jar,
	}

	// Create a new request
	req, err := http.NewRequest("GET", urltarget, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add cookies to the request
	for name, value := range cookies {
		req.AddCookie(&http.Cookie{
			Name:  name,
			Value: value,
		})
	}

	// Add additional headers to the request
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	// Make the GET request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make GET request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading data from response: %w", err)
	}

	return body, nil
}

func GetMahasiswa(urltarget string, cookies map[string]string, headers map[string]string) (string, error) {
	// Buat cookie jar
	jar, err := cookiejar.New(nil)
	if err != nil {
		return "", fmt.Errorf("failed to create cookie jar: %w", err)
	}

	// Buat HTTP client dengan cookie jar
	client := &http.Client{
		Jar: jar,
	}

	// Buat permintaan baru
	req, err := http.NewRequest("GET", urltarget, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Tambahkan cookies ke permintaan
	for name, value := range cookies {
		req.AddCookie(&http.Cookie{
			Name:  name,
			Value: value,
		})
	}

	// Tambahkan headers tambahan ke permintaan
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	// Lakukan permintaan GET
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make GET request: %w", err)
	}
	defer resp.Body.Close()

	// Parse response body dengan goquery
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to parse response body: %w", err)
	}

	// Ekstrak informasi mahasiswa dengan trim spasi
	nim := strings.TrimSpace(doc.Find("#block-nim .input-nim").Text())
	nama := strings.TrimSpace(doc.Find("#block-nama .input-nama").Text())
	programStudi := strings.TrimSpace(doc.Find("#block-idunit .input-idunit").Text())
	noHp := strings.TrimSpace(doc.Find("#block-hp .input-hp").Text())

	// Buat instance Mahasiswa
	mahasiswa := model.Mahasiswa{
		ID:           primitive.NewObjectID(),
		NIM:          nim,
		Nama:         nama,
		ProgramStudi: programStudi,
		NomorHp:      noHp,
	}

	// Cek apakah data mahasiswa sudah ada di database
	filter := bson.M{"nim": mahasiswa.NIM}
	var existingMahasiswa model.Mahasiswa
	err = config.Mongoconn.Collection("mahasiswa").FindOne(context.TODO(), filter).Decode(&existingMahasiswa)
	if err == nil {
		// Data sudah ada, tidak perlu menambah data baru
		jsonData, err := json.Marshal(existingMahasiswa)
		if err != nil {
			return "", fmt.Errorf("failed to marshal JSON: %w", err)
		}
		return string(jsonData), nil
	}

	// Simpan ke MongoDB jika data belum ada
	if _, err := atdb.InsertOneDoc(config.Mongoconn, "mahasiswa", mahasiswa); err != nil {
		return "", fmt.Errorf("failed to insert document into MongoDB: %w", err)
	}

	// Konversi ke JSON
	jsonData, err := json.Marshal(mahasiswa)
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return string(jsonData), nil
}

func GetDosen(urltarget string, cookies map[string]string, headers map[string]string) (string, error) {
	// Buat cookie jar
	jar, err := cookiejar.New(nil)
	if err != nil {
		return "", fmt.Errorf("failed to create cookie jar: %w", err)
	}

	// Buat HTTP client dengan cookie jar
	client := &http.Client{
		Jar: jar,
	}

	// Buat permintaan baru
	req, err := http.NewRequest("GET", urltarget, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Tambahkan cookies ke permintaan
	for name, value := range cookies {
		req.AddCookie(&http.Cookie{
			Name:  name,
			Value: value,
		})
	}

	// Tambahkan headers tambahan ke permintaan
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	// Lakukan permintaan GET
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make GET request: %w", err)
	}
	defer resp.Body.Close()

	// Parse response body dengan goquery
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to parse response body: %w", err)
	}

	// Ekstrak informasi dosen dengan trim spasi
	nip := strings.TrimSpace(doc.Find("#block-nip .input-nip").Text())
	nidn := strings.TrimSpace(doc.Find("#block-nidn .input-nidn").Text())
	nama := strings.TrimSpace(doc.Find("#block-nama .input-nama").Text())
	noHp := strings.TrimSpace(doc.Find("#block-nohp .input-nohp").Text())

	// Buat instance Dosen
	dosen := model.Dosen{
		ID:   primitive.NewObjectID(),
		NIP:  nip,
		NIDN: nidn,
		Nama: nama,
		NoHp: noHp,
	}

	// Cek apakah data dosen sudah ada di database
	filter := bson.M{"nip": dosen.NIP}
	var existingDosen model.Dosen
	err = config.Mongoconn.Collection("dosen").FindOne(context.TODO(), filter).Decode(&existingDosen)
	if err == nil {
		// Data sudah ada, tidak perlu menambah data baru
		jsonData, err := json.Marshal(existingDosen)
		if err != nil {
			return "", fmt.Errorf("failed to marshal JSON: %w", err)
		}
		return string(jsonData), nil
	}

	// Simpan ke MongoDB jika data belum ada
	if _, err := atdb.InsertOneDoc(config.Mongoconn, "dosen", dosen); err != nil {
		return "", fmt.Errorf("failed to insert document into MongoDB: %w", err)
	}

	// Konversi ke JSON
	jsonData, err := json.Marshal(dosen)
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return string(jsonData), nil
}

func GetStruct(structname interface{}, urltarget string) (errormessage string) {
	v, _ := query.Values(structname)
	resp, err := http.Get(urltarget + "?" + v.Encode())
	if err != nil {
		errormessage = "GetStruct http.get: " + err.Error()
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		errormessage = "Error Read data from Response: " + err.Error()
		return
	}

	// Print the response body for debugging
	fmt.Println("Status Code: ", resp.StatusCode)
	fmt.Println("Response Body: ", string(body)) // Print the response body

	errormessage = "Request successful"
	return
}

func GetStructWithBearer[T any](tokenbearer string, structname interface{}, urltarget string) (result T, errormessage string) {
	client := http.Client{}
	v, _ := query.Values(structname)
	req, err := http.NewRequest("GET", urltarget+"?"+v.Encode(), nil)
	if err != nil {
		errormessage = "http.NewRequest Got error : " + err.Error()
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+tokenbearer)
	resp, err := client.Do(req)
	if err != nil {
		errormessage = "client.Do(req) Error occured. Error is :" + err.Error()
		return
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		errormessage = "Error Read Data data from request : " + err.Error()
		return
	}
	if er := json.Unmarshal(respBody, &result); er != nil {
		errormessage = "Error Unmarshal from Response." + er.Error()
	}
	return
}

func GetStructWithToken[T any](tokenkey string, tokenvalue string, structname interface{}, urltarget string) (result T, errormessage string) {
	client := http.Client{}
	v, _ := query.Values(structname)
	req, err := http.NewRequest("GET", urltarget+"?"+v.Encode(), nil)
	if err != nil {
		errormessage = "http.NewRequest Got error : " + err.Error()
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add(tokenkey, tokenvalue)
	resp, err := client.Do(req)
	if err != nil {
		errormessage = "client.Do(req) Error occured. Error is :" + err.Error()
		return
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		errormessage = "Error Read Data data from request : " + err.Error()
		return
	}
	if er := json.Unmarshal(respBody, &result); er != nil {
		errormessage = "Error Unmarshal from Response." + er.Error()
	}
	return
}

func PostStruct[T any](structname interface{}, urltarget string) (result T, errormessage string) {
	mJson, _ := json.Marshal(structname)
	fmt.Println("Request JSON: ", string(mJson)) // Print request JSON
	resp, err := http.Post(urltarget, "application/json", bytes.NewBuffer(mJson))
	if err != nil {
		errormessage = "Could not make POST request to server : " + err.Error()
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		errormessage = "Error Read Data from request : " + err.Error()
		return
	}

	fmt.Println("Status Code: ", resp.StatusCode)
	fmt.Println("Response Body: ", string(body)) // Print the response body

	if resp.Header.Get("Content-Type") == "application/json" {
		if er := json.Unmarshal(body, &result); er != nil {
			errormessage = "Error Unmarshal from Response." + er.Error()
		}
	} else {
		errormessage = "Received non-JSON response: " + string(body)
	}
	return
}

func PostStructWithBearer[T any](tokenbearer string, structname interface{}, urltarget string) (result T, errormessage string) {
	client := http.Client{}
	mJson, _ := json.Marshal(structname)
	req, err := http.NewRequest("POST", urltarget, bytes.NewBuffer(mJson))
	if err != nil {
		errormessage = "http.NewRequest Got error :" + err.Error()
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+tokenbearer)
	resp, err := client.Do(req)
	if err != nil {
		errormessage = "client.Do(req) Error occured. Error is :" + err.Error()
		return
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		errormessage = "Error Read Data data from request." + err.Error()
		return
	}
	if er := json.Unmarshal(respBody, &result); er != nil {
		errormessage = "Error Unmarshal from Response." + er.Error()
	}
	return
}

func PostStructWithToken[T any](tokenkey string, tokenvalue string, structname interface{}, urltarget string) (result T, errormessage string) {
	client := http.Client{}
	mJson, _ := json.Marshal(structname)
	req, err := http.NewRequest("POST", urltarget, bytes.NewBuffer(mJson))
	if err != nil {
		errormessage = "http.NewRequest Got error :" + err.Error()
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add(tokenkey, tokenvalue)
	resp, err := client.Do(req)
	if err != nil {
		errormessage = "client.Do(req) Error occured. Error is :" + err.Error()
		return
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		errormessage = "Error Read Data data from request." + err.Error()
		return
	}
	if er := json.Unmarshal(respBody, &result); er != nil {
		errormessage = string(respBody) + "Error Unmarshal from Response : " + er.Error()
	}
	return
}

func PutStructWithBearer[T any](tokenbearer string, structname interface{}, urltarget string) (result T, errormessage string) {
	client := http.Client{}
	mJson, _ := json.Marshal(structname)
	req, err := http.NewRequest("PUT", urltarget, bytes.NewBuffer(mJson))
	if err != nil {
		errormessage = "http.NewRequest Got error :" + err.Error()
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+tokenbearer)
	resp, err := client.Do(req)
	if err != nil {
		errormessage = "client.Do(req) Error occured. Error is :" + err.Error()
		return
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		errormessage = "Error Read Data data from request : " + err.Error()
		return
	}
	if er := json.Unmarshal(respBody, &result); er != nil {
		errormessage = "Error Unmarshal from Response." + er.Error()
	}
	return
}

func NotFound(respw http.ResponseWriter, req *http.Request) {
	var resp model.Response
	resp.Response = "Not Found"
	at.WriteJSON(respw, http.StatusNotFound, resp)
}
