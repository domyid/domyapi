package domyApi

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
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

		disetujui := false
		if s.Find("td").Eq(4).Find("i").HasClass("fa-check") {
			disetujui = true
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

// Fungsi untuk mengambil data dari halaman edit
func GetDetailBimbingan(bimbinganID, token string) (model.DetailBimbingan, error) {
	// URL target untuk mendapatkan data dari halaman edit
	urlTarget := fmt.Sprintf("https://siakad.ulbi.ac.id/siakad/data_bimbingan/edit/%s", bimbinganID)

	// Buat payload berisi informasi token
	payload := map[string]string{
		"SIAKAD_CLOUD_ACCESS": token,
	}

	// Mengirim permintaan untuk mengambil data dari halaman edit
	doc, err := GetData(urlTarget, payload, nil)
	if err != nil {
		return model.DetailBimbingan{}, fmt.Errorf("error fetching data: %v", err)
	}

	// Ekstrak informasi dari respon
	detail := model.DetailBimbingan{
		BimbinganKe:    doc.Find("input[name='bimbinganke']").AttrOr("value", ""),
		NIP:            doc.Find("select[name='nip'] option[selected]").AttrOr("value", ""),
		TglBimbingan:   doc.Find("input[name='tglbimbingan']").AttrOr("value", ""),
		TopikBimbingan: doc.Find("input[name='topikbimbingan']").AttrOr("value", ""),
		Bahasan:        doc.Find("textarea[name='bahasan']").Text(),
		Link:           "",
		Lampiran:       "",
		Key:            "",
		Act:            "save",
	}

	return detail, nil
}

// Fungsi untuk mengirim data yang telah disetujui
func ApproveBimbingan(bimbinganID, token string, data model.DetailBimbingan) error {
	postURL := fmt.Sprintf("https://siakad.ulbi.ac.id/siakad/data_bimbingan/edit/%s", bimbinganID)

	form := url.Values{}
	form.Add("bimbinganke", data.BimbinganKe)
	form.Add("nip", data.NIP)
	form.Add("tglbimbingan", data.TglBimbingan)
	form.Add("topikbimbingan", data.TopikBimbingan)
	form.Add("bahasan", data.Bahasan)
	form.Add("disetujui", "1") // Set the approved flag
	form.Add("link[]", data.Link)
	form.Add("lampiran[]", data.Lampiran)
	form.Add("key", data.Key)
	form.Add("act", data.Act)

	req, err := http.NewRequest("POST", postURL, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Cookie", fmt.Sprintf("SIAKAD_CLOUD_ACCESS=%s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusSeeOther && resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}
