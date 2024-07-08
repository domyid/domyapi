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

func FetchListTugasAkhirMahasiswa(noHp string) ([]model.TugasAkhirAllMahasiswa, error) {
	urlTarget := "https://siakad.ulbi.ac.id/siakad/list_ta"

	// Mengambil token dari database berdasarkan no_hp
	tokenData, err := atdb.GetOneDoc[model.TokenData](config.Mongoconn, "tokens", primitive.M{"nohp": noHp})
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

func GetDataIDFromTugasAkhir(noHp, nim string) (string, error) {
	// Fetch the list of Tugas Akhir for all students
	listTA, err := FetchListTugasAkhirMahasiswa(noHp)
	if err != nil {
		return "", err
	}

	// Search for the data ID based on NIM
	for _, ta := range listTA {
		if ta.NIM == nim {
			return ta.DataID, nil
		}
	}

	return "", fmt.Errorf("no valid data ID found for the given NIM")
}

// FetchListBimbingan retrieves the list of Bimbingan based on the given dataID and token.
func FetchListBimbingan(dataID, token string) ([]model.ListBimbingan, error) {
	// URL target untuk mendapatkan data list bimbingan
	urlTarget := fmt.Sprintf("https://siakad.ulbi.ac.id/siakad/list_bimbingan/%s", dataID)

	// Buat payload berisi informasi token
	payload := map[string]string{
		"SIAKAD_CLOUD_ACCESS": token,
	}

	// Mengirim permintaan untuk mengambil data list bimbingan
	doc, err := GetData(urlTarget, payload, nil)
	if err != nil {
		return nil, fmt.Errorf("error fetching data: %v", err)
	}

	// Ekstrak informasi dari respon
	var listBimbingan []model.ListBimbingan
	doc.Find("tbody tr").Each(func(i int, s *goquery.Selection) {
		no := strings.TrimSpace(s.Find("td").Eq(0).Text())
		tanggal := strings.TrimSpace(s.Find("td").Eq(1).Text())
		dosenPembimbing := strings.TrimSpace(s.Find("td").Eq(2).Text())
		topik := strings.TrimSpace(s.Find("td").Eq(3).Text())
		disetujui := s.Find("td").Eq(4).Text() != ""
		if s.Find("td").Eq(4).Find("i").HasClass("fa-check") {
			disetujui = true
		} else {
			disetujui = false
		}
		dataID, _ := s.Find("td").Eq(5).Find("button").Attr("data-id")

		listbimbingan := model.ListBimbingan{
			No:              no,
			Tanggal:         tanggal,
			DosenPembimbing: dosenPembimbing,
			Topik:           topik,
			Disetujui:       disetujui,
			DataID:          dataID,
		}
		listBimbingan = append(listBimbingan, listbimbingan)
	})

	return listBimbingan, nil
}
