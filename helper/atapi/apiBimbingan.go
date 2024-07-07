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

func FetchListTugasAkhirAllMahasiswa(noHp string) ([]model.TugasAkhirAllMahasiswa, error) {
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

func FetchListTugasAkhirMahasiswa(noHp string) ([]model.TugasAkhirMahasiswa, error) {
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

func FetchListBimbingan(noHp, nim string) ([]model.ListBimbingan, error) {
	// Mengambil daftar tugas akhir mahasiswa
	listTA, err := FetchListTugasAkhirAllMahasiswa(noHp)
	if err != nil {
		return nil, fmt.Errorf("error Fetching List Tugas Akhir: %v", err)
	}

	// Cari dataID berdasarkan NIM
	var dataID string
	for _, ta := range listTA {
		if ta.NIM == nim {
			dataID = ta.DataID
			break
		}
	}

	if dataID == "" {
		return nil, fmt.Errorf("no valid data ID found for the provided NIM")
	}

	// URL target untuk list bimbingan
	urlTarget := fmt.Sprintf("https://siakad.ulbi.ac.id/siakad/list_bimbingan/%s", dataID)

	// Mengambil token dari database berdasarkan nohp
	tokenData, err := atdb.GetOneDoc[model.TokenData](config.Mongoconn, "tokens", primitive.M{"nohp": noHp})
	if err != nil {
		return nil, fmt.Errorf("error Fetching Token: %v", err)
	}

	// Buat payload berisi informasi token
	payload := map[string]string{
		"SIAKAD_CLOUD_ACCESS": tokenData.Token,
	}

	// Mengirim permintaan untuk mengambil data list bimbingan
	doc, err := GetData(urlTarget, payload, nil)
	if err != nil {
		return nil, fmt.Errorf("error Fetching Data: %v", err)
	}

	// Ekstrak informasi dari respon
	var listBimbingan []model.ListBimbingan
	doc.Find("tbody tr").Each(func(i int, s *goquery.Selection) {
		no := strings.TrimSpace(s.Find("td").Eq(0).Text())
		tanggal := strings.TrimSpace(s.Find("td").Eq(1).Text())
		dosenPembimbing := strings.TrimSpace(s.Find("td").Eq(2).Text())
		topik := strings.TrimSpace(s.Find("td").Eq(3).Text())
		disetujui := s.Find("td").Eq(4).Find("i").HasClass("fa-check")
		dataID, _ := s.Find("td").Eq(5).Find("button").Attr("data-id")

		bimbingan := model.ListBimbingan{
			No:              no,
			Tanggal:         tanggal,
			DosenPembimbing: dosenPembimbing,
			Topik:           topik,
			Disetujui:       disetujui,
			DataID:          dataID,
		}
		listBimbingan = append(listBimbingan, bimbingan)
	})

	return listBimbingan, nil
}
