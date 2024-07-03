package domyApi

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/cookiejar"
	"strings"
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

	// Cek apakah user_id sudah ada di database
	existingTokenData, err := atdb.GetOneDoc[model.TokenData](config.Mongoconn, "tokens", primitive.M{"user_id": reqLogin.Email})
	if err != nil {
		// Jika terjadi kesalahan selain tidak menemukan dokumen, kembalikan error
		var respn model.Response
		respn.Status = "Gagal memeriksa database"
		respn.Response = err.Error()
		at.WriteJSON(w, http.StatusInternalServerError, respn)
		return
	}

	if existingTokenData.UserID == "" {
		// Jika user_id tidak ditemukan, insert data baru
		tokenData := model.TokenData{
			UserID:    reqLogin.Email,
			Token:     res.Session,
			Role:      reqLogin.Role,
			Password:  reqLogin.Password,
			UpdatedAt: time.Now(),
		}
		_, insertErr := atdb.InsertOneDoc(config.Mongoconn, "tokens", tokenData)
		if insertErr != nil {
			var respn model.Response
			respn.Status = "Gagal Insert Database"
			respn.Response = insertErr.Error()
			at.WriteJSON(w, http.StatusNotModified, respn)
			return
		}
		at.WriteJSON(w, http.StatusOK, tokenData)
	} else {
		// Jika user_id ditemukan, perbarui token yang ada
		update := bson.M{
			"$set": bson.M{
				"token":      res.Session,
				"updated_at": time.Now(),
			},
		}
		_, updateErr := atdb.UpdateOneDoc(config.Mongoconn, "tokens", primitive.M{"user_id": reqLogin.Email}, update)
		if updateErr != nil {
			var respn model.Response
			respn.Status = "Gagal Update Database"
			respn.Response = updateErr.Error()
			at.WriteJSON(w, http.StatusInternalServerError, respn)
			return
		}
		at.WriteJSON(w, http.StatusOK, existingTokenData)
	}

	// Ambil dan simpan data mahasiswa atau dosen
	if reqLogin.Role == "mahasiswa" {
		err = saveMahasiswaData(client, res.Session, reqLogin.Email)
		if err != nil {
			at.WriteJSON(w, http.StatusInternalServerError, err.Error())
			return
		}
	} else if reqLogin.Role == "dosen" {
		err = saveDosenData(client, res.Session, reqLogin.Email)
		if err != nil {
			at.WriteJSON(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
}

func saveMahasiswaData(_ *http.Client, token, email string) error {
	urlTarget := "https://siakad.ulbi.ac.id/siakad/data_mahasiswa"

	cookies := map[string]string{
		"SIAKAD_CLOUD_ACCESS": token,
	}

	doc, err := helper.GetData(urlTarget, cookies, nil)
	if err != nil {
		return err
	}

	nim := strings.TrimSpace(doc.Find("#block-nim .input-nim").Text())
	nama := strings.TrimSpace(doc.Find("#block-nama .input-nama").Text())
	programStudi := strings.TrimSpace(doc.Find("#block-idunit .input-idunit").Text())
	noHp := strings.TrimSpace(doc.Find("#block-hp .input-hp").Text())

	mahasiswa := model.Mahasiswa{
		Email:        email,
		NIM:          nim,
		Nama:         nama,
		ProgramStudi: programStudi,
		NomorHp:      noHp,
	}

	_, err = atdb.InsertOneDoc(config.Mongoconn, "mahasiswa", mahasiswa)
	return err
}

func saveDosenData(_ *http.Client, token, email string) error {
	urlTarget := "https://siakad.ulbi.ac.id/siakad/data_pegawai"

	cookies := map[string]string{
		"SIAKAD_CLOUD_ACCESS": token,
	}

	doc, err := helper.GetData(urlTarget, cookies, nil)
	if err != nil {
		return err
	}

	nip := strings.TrimSpace(doc.Find("#block-nip .input-nip").Text())
	nidn := strings.TrimSpace(doc.Find("#block-nidn .input-nidn").Text())
	nama := strings.TrimSpace(doc.Find("#block-nama .input-nama").Text())
	noHp := strings.TrimSpace(doc.Find("#block-nohp .input-nohp").Text())

	dosen := model.Dosen{
		Email: email,
		NIP:   nip,
		NIDN:  nidn,
		Nama:  nama,
		NoHp:  noHp,
	}

	_, err = atdb.InsertOneDoc(config.Mongoconn, "dosen", dosen)
	return err
}

func RefreshTokenDosen(w http.ResponseWriter, req *http.Request) {
	jar, _ := cookiejar.New(nil)

	// Create a new HTTP client with the cookie jar
	client := &http.Client{
		Jar: jar,
	}

	// Ambil login dari header
	login := req.Header.Get("login")
	if login == "" {
		http.Error(w, "No valid login found", http.StatusForbidden)
		return
	}

	// Ambil token dari database berdasarkan user_id
	tokenData, err := atdb.GetOneDoc[model.TokenData](config.Mongoconn, "tokens", primitive.M{"user_id": login})
	if err != nil {
		http.Error(w, "Token not found for user", http.StatusNotFound)
		return
	}

	// Gunakan GetRefreshTokenDosen untuk memperbarui token
	newToken, err := helper.GetRefreshTokenDosen(client, tokenData.Token)
	if err != nil {
		if errors.Is(err, errors.New("no token found")) {
			http.Error(w, "token is invalid", http.StatusForbidden)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Memperbarui token di database
	update := bson.M{
		"$set": bson.M{
			"token":      newToken,
			"updated_at": time.Now(),
		},
	}
	_, err = atdb.UpdateOneDoc(config.Mongoconn, "tokens", primitive.M{"user_id": login}, update)
	if err != nil {
		http.Error(w, "Failed to update database", http.StatusInternalServerError)
		return
	}

	result := &model.ResponseAct{
		Login:     true,
		SxSession: newToken,
	}

	at.WriteJSON(w, http.StatusOK, result)
}

func RefreshTokenMahasiswa(w http.ResponseWriter, req *http.Request) {
	jar, _ := cookiejar.New(nil)

	// Create a new HTTP client with the cookie jar
	client := &http.Client{
		Jar: jar,
	}

	// Ambil login dari header
	login := req.Header.Get("login")
	if login == "" {
		http.Error(w, "No valid login found", http.StatusForbidden)
		return
	}

	// Ambil token dari database berdasarkan user_id
	tokenData, err := atdb.GetOneDoc[model.TokenData](config.Mongoconn, "tokens", primitive.M{"user_id": login})
	if err != nil {
		http.Error(w, "Token not found for user", http.StatusNotFound)
		return
	}

	// Gunakan GetRefreshTokenMahasiswa untuk memperbarui token
	newToken, err := helper.GetRefreshTokenMahasiswa(client, tokenData.Token)
	if err != nil {
		if errors.Is(err, errors.New("no token found")) {
			http.Error(w, "token is invalid", http.StatusForbidden)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Memperbarui token di database
	update := bson.M{
		"$set": bson.M{
			"token":      newToken,
			"updated_at": time.Now(),
		},
	}
	_, err = atdb.UpdateOneDoc(config.Mongoconn, "tokens", primitive.M{"user_id": login}, update)
	if err != nil {
		http.Error(w, "Failed to update database", http.StatusInternalServerError)
		return
	}

	result := &model.ResponseAct{
		Login:     true,
		SxSession: newToken,
	}

	at.WriteJSON(w, http.StatusOK, result)
}
