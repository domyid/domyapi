package domyApi

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/cookiejar"
	"time"

	config "github.com/domyid/domyapi/config"
	at "github.com/domyid/domyapi/helper/at"
	helper "github.com/domyid/domyapi/helper/atapi"
	atdb "github.com/domyid/domyapi/helper/atdb"
	model "github.com/domyid/domyapi/model"
)

func LoginSiakad(w http.ResponseWriter, req *http.Request) {
	jar, _ := cookiejar.New(nil)

	// Create a new HTTP client with the cookie jar
	client := &http.Client{
		Jar: jar,
	}

	var reqLogin model.RequestLoginSiakad

	if err := json.NewDecoder(req.Body).Decode(&reqLogin); err != nil {
		at.WriteJSON(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	resp, err := helper.LoginAct(*client, reqLogin)
	if err != nil {
		at.WriteJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	res, err := helper.LoginRequest(client, *resp)
	if err != nil {
		at.WriteJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	// save token to db
	at.WriteJSON(w, http.StatusOK, res)

}

func SaveTokenString(w http.ResponseWriter, req *http.Request) {
	jar, _ := cookiejar.New(nil)

	// Create a new HTTP client with the cookie jar
	client := &http.Client{
		Jar: jar,
	}

	login := at.GetLoginFromHeader(req)
	if login == "" {
		at.WriteJSON(w, http.StatusForbidden, "No valid login found")
		return
	}

	token, err := helper.GetRefreshToken(client, login)
	if err != nil {
		if errors.Is(err, errors.New("no token found")) {
			at.WriteJSON(w, http.StatusForbidden, "token is invalid")
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		at.WriteJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Simpan token ke database
	tokenData := model.TokenData{
		UserID:    login, // Asumsikan login adalah userID
		Token:     token,
		UpdatedAt: time.Now(),
	}

	_, err = atdb.InsertOneDoc(config.Mongoconn, "tokens", tokenData)
	if err != nil {
		var respn model.Response
		respn.Status = "Gagal Insert Database"
		respn.Response = err.Error()
		at.WriteJSON(w, http.StatusNotModified, respn)
		return
	}

	result := &model.ResponseAct{
		Login:     true,
		SxSession: token,
	}

	at.WriteJSON(w, http.StatusOK, result)
}
