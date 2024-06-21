package domyApi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"

	"github.com/PuerkitoBio/goquery"
	"github.com/google/go-querystring/query"
)

func Get[T any](urltarget string) (statusCode int, result T, err error) {
	resp, err := http.Get(urltarget)
	if err != nil {
		return
	}
	statusCode = resp.StatusCode
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	if er := json.Unmarshal(body, &result); er != nil {
		return
	}
	return
}

func GetData(urltarget string, cookies map[string]string, headers map[string]string) (result []byte, err error) {
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

func PostStructWithToken[T any](tokenkey string, tokenvalue string, structname interface{}, urltarget string) (result T, err error) {
	client := http.Client{}
	mJson, _ := json.Marshal(structname)
	req, err := http.NewRequest("POST", urltarget, bytes.NewBuffer(mJson))
	if err != nil {
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add(tokenkey, tokenvalue)
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	if err = json.Unmarshal(respBody, &result); err != nil {
		rawstring := string(respBody)
		err = errors.New("Not A Valid JSON Response from " + urltarget + ". CONTENT: " + rawstring)
		return
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

func FetchDataFromURL(urltarget string, cookies map[string]string, headers map[string]string) (*goquery.Document, error) {
	// Buat cookie jar
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create cookie jar: %w", err)
	}

	// Buat HTTP client dengan cookie jar
	client := &http.Client{
		Jar: jar,
	}

	// Buat permintaan baru
	req, err := http.NewRequest("GET", urltarget, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
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
		return nil, fmt.Errorf("failed to make GET request: %w", err)
	}
	defer resp.Body.Close()

	// Baca isi body respon
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Log isi HTML
	log.Printf("Fetched HTML: %s", string(body))

	// Parse response body dengan goquery
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response body: %w", err)
	}

	return doc, nil
}
