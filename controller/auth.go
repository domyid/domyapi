package domyApi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"time"

	config "github.com/domyid/domyapi/config"
	at "github.com/domyid/domyapi/helper/at"
	helper "github.com/domyid/domyapi/helper/atapi"
	atdb "github.com/domyid/domyapi/helper/atdb"
	model "github.com/domyid/domyapi/model"
	"go.mongodb.org/mongo-driver/bson"
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
	jar, _ := cookiejar.New(nil)

	// Create a new HTTP client with the cookie jar
	client := &http.Client{
		Jar: jar,
	}

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

	// Menggunakan GetRefreshToken untuk memperbarui token
	newToken, err := helper.GetRefreshToken(client, tokenData.Token)
	if err != nil {
		if errors.Is(err, errors.New("no token found")) {
			at.WriteJSON(w, http.StatusForbidden, "token is invalid")
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		at.WriteJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Memperbarui token di database
	update := bson.M{
		"$set": bson.M{
			"token":      newToken,
			"updated_at": time.Now(),
		},
	}
	_, err = atdb.UpdateDoc(config.Mongoconn, "tokens", primitive.M{"user_id": login}, update)
	if err != nil {
		fmt.Println("Error Updating Token:", err)
		var respn model.Response
		respn.Status = "Gagal Update Database"
		respn.Response = err.Error()
		at.WriteJSON(w, http.StatusNotModified, respn)
		return
	}

	// Mengembalikan token yang diperbarui
	result := &model.ResponseAct{
		Login:     true,
		SxSession: newToken,
	}

	at.WriteJSON(w, http.StatusOK, result)
}
