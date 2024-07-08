package domyApi

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	config "github.com/domyid/domyapi/config"
	atdb "github.com/domyid/domyapi/helper/atdb"
	model "github.com/domyid/domyapi/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Fungsi untuk mengekstrak informasi dosen dari dokumen HTML dan mendapatkan dataid
func ExtractDataid(cookies map[string]string) (string, error) {
	urlTarget := "https://siakad.ulbi.ac.id/siakad/data_pegawai"

	// Mengirim permintaan untuk mengambil data dosen
	doc, err := GetData(urlTarget, cookies, nil)
	if err != nil {
		return "", err
	}

	// Ekstrak dataid dari elemen yang sesuai
	dataid, exists := doc.Find(".profile-nav li.active a").Attr("href")
	if !exists {
		return "", fmt.Errorf("dataid not found")
	}

	// Ambil angka unik dari href
	parts := strings.Split(dataid, "/")
	if len(parts) == 0 {
		return "", fmt.Errorf("invalid dataid format")
	}
	dataid = parts[len(parts)-1]

	return dataid, nil
}

// Fungsi untuk mendapatkan data jadwal mengajar
func FetchJadwalMengajar(noHp, periode string) ([]model.JadwalMengajar, error) {
	// Mengambil token dari database berdasarkan no_hp
	tokenData, err := atdb.GetOneDoc[model.TokenData](config.Mongoconn, "tokens", primitive.M{"nohp": noHp})
	if err != nil {
		return nil, fmt.Errorf("error fetching token: %v", err)
	}

	// Buat payload berisi informasi token
	payload := map[string]string{
		"SIAKAD_CLOUD_ACCESS": tokenData.Token,
	}

	// Ambil dataid dari halaman detail dosen
	dataid, err := ExtractDosenData(payload)
	if err != nil {
		return nil, fmt.Errorf("error extracting dataid: %v", err)
	}

	// URL target untuk mendapatkan data jadwal mengajar
	urlTarget := fmt.Sprintf("https://siakad.ulbi.ac.id/siakad/list_jadwalmengajar/%s", dataid)

	// Mengirim permintaan untuk mendapatkan data jadwal mengajar dengan metode POST
	formData := url.Values{}
	formData.Set("periode", periode)

	doc, err := GetDataPOST(urlTarget, payload, formData, nil)
	if err != nil {
		return nil, fmt.Errorf("error fetching data: %v", err)
	}

	// Ekstrak informasi dari respon
	var listJadwal []model.JadwalMengajar
	doc.Find("tbody tr").Each(func(i int, s *goquery.Selection) {
		no := strings.TrimSpace(s.Find("td").Eq(0).Text())
		kode := strings.TrimSpace(s.Find("td").Eq(1).Text())
		mataKuliah := strings.TrimSpace(s.Find("td").Eq(2).Text())
		sks := strings.TrimSpace(s.Find("td").Eq(3).Text())
		smt := strings.TrimSpace(s.Find("td").Eq(4).Text())
		kelas := strings.TrimSpace(s.Find("td").Eq(5).Text())
		programStudi := strings.TrimSpace(s.Find("td").Eq(6).Text())
		hari := strings.TrimSpace(s.Find("td").Eq(7).Text())
		waktu := strings.TrimSpace(s.Find("td").Eq(8).Text())
		ruang := strings.TrimSpace(s.Find("td").Eq(9).Text())

		jadwal := model.JadwalMengajar{
			No:           no,
			Kode:         kode,
			MataKuliah:   mataKuliah,
			SKS:          sks,
			Smt:          smt,
			Kelas:        kelas,
			ProgramStudi: programStudi,
			Hari:         hari,
			Waktu:        waktu,
			Ruang:        ruang,
		}
		listJadwal = append(listJadwal, jadwal)
	})

	return listJadwal, nil
}
