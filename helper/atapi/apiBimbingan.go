package domyApi

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	config "github.com/domyid/domyapi/config"
	atdb "github.com/domyid/domyapi/helper/atdb"
	model "github.com/domyid/domyapi/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetListTa(client *http.Client, token string) (string, error) {
	homeURL := "https://siakad.ulbi.ac.id/siakad/list_ta"

	req, err := http.NewRequest("GET", homeURL, nil)
	if err != nil {
		return "", nil
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Referer", "https://siakad.ulbi.ac.id/gate/menu")
	req.Header.Set("Cookie", fmt.Sprintf("SIAKAD_CLOUD_ACCESS=%s", token))

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	tokenStr := resp.Header.Get("Sx-Session")

	if tokenStr == "" {
		return "", fmt.Errorf("no token found")
	}

	return tokenStr, nil
}

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
