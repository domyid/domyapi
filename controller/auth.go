package domyApi

import (
	"encoding/json"
	"fmt"
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

func saveMahasiswaData(_ *http.Client, token, email string) (string, error) {
	cookies := map[string]string{
		"SIAKAD_CLOUD_ACCESS": token,
	}

	mahasiswa, err := helper.ExtractMahasiswaData(cookies)
	if err != nil {
		return "", err
	}

	mahasiswa.Email = email

	// Ubah nomor HP yang diawali dengan angka 0 menjadi 62
	if strings.HasPrefix(mahasiswa.NomorHp, "0") {
		mahasiswa.NomorHp = "62" + mahasiswa.NomorHp[1:]
	}

	// Cek apakah data mahasiswa sudah ada berdasarkan email
	existingMahasiswa, err := atdb.GetOneDoc[model.Mahasiswa](config.Mongoconn, "mahasiswa", primitive.M{"email": email})
	if err == nil && existingMahasiswa.Email != "" {
		// Data mahasiswa sudah ada, tidak perlu disimpan lagi
		return mahasiswa.NomorHp, nil
	}

	_, err = atdb.InsertOneDoc(config.Mongoconn, "mahasiswa", mahasiswa)
	return mahasiswa.NomorHp, err
}

func saveDosenData(_ *http.Client, token, email string) (string, error) {
	cookies := map[string]string{
		"SIAKAD_CLOUD_ACCESS": token,
	}

	dosen, err := helper.ExtractDosenData(cookies)
	if err != nil {
		return "", err
	}

	dosen.Email = email

	// Ubah nomor HP yang diawali dengan angka 0 menjadi 62
	if strings.HasPrefix(dosen.NoHp, "0") {
		dosen.NoHp = "62" + dosen.NoHp[1:]
	}

	// Cek apakah data dosen sudah ada berdasarkan email
	existingDosen, err := atdb.GetOneDoc[model.Dosen](config.Mongoconn, "dosen", primitive.M{"email": email})
	if err == nil && existingDosen.Email != "" {
		// Data dosen sudah ada, tidak perlu disimpan lagi
		return dosen.NoHp, nil
	}

	_, err = atdb.InsertOneDoc(config.Mongoconn, "dosen", dosen)
	return dosen.NoHp, err
}

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

	var noHp string

	// Ambil dan simpan data mahasiswa atau dosen
	if reqLogin.Role == "mhs" {
		noHp, err = saveMahasiswaData(client, res.Session, reqLogin.Email)
		if err != nil {
			at.WriteJSON(w, http.StatusInternalServerError, err.Error())
			return
		}
	} else if reqLogin.Role == "dosen" {
		noHp, err = saveDosenData(client, res.Session, reqLogin.Email)
		if err != nil {
			at.WriteJSON(w, http.StatusInternalServerError, err.Error())
			return
		}

		// Insert or update ApprovalBAP collection
		dosen, err := atdb.GetOneDoc[model.Dosen](config.Mongoconn, "dosen", primitive.M{"email": reqLogin.Email})
		if err != nil {
			at.WriteJSON(w, http.StatusInternalServerError, "Failed to fetch dosen data")
			return
		}

		approvalBAP := model.ApprovalBAP{
			Status:     false,
			DataID:     dosen.DataID,
			EmailDosen: reqLogin.Email,
		}

		_, err = atdb.InsertOneDoc(config.Mongoconn, "approvalbap", approvalBAP)
		if err != nil {
			at.WriteJSON(w, http.StatusInternalServerError, "Failed to insert approval BAP data")
			return
		}
	}

	// Cek apakah user_id sudah ada di database
	existingTokenData, err := atdb.GetOneDoc[model.TokenData](config.Mongoconn, "tokens", primitive.M{"user_id": reqLogin.Email})
	if err != nil {
		// Jika user_id tidak ditemukan, insert data baru
		tokenData := model.TokenData{
			UserID:    reqLogin.Email,
			Token:     res.Session,
			Role:      reqLogin.Role,
			Password:  reqLogin.Password,
			NoHp:      noHp,
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
				"nohp":       noHp,
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
}

// Refresh tokens function
func RefreshTokens(w http.ResponseWriter, req *http.Request) {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}

	tokens, err := atdb.GetAllDoc[[]model.TokenData](config.Mongoconn, "tokens", bson.M{})
	if err != nil {
		http.Error(w, "Failed to fetch tokens from database", http.StatusInternalServerError)
		return
	}

	for _, tokenData := range tokens {
		var newToken string
		if tokenData.Role == "dosen" {
			newToken, err = helper.GetRefreshTokenDosen(client, tokenData.Token)
		} else if tokenData.Role == "mhs" {
			newToken, err = helper.GetRefreshTokenMahasiswa(client, tokenData.Token)
		} else {
			continue
		}

		if err != nil {
			if strings.Contains(err.Error(), "no token found") {
				err := helper.Logout(client, tokenData)
				if err != nil {
					http.Error(w, fmt.Sprintf("Logout failed: %v", err), http.StatusInternalServerError)
					return
				}
				delErr := atdb.DeleteOneDoc(config.Mongoconn, "tokens", bson.M{"user_id": tokenData.UserID})
				if delErr != nil {
					http.Error(w, "Failed to delete invalid token", http.StatusInternalServerError)
					return
				}
				continue
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		update := bson.M{
			"$set": bson.M{
				"token":      newToken,
				"updated_at": time.Now(),
			},
		}
		_, err = atdb.UpdateOneDoc(config.Mongoconn, "tokens", bson.M{"user_id": tokenData.UserID}, update)
		if err != nil {
			http.Error(w, "Failed to update database", http.StatusInternalServerError)
			return
		}
	}

	result := &model.ResponseAct{
		Login:     true,
		SxSession: "All tokens refreshed successfully",
	}

	at.WriteJSON(w, http.StatusOK, result)
}
