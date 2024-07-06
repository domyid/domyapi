package domyApi

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	config "github.com/domyid/domyapi/config"
	atdb "github.com/domyid/domyapi/helper/atdb"
	model "github.com/domyid/domyapi/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func FetchListTugasAkhirMahasiswa(userID string) ([]model.TugasAkhirMahasiswa, error) {
	urlTarget := "https://siakad.ulbi.ac.id/siakad/list_ta"

	// Mengambil token dari database berdasarkan user_id
	tokenData, err := atdb.GetOneDoc[model.TokenData](config.Mongoconn, "tokens", primitive.M{"user_id": userID})
	if err != nil {
		return nil, fmt.Errorf("error Fetching Token: %v", err)
	}

	// Buat payload berisi informasi token
	payload := map[string]string{
		"SIAKAD_CLOUD_ACCESS": tokenData.Token,
	}

	// Mengirim permintaan untuk mengambil data list TA
	doc, err := GetData(urlTarget, payload, nil)
	if err != nil {
		return nil, fmt.Errorf("error Fetching Data: %v", err)
	}

	// Ekstrak informasi dari respon
	var listTA []model.TugasAkhirMahasiswa
	doc.Find("tbody tr").Each(func(i int, s *goquery.Selection) {
		judul := strings.TrimSpace(s.Find("td").Eq(1).Text())
		pembimbing1 := strings.TrimSpace(s.Find("td").Eq(2).Find("ol li").Eq(0).Text())
		pembimbing2 := strings.TrimSpace(s.Find("td").Eq(2).Find("ol li").Eq(1).Text())
		tglMulai := strings.TrimSpace(s.Find("td").Eq(3).Text())
		status := strings.TrimSpace(s.Find("td").Eq(4).Find("h3").Text())
		dataID, _ := s.Find("td").Eq(5).Find(".btn-group .action-link").Attr("data-id")

		ta := model.TugasAkhirMahasiswa{
			Judul:        judul,
			Pembimbing1:  pembimbing1,
			Pembimbing2:  pembimbing2,
			TanggalMulai: tglMulai,
			Status:       status,
			DataID:       dataID,
		}
		listTA = append(listTA, ta)
	})

	return listTA, nil
}

func FetchListTugasAkhirAllMahasiswa(userID string) ([]model.TugasAkhirAllMahasiswa, error) {
	urlTarget := "https://siakad.ulbi.ac.id/siakad/list_ta"

	// Mengambil token dari database berdasarkan user_id
	tokenData, err := atdb.GetOneDoc[model.TokenData](config.Mongoconn, "tokens", primitive.M{"user_id": userID})
	if err != nil {
		return nil, fmt.Errorf("error Fetching Token: %v", err)
	}

	// Buat payload berisi informasi token
	payload := map[string]string{
		"SIAKAD_CLOUD_ACCESS": tokenData.Token,
	}

	// Mengirim permintaan untuk mengambil data list TA
	doc, err := GetData(urlTarget, payload, nil)
	if err != nil {
		return nil, fmt.Errorf("error Fetching Data: %v", err)
	}

	// Ekstrak informasi dari respon
	var listTA []model.TugasAkhirAllMahasiswa
	doc.Find("tbody tr").Each(func(i int, s *goquery.Selection) {
		nama := strings.TrimSpace(s.Find("td").Eq(0).Find("strong").Text())
		nim := strings.TrimSpace(s.Find("td").Eq(0).Contents().FilterFunction(func(_ int, selection *goquery.Selection) bool {
			return goquery.NodeName(selection) == "#text"
		}).Text())
		judul := strings.TrimSpace(s.Find("td").Eq(1).Text())
		pembimbing1 := strings.TrimSpace(s.Find("td").Eq(2).Find("ol li").Eq(0).Text())
		pembimbing2 := strings.TrimSpace(s.Find("td").Eq(2).Find("ol li").Eq(1).Text())
		tglMulai := strings.TrimSpace(s.Find("td").Eq(3).Text())
		status := strings.TrimSpace(s.Find("td").Eq(4).Find("h3").Text())
		dataID, _ := s.Find("td").Eq(5).Find(".btn-group .action-link").Attr("data-id")

		ta := model.TugasAkhirAllMahasiswa{
			Nama:         nama,
			NIM:          nim,
			Judul:        judul,
			Pembimbing1:  pembimbing1,
			Pembimbing2:  pembimbing2,
			TanggalMulai: tglMulai,
			Status:       status,
			DataID:       dataID,
		}
		listTA = append(listTA, ta)
	})

	return listTA, nil
}
