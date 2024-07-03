package domyApi

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	model "github.com/domyid/domyapi/model"
)

func LoginAct(client http.Client, reqLogin model.RequestLoginSiakad) (*model.ResponseLogin, error) {
	loginURL := "https://siakad.ulbi.ac.id/gate/login"
	resp, err := client.Get(loginURL)
	if err != nil {
		fmt.Println("Error fetching login page:", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Read the body of the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil, err
	}

	// Parse the HTML to find the token
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		fmt.Println("Error parsing HTML:", err)
		return nil, err
	}

	// Find the token value
	token, exists := doc.Find("input[name=__token]").Attr("value")
	if !exists {
		fmt.Println("Token not found")
		return nil, errors.New("token not found")
	}

	clientID, exist := doc.Find("input[name=client_id]").Attr("value")
	if !exist {
		fmt.Println("Client ID not found")
		return nil, errors.New("client ID not found")
	}

	// Form data for login request
	formData := url.Values{
		"email":        {reqLogin.Email},
		"password":     {reqLogin.Password},
		"__token":      {token},
		"_token":       {""},
		"client_id":    {clientID},
		"redirect_uri": {"https://siakad.ulbi.ac.id/gate/authsso"},
	}

	// Create a new request for login
	req, err := http.NewRequest("POST", loginURL, strings.NewReader(formData.Encode()))
	if err != nil {
		fmt.Println("Error creating login request:", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")

	// Send the login request
	loginResp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending login request:", err)
		return nil, err
	}
	defer loginResp.Body.Close()
	// Print all headers from the login response

	redirectURL := loginResp.Header.Get("Sx-Referer")

	u, err := url.Parse(redirectURL)
	if err != nil {
		fmt.Println("Error parsing redirect URL:", err)
		return nil, err
	}
	code := u.Query().Get("code")

	if code == "" {
		fmt.Println("Code/token not found in redirect URL")
		return nil, errors.New("code/token not found in redirect URL")
	}

	if loginResp.Header.Get("Sx-Session") == "" {
		fmt.Println("No session found")
		return nil, errors.New("no session found")
	}

	result := &model.ResponseLogin{
		Session: loginResp.Header.Get("Sx-Session"),
		Code:    code,
		Role:    reqLogin.Role,
	}
	return result, nil

}

func LoginRequest(client *http.Client, userReq model.ResponseLogin) (*model.ResponseLogin, error) {
	loginURL := "https://siakad.ulbi.ac.id/siakad/login"

	// Form data for login request
	formData := url.Values{
		"oldpass":   {""},
		"newpass":   {""},
		"renewpass": {""},
		"act":       {""},
		"sessdata":  {""},
		"kodemodul": {"siakad"},
		"koderole":  {userReq.Role},
		"kodeunit":  {"55301"},
	}

	req, err := http.NewRequest("POST", loginURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, err
	}

	// Set necessary headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Referer", "https://siakad.ulbi.ac.id/gate/menu")
	req.Header.Set("Cookie", fmt.Sprintf("XSRF-TOKEN=%s", userReq.Code))

	// Add the necessary cookies
	req.AddCookie(&http.Cookie{Name: "SIAKAD_CLOUD_ACCESS", Value: userReq.Session})

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	for resp.StatusCode == http.StatusFound {
		redirectURL := resp.Header.Get("Location")
		if redirectURL == "" {
			fmt.Println("No redirect URL found")
			break
		}

		req, err = http.NewRequest("GET", redirectURL, nil)
		if err != nil {
			return nil, err
		}
		resp, err = client.Do(req)
		if err != nil {
			return nil, err
		}
	}

	result := &model.ResponseLogin{
		Code:    userReq.Code,
		Session: userReq.Session,
		Role:    userReq.Role,
	}
	return result, nil
}

func GetRefreshToken(client *http.Client, token string) (string, error) {
	homeURL := "https://siakad.ulbi.ac.id/siakad/home"

	req, err := http.NewRequest("GET", homeURL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Referer", "https://siakad.ulbi.ac.id/gate/menu")
	req.Header.Set("Cookie", fmt.Sprintf("SIAKAD_CLOUD_ACCESS=%s", token))

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Ambil cookies dari response header
	cookies := resp.Cookies()
	var newToken string
	for _, cookie := range cookies {
		if cookie.Name == "SIAKAD_CLOUD_ACCESS" {
			newToken = cookie.Value
			break
		}
	}

	// Jika tidak ada token baru di cookies, gunakan token yang lama
	if newToken == "" {
		newToken = token
	}

	return newToken, nil
}
