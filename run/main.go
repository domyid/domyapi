package main

import (
	"bytes"
	"fmt"
	"os"

	pdf "github.com/domyid/domyapi/helper/pdf"
	model "github.com/domyid/domyapi/model"
)

func main() {
	// Data BAP contoh
	data := model.BAP{
		Kode:       "IF123",
		MataKuliah: "Pemrograman Go",
		Kelas:      "TI-1A",
		SMT:        "6",
		SKS:        "3",
		RiwayatMengajar: []model.RiwayatMengajar{
			{
				Pertemuan:       "1",
				Tanggal:         "2023-09-01",
				Jam:             "08:00-10:00",
				RencanaMateri:   "Pendahuluan",
				RealisasiMateri: "Pendahuluan",
			},
			{
				Pertemuan:       "2",
				Tanggal:         "2023-09-08",
				Jam:             "08:00-10:00",
				RencanaMateri:   "Dasar Pemrograman",
				RealisasiMateri: "Dasar Pemrograman",
			},
		},
		AbsensiKelas: []model.Absensi{
			{
				No:         "1",
				NIM:        "1214001",
				Nama:       "Ahmad",
				Pertemuan:  "14",
				Alfa:       "0",
				Hadir:      "14",
				Ijin:       "0",
				Sakit:      "0",
				Presentase: "100%",
			},
			{
				No:         "2",
				NIM:        "1214002",
				Nama:       "Budi",
				Pertemuan:  "14",
				Alfa:       "1",
				Hadir:      "13",
				Ijin:       "0",
				Sakit:      "0",
				Presentase: "93%",
			},
		},
		ListNilai: []model.Nilai{
			{
				No:    "1",
				NIM:   "1214001",
				Nama:  "Ahmad",
				Hadir: "100",
				ATS:   "90",
				AAS:   "95",
				Nilai: "92.5",
				Grade: "A",
			},
			{
				No:    "2",
				NIM:   "1214002",
				Nama:  "Budi",
				Hadir: "93",
				ATS:   "80",
				AAS:   "85",
				Nilai: "86.0",
				Grade: "A",
			},
		},
	}

	// Generate PDF
	buf, fileName, err := pdf.GenerateBAPPDFwithoutsignature(data)
	if err != nil {
		fmt.Println("Error generating PDF:", err)
		return
	}

	// Save to file
	err = saveToFile(buf, fileName)
	if err != nil {
		fmt.Println("Error saving PDF:", err)
		return
	}

	fmt.Println("PDF saved as:", fileName)
}

// saveToFile saves the buffer to a file
func saveToFile(buf *bytes.Buffer, fileName string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = buf.WriteTo(file)
	return err
}

// func generateDocID(time string) string {
// 	hash := sha256.New()
// 	hash.Write([]byte(time))
// 	hashedBytes := hash.Sum(nil)
// 	return hex.EncodeToString(hashedBytes)
// }
