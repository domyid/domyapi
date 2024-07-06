package domyApi

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	config "github.com/domyid/domyapi/config"
	at "github.com/domyid/domyapi/helper/at"
	api "github.com/domyid/domyapi/helper/atapi"
	atdb "github.com/domyid/domyapi/helper/atdb"
	model "github.com/domyid/domyapi/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetListTugasAkhir(respw http.ResponseWriter, req *http.Request) {
	urlTarget := "https://siakad.ulbi.ac.id/siakad/list_ta"

	// Mengambil user_id dari header
	userID := req.Header.Get("user_id")
	if userID == "" {
		http.Error(respw, "No valid user ID found", http.StatusForbidden)
		return
	}

	// Mengambil token dari database berdasarkan user_id
	tokenData, err := atdb.GetOneDoc[model.TokenData](config.Mongoconn, "tokens", primitive.M{"user_id": userID})
	if err != nil {
		fmt.Println("Error Fetching Token:", err)
		at.WriteJSON(respw, http.StatusNotFound, "Token not found for user")
		return
	}

	// Buat payload berisi informasi token
	payload := map[string]string{
		"SIAKAD_CLOUD_ACCESS": tokenData.Token,
	}

	// Mengirim permintaan untuk mengambil data list TA
	doc, err := api.GetData(urlTarget, payload, nil)
	if err != nil {
		at.WriteJSON(respw, http.StatusInternalServerError, err.Error())
		return
	}

	// Ekstrak informasi dari respon
	var listTA []model.TugasAkhir
	doc.Find("tbody tr").Each(func(i int, s *goquery.Selection) {
		nama := strings.TrimSpace(s.Find("td").Eq(0).Find("strong").Text())
		nim := strings.TrimSpace(s.Find("td").Eq(0).Contents().FilterFunction(func(_ int, selection *goquery.Selection) bool {
			return goquery.NodeName(selection) == "#text"
		}).Text())
		judul := strings.TrimSpace(s.Find("td").Eq(1).Text())
		pembimbing := strings.TrimSpace(s.Find("td").Eq(2).Text())
		tglMulai := strings.TrimSpace(s.Find("td").Eq(3).Text())
		status := strings.TrimSpace(s.Find("td").Eq(4).Text())
		dataID, _ := s.Find("td").Eq(5).Find(".btn-group .action-link").Attr("data-id")

		ta := model.TugasAkhir{
			Nama:         nama,
			NIM:          nim,
			Judul:        judul,
			Pembimbing:   pembimbing,
			TanggalMulai: tglMulai,
			Status:       status,
			DataID:       dataID,
		}
		listTA = append(listTA, ta)
	})

	// Kembalikan daftar TA sebagai respon JSON
	at.WriteJSON(respw, http.StatusOK, listTA)
}
