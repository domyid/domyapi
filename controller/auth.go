package domyApi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"time"

	config "github.com/domyid/domyapi/config"
	at "github.com/domyid/domyapi/helper/at"
	helper "github.com/domyid/domyapi/helper/atapi"
	atdb "github.com/domyid/domyapi/helper/atdb"
	model "github.com/domyid/domyapi/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	// Simpan session token ke database
	tokenData := model.TokenData{
		UserID:    reqLogin.Email, // Asumsikan email digunakan sebagai userID
		Token:     res.Session,
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

	at.WriteJSON(w, http.StatusOK, res)
}

func RefreshToken(w http.ResponseWriter, req *http.Request) {
	// Mengambil login dari header
	login := req.Header.Get("login")
	if login == "" {
		at.WriteJSON(w, http.StatusForbidden, "No valid login found")
		return
	}

	// Mengambil token dari database berdasarkan user_id
	tokenData, err := atdb.GetOneDoc[model.TokenData](config.Mongoconn, "tokens", primitive.M{"user_id": login})
	if err != nil {
		fmt.Println("Error Fetching Token:", err)
		at.WriteJSON(w, http.StatusNotFound, "Token not found for user")
		return
	}

	// Mengembalikan token dari database
	result := &model.ResponseAct{
		Login:     true,
		SxSession: tokenData.Token,
	}

	at.WriteJSON(w, http.StatusOK, result)
}
