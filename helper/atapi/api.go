package domyApi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"

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

func GetwithUrl(urltarget string, headers map[string]string) (result []byte, err error) {
	// Create an HTTP client
	client := &http.Client{}

	// Create a new request
	req, err := http.NewRequest("GET", urltarget, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
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

// GetData fetches data from the specified URL using the provided cookies.
func GetData(url string, cookies map[string]string, headers map[string]string) (*goquery.Document, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	for name, value := range cookies {
		req.AddCookie(&http.Cookie{
			Name:  name,
			Value: value,
		})
	}

	for name, value := range headers {
		req.Header.Set(name, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to create document from response: %w", err)
	}

	return doc, nil
}

// GetDataPOST fetches data from the specified URL using the provided cookies with POST method.
func GetDataPOST(url string, cookies map[string]string, formData url.Values, headers map[string]string) (*goquery.Document, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	for name, value := range cookies {
		req.AddCookie(&http.Cookie{
			Name:  name,
			Value: value,
		})
	}

	for name, value := range headers {
		req.Header.Set(name, value)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to create document from response: %w", err)
	}

	return doc, nil
}

// PostDataToURL sends a POST request with form data and file to the specified URL.
func PostData(urltarget string, cookies map[string]string, formData map[string]string, fileFieldName, filePath string) (*http.Response, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	for key, val := range formData {
		if err := writer.WriteField(key, val); err != nil {
			return nil, fmt.Errorf("failed to write field %s: %w", key, err)
		}
	}

	if filePath != "" {
		file, err := os.Open(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to open file: %w", err)
		}
		defer file.Close()

		part, err := writer.CreateFormFile(fileFieldName, file.Name())
		if err != nil {
			return nil, fmt.Errorf("failed to create form file: %w", err)
		}

		_, err = io.Copy(part, file)
		if err != nil {
			return nil, fmt.Errorf("failed to copy file content: %w", err)
		}
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close writer: %w", err)
	}

	req, err := http.NewRequest("POST", urltarget, &body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	for name, value := range cookies {
		req.AddCookie(&http.Cookie{
			Name:  name,
			Value: value,
		})
	}

	client := &http.Client{}
	return client.Do(req)
}
