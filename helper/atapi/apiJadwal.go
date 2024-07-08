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

// FetchJadwalMengajar retrieves the teaching schedule based on noHp and periode.
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

	// URL target untuk mendapatkan data jadwal mengajar
	urlTarget := fmt.Sprintf("https://siakad.ulbi.ac.id/siakad/list_jadwalmengajar/%s", tokenData.UserID)

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

		// Ambil data ID dari elemen href
		href, exists := s.Find("td").Eq(10).Find("a").Attr("href")
		var dataID string
		if exists {
			dataID = strings.TrimPrefix(href, "/siakad/data_kelas/detail/")
		}

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
			DataID:       dataID,
		}
		listJadwal = append(listJadwal, jadwal)
	})

	return listJadwal, nil
}

func FetchListAbsensi(dataID, token string) ([]model.Absensi, error) {
	// URL target untuk mendapatkan data list nilai
	urlTarget := fmt.Sprintf("https://siakad.ulbi.ac.id/siakad/list_absensi/%s", dataID)

	// Buat payload berisi informasi token
	payload := map[string]string{
		"SIAKAD_CLOUD_ACCESS": token,
	}

	// Mengirim permintaan untuk mengambil data list nilai
	doc, err := GetData(urlTarget, payload, nil)
	if err != nil {
		return nil, fmt.Errorf("error fetching data: %v", err)
	}

	var listAbsensi []model.Absensi

	/// Select the specific table inside the div with class 'table-responsive'
	doc.Find(".table-responsive table.dataTable tbody tr").Each(func(i int, s *goquery.Selection) {
		pertemuan := strings.TrimSpace(s.Find("td.text-center").Eq(0).Text())
		tanggalJam := strings.TrimSpace(s.Find("td.text-center").Eq(1).Text())
		tanggalJamSplit := strings.Split(tanggalJam, "\n")
		tanggal := strings.TrimSpace(tanggalJamSplit[0])
		jam := ""
		if len(tanggalJamSplit) > 1 {
			jam = strings.TrimSpace(tanggalJamSplit[1])
		}
		materi := strings.TrimSpace(s.Find("td.word-wrap").Eq(0).Text())
		pengajar := strings.TrimSpace(s.Find("td.word-wrap").Eq(1).Text())
		ruang := strings.TrimSpace(s.Find("td.text-center").Eq(2).Text())
		hadir := strings.TrimSpace(s.Find("td.text-right").Eq(0).Text())
		persentase := strings.TrimSpace(s.Find("td.text-right").Eq(1).Text())

		absensi := model.Absensi{
			Pertemuan:  pertemuan,
			Tanggal:    tanggal,
			Jam:        jam,
			Materi:     materi,
			Pengajar:   pengajar,
			Ruang:      ruang,
			Hadir:      hadir,
			Persentase: persentase,
		}

		listAbsensi = append(listAbsensi, absensi)
	})

	return listAbsensi, nil
}

func FetchListNilai(dataID, token string) ([]model.ListNilai, error) {
	// URL target untuk mendapatkan data list nilai
	urlTarget := fmt.Sprintf("https://siakad.ulbi.ac.id/siakad/set_nilai/%s", dataID)

	// Buat payload berisi informasi token
	payload := map[string]string{
		"SIAKAD_CLOUD_ACCESS": token,
	}

	// Mengirim permintaan untuk mengambil data list nilai
	doc, err := GetData(urlTarget, payload, nil)
	if err != nil {
		return nil, fmt.Errorf("error fetching data: %v", err)
	}

	// Ekstrak informasi dari respon
	var listNilai []model.ListNilai
	doc.Find("tbody tr").Each(func(i int, s *goquery.Selection) {
		no := strings.TrimSpace(s.Find("td").Eq(0).Text())
		nim := strings.TrimSpace(s.Find("td").Eq(1).Text())
		nama := strings.TrimSpace(s.Find("td").Eq(2).Text())
		hadir := strings.TrimSpace(s.Find("td").Eq(3).Text())
		ats := strings.TrimSpace(s.Find("td").Eq(4).Text())
		aas := strings.TrimSpace(s.Find("td").Eq(5).Text())
		nilai := strings.TrimSpace(s.Find("td").Eq(6).Text())
		grade := strings.TrimSpace(s.Find("td").Eq(7).Text())
		lulus := strings.TrimSpace(s.Find("td").Eq(8).Text())
		keterangan := strings.TrimSpace(s.Find("td").Eq(9).Text())

		nilaiRecord := model.ListNilai{
			No:         no,
			NIM:        nim,
			Nama:       nama,
			Hadir:      hadir,
			ATS:        ats,
			AAS:        aas,
			Nilai:      nilai,
			Grade:      grade,
			Lulus:      lulus,
			Keterangan: keterangan,
		}
		listNilai = append(listNilai, nilaiRecord)
	})

	return listNilai, nil
}
