package domyApi

import (
	"errors"
	"strings"

	model "github.com/domyid/domyapi/model"
)

// Fungsi untuk mengekstrak informasi mahasiswa dari dokumen HTML
func ExtractMahasiswaData(cookies map[string]string) (model.Mahasiswa, error) {
	urlTarget := "https://siakad.ulbi.ac.id/siakad/data_mahasiswa"

	// Mengirim permintaan untuk mengambil data mahasiswa
	doc, err := GetData(urlTarget, cookies, nil)
	if err != nil {
		return model.Mahasiswa{}, err
	}

	nim := strings.TrimSpace(doc.Find("#block-nim .input-nim").Text())
	if nim == "" {
		return model.Mahasiswa{}, errors.New("NIM not found")
	}

	nama := strings.TrimSpace(doc.Find("#block-nama .input-nama").Text())
	programStudi := strings.TrimSpace(doc.Find("#block-idunit .input-idunit").Text())
	noHp := strings.TrimSpace(doc.Find("#block-hp .input-hp").Text())
	nirm := strings.TrimSpace(doc.Find(".input-nirm").Text())
	periodeMasuk := strings.TrimSpace(doc.Find(".input-idperiode").Text())
	tahunKurikulum := strings.TrimSpace(doc.Find(".input-idkurikulum").Text())
	sistemKuliah := strings.TrimSpace(doc.Find(".input-idsistemkuliah").Text())
	kelas := strings.TrimSpace(doc.Find(".input-idkelasperkuliahan").Text())
	jenisPendaftaran := strings.TrimSpace(doc.Find(".input-istransfer").Text())
	jalurPendaftaran := strings.TrimSpace(doc.Find(".input-idjalurpendaftaran").Text())
	gelombang := strings.TrimSpace(doc.Find(".input-idgelombang").Text())
	tanggalMasuk := strings.TrimSpace(doc.Find(".input-tgldaftar").Text())
	kebutuhanKhusus := strings.TrimSpace(doc.Find(".input-isdisabilitas").Text())
	statusMahasiswa := strings.TrimSpace(doc.Find(".input-idstatusmhs").Text())

	mahasiswa := model.Mahasiswa{
		NIM:              nim,
		Nama:             nama,
		ProgramStudi:     programStudi,
		NomorHp:          noHp,
		NIRM:             nirm,
		PeriodeMasuk:     periodeMasuk,
		TahunKurikulum:   tahunKurikulum,
		SistemKuliah:     sistemKuliah,
		Kelas:            kelas,
		JenisPendaftaran: jenisPendaftaran,
		JalurPendaftaran: jalurPendaftaran,
		Gelombang:        gelombang,
		TanggalMasuk:     tanggalMasuk,
		KebutuhanKhusus:  kebutuhanKhusus,
		StatusMahasiswa:  statusMahasiswa,
	}

	return mahasiswa, nil
}

// Fungsi untuk mengekstrak informasi dosen dari dokumen HTML
func ExtractDosenData(cookies map[string]string) (model.Dosen, error) {
	urlTarget := "https://siakad.ulbi.ac.id/siakad/data_pegawai"

	// Mengirim permintaan untuk mengambil data dosen
	doc, err := GetData(urlTarget, cookies, nil)
	if err != nil {
		return model.Dosen{}, err
	}

	// Ekstrak informasi dosen dan hapus spasi berlebih
	nip := strings.TrimSpace(doc.Find(".input-nip").Text())
	nidn := strings.TrimSpace(doc.Find(".input-nidn").Text())
	nama := strings.TrimSpace(doc.Find(".input-nama").Text())
	gelarDepan := strings.TrimSpace(doc.Find(".input-gelardepan").Text())
	gelarBelakang := strings.TrimSpace(doc.Find(".input-gelarbelakang").Text())
	jenisKelamin := strings.TrimSpace(doc.Find(".input-jeniskelamin").Text())
	tempatLahir := strings.TrimSpace(doc.Find(".input-tmplahir").Text())
	tanggalLahir := strings.TrimSpace(doc.Find(".input-tgllahir").Text())
	agama := strings.TrimSpace(doc.Find(".input-idagama").Text())
	noHp := strings.TrimSpace(doc.Find(".input-nohp").Text())
	emailKampus := strings.TrimSpace(doc.Find(".input-emailkampus").Text())
	emailPribadi := strings.TrimSpace(doc.Find(".input-email").Text())

	// Ekstrak angka unik dari href
	href, exists := doc.Find(".profile-nav li.active a").Attr("href")
	var dataid string
	if exists {
		parts := strings.Split(href, "/")
		dataid = parts[len(parts)-1]
	}

	// Buat instance Dosen
	dosen := model.Dosen{
		NIP:           nip,
		NIDN:          nidn,
		Nama:          nama,
		GelarDepan:    gelarDepan,
		GelarBelakang: gelarBelakang,
		JenisKelamin:  jenisKelamin,
		TempatLahir:   tempatLahir,
		TanggalLahir:  tanggalLahir,
		Agama:         agama,
		NoHp:          noHp,
		EmailKampus:   emailKampus,
		EmailPribadi:  emailPribadi,
		DataId:        dataid,
	}

	return dosen, nil
}
