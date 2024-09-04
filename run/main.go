package main

import (
	"bytes"
	"fmt"
	"log"
	"os"

	pdf "github.com/domyid/domyapi/helper/pdf"
	model "github.com/domyid/domyapi/model"
)

func main() {
	// Data BAP contoh
	// data := model.BAP{
	// 	Kode:       "IF123",
	// 	MataKuliah: "Pemrograman Go",
	// 	Kelas:      "TI-1A",
	// 	SMT:        "6",
	// 	SKS:        "3",
	// 	RiwayatMengajar: []model.RiwayatMengajar{
	// 		{Pertemuan: "1", Tanggal: "2023-09-01", Jam: "08:00-10:00", RencanaMateri: "Pendahuluan", RealisasiMateri: "Pendahuluan"},
	// 		{Pertemuan: "2", Tanggal: "2023-09-08", Jam: "08:00-10:00", RencanaMateri: "Dasar Pemrograman", RealisasiMateri: "Dasar Pemrograman"},
	// 		{Pertemuan: "3", Tanggal: "2023-09-15", Jam: "08:00-10:00", RencanaMateri: "Struktur Data", RealisasiMateri: "Struktur Data"},
	// 		{Pertemuan: "4", Tanggal: "2023-09-22", Jam: "08:00-10:00", RencanaMateri: "Algoritma", RealisasiMateri: "Algoritma"},
	// 		{Pertemuan: "5", Tanggal: "2023-09-29", Jam: "08:00-10:00", RencanaMateri: "Pemrograman Lanjutan", RealisasiMateri: "Pemrograman Lanjutan"},
	// 		{Pertemuan: "6", Tanggal: "2023-10-06", Jam: "08:00-10:00", RencanaMateri: "Basis Data", RealisasiMateri: "Basis Data"},
	// 		{Pertemuan: "7", Tanggal: "2023-10-13", Jam: "08:00-10:00", RencanaMateri: "Jaringan Komputer", RealisasiMateri: "Jaringan Komputer"},
	// 		{Pertemuan: "8", Tanggal: "2023-10-20", Jam: "08:00-10:00", RencanaMateri: "Keamanan Data", RealisasiMateri: "Keamanan Data"},
	// 		{Pertemuan: "9", Tanggal: "2023-10-27", Jam: "08:00-10:00", RencanaMateri: "Pemrograman Web", RealisasiMateri: "Pemrograman Web"},
	// 		{Pertemuan: "10", Tanggal: "2023-11-03", Jam: "08:00-10:00", RencanaMateri: "Pemrograman Mobile", RealisasiMateri: "Pemrograman Mobile"},
	// 		{Pertemuan: "11", Tanggal: "2023-11-10", Jam: "08:00-10:00", RencanaMateri: "Pemrograman GUI", RealisasiMateri: "Pemrograman GUI"},
	// 		{Pertemuan: "12", Tanggal: "2023-11-17", Jam: "08:00-10:00", RencanaMateri: "Pemrograman Jaringan", RealisasiMateri: "Pemrograman Jaringan"},
	// 		{Pertemuan: "13", Tanggal: "2023-11-24", Jam: "08:00-10:00", RencanaMateri: "Pemrograman Paralel", RealisasiMateri: "Pemrograman Paralel"},
	// 		{Pertemuan: "14", Tanggal: "2023-12-01", Jam: "08:00-10:00", RencanaMateri: "Pemrograman Asynchronous", RealisasiMateri: "Pemrograman Asynchronous"},
	// 		{Pertemuan: "15", Tanggal: "2023-12-08", Jam: "08:00-10:00", RencanaMateri: "Pemrograman Fungsi", RealisasiMateri: "Pemrograman Fungsi"},
	// 		{Pertemuan: "16", Tanggal: "2023-12-15", Jam: "08:00-10:00", RencanaMateri: "Pemrograman Objek", RealisasiMateri: "Pemrograman Objek"},
	// 	},
	// 	AbsensiKelas: []model.Absensi{
	// 		{No: "1", NIM: "1214001", Nama: "Ahmad", Pertemuan: "14", Alfa: "0", Hadir: "14", Ijin: "0", Sakit: "0", Presentase: "100%"},
	// 		{No: "2", NIM: "1214002", Nama: "Budi", Pertemuan: "14", Alfa: "1", Hadir: "13", Ijin: "0", Sakit: "0", Presentase: "93%"},
	// 		{No: "3", NIM: "1214003", Nama: "Cindy", Pertemuan: "14", Alfa: "0", Hadir: "14", Ijin: "0", Sakit: "0", Presentase: "100%"},
	// 		{No: "4", NIM: "1214004", Nama: "Dewi", Pertemuan: "14", Alfa: "2", Hadir: "12", Ijin: "0", Sakit: "0", Presentase: "86%"},
	// 		{No: "5", NIM: "1214005", Nama: "Eko", Pertemuan: "14", Alfa: "0", Hadir: "14", Ijin: "0", Sakit: "0", Presentase: "100%"},
	// 		{No: "6", NIM: "1214006", Nama: "Fajar", Pertemuan: "14", Alfa: "1", Hadir: "13", Ijin: "0", Sakit: "0", Presentase: "93%"},
	// 		{No: "7", NIM: "1214007", Nama: "Gina", Pertemuan: "14", Alfa: "0", Hadir: "14", Ijin: "0", Sakit: "0", Presentase: "100%"},
	// 		{No: "8", NIM: "1214008", Nama: "Hadi", Pertemuan: "14", Alfa: "2", Hadir: "12", Ijin: "0", Sakit: "0", Presentase: "86%"},
	// 		{No: "9", NIM: "1214009", Nama: "Intan", Pertemuan: "14", Alfa: "0", Hadir: "14", Ijin: "0", Sakit: "0", Presentase: "100%"},
	// 		{No: "10", NIM: "1214010", Nama: "Joko", Pertemuan: "14", Alfa: "1", Hadir: "13", Ijin: "0", Sakit: "0", Presentase: "93%"},
	// 		{No: "11", NIM: "1214011", Nama: "Kiki", Pertemuan: "14", Alfa: "0", Hadir: "14", Ijin: "0", Sakit: "0", Presentase: "100%"},
	// 		{No: "12", NIM: "1214012", Nama: "Lina", Pertemuan: "14", Alfa: "2", Hadir: "12", Ijin: "0", Sakit: "0", Presentase: "86%"},
	// 		{No: "13", NIM: "1214013", Nama: "Maya", Pertemuan: "14", Alfa: "0", Hadir: "14", Ijin: "0", Sakit: "0", Presentase: "100%"},
	// 		{No: "14", NIM: "1214014", Nama: "Nina", Pertemuan: "14", Alfa: "1", Hadir: "13", Ijin: "0", Sakit: "0", Presentase: "93%"},
	// 		{No: "15", NIM: "1214015", Nama: "Omar", Pertemuan: "14", Alfa: "0", Hadir: "14", Ijin: "0", Sakit: "0", Presentase: "100%"},
	// 		{No: "16", NIM: "1214016", Nama: "Putu", Pertemuan: "14", Alfa: "2", Hadir: "12", Ijin: "0", Sakit: "0", Presentase: "86%"},
	// 		{No: "17", NIM: "1214017", Nama: "Qori", Pertemuan: "14", Alfa: "0", Hadir: "14", Ijin: "0", Sakit: "0", Presentase: "100%"},
	// 		{No: "18", NIM: "1214018", Nama: "Rama", Pertemuan: "14", Alfa: "1", Hadir: "13", Ijin: "0", Sakit: "0", Presentase: "93%"},
	// 		{No: "19", NIM: "1214019", Nama: "Susi", Pertemuan: "14", Alfa: "0", Hadir: "14", Ijin: "0", Sakit: "0", Presentase: "100%"},
	// 		{No: "20", NIM: "1214020", Nama: "Tina", Pertemuan: "14", Alfa: "2", Hadir: "12", Ijin: "0", Sakit: "0", Presentase: "86%"},
	// 	},
	// 	ListNilai: []model.Nilai{
	// 		{No: "1", NIM: "1214001", Nama: "Ahmad", Hadir: "100", ATS: "90", AAS: "95", Nilai: "92.5", Grade: "A"},
	// 		{No: "2", NIM: "1214002", Nama: "Budi", Hadir: "93", ATS: "80", AAS: "85", Nilai: "86.0", Grade: "A"},
	// 		{No: "3", NIM: "1214003", Nama: "Cindy", Hadir: "100", ATS: "85", AAS: "90", Nilai: "90.5", Grade: "A"},
	// 		{No: "4", NIM: "1214004", Nama: "Dewi", Hadir: "86", ATS: "78", AAS: "82", Nilai: "82.0", Grade: "B"},
	// 		{No: "5", NIM: "1214005", Nama: "Eko", Hadir: "100", ATS: "90", AAS: "95", Nilai: "92.5", Grade: "A"},
	// 		{No: "6", NIM: "1214006", Nama: "Fajar", Hadir: "93", ATS: "80", AAS: "85", Nilai: "86.0", Grade: "A"},
	// 		{No: "7", NIM: "1214007", Nama: "Gina", Hadir: "100", ATS: "85", AAS: "90", Nilai: "90.5", Grade: "A"},
	// 		{No: "8", NIM: "1214008", Nama: "Hadi", Hadir: "86", ATS: "78", AAS: "82", Nilai: "82.0", Grade: "B"},
	// 		{No: "9", NIM: "1214009", Nama: "Intan", Hadir: "100", ATS: "90", AAS: "95", Nilai: "92.5", Grade: "A"},
	// 		{No: "10", NIM: "1214010", Nama: "Joko", Hadir: "93", ATS: "80", AAS: "85", Nilai: "86.0", Grade: "A"},
	// 		{No: "11", NIM: "1214011", Nama: "Kiki", Hadir: "100", ATS: "85", AAS: "90", Nilai: "90.5", Grade: "A"},
	// 		{No: "12", NIM: "1214012", Nama: "Lina", Hadir: "86", ATS: "78", AAS: "82", Nilai: "82.0", Grade: "B"},
	// 		{No: "13", NIM: "1214013", Nama: "Maya", Hadir: "100", ATS: "90", AAS: "95", Nilai: "92.5", Grade: "A"},
	// 		{No: "14", NIM: "1214014", Nama: "Nina", Hadir: "93", ATS: "80", AAS: "85", Nilai: "86.0", Grade: "A"},
	// 		{No: "15", NIM: "1214015", Nama: "Omar", Hadir: "100", ATS: "85", AAS: "90", Nilai: "90.5", Grade: "A"},
	// 		{No: "16", NIM: "1214016", Nama: "Putu", Hadir: "86", ATS: "78", AAS: "82", Nilai: "82.0", Grade: "B"},
	// 		{No: "17", NIM: "1214017", Nama: "Qori", Hadir: "100", ATS: "90", AAS: "95", Nilai: "92.5", Grade: "A"},
	// 		{No: "18", NIM: "1214018", Nama: "Rama", Hadir: "93", ATS: "80", AAS: "85", Nilai: "86.0", Grade: "A"},
	// 		{No: "19", NIM: "1214019", Nama: "Susi", Hadir: "100", ATS: "85", AAS: "90", Nilai: "90.5", Grade: "A"},
	// 		{No: "20", NIM: "1214020", Nama: "Tina", Hadir: "86", ATS: "78", AAS: "82", Nilai: "82.0", Grade: "B"},
	// 	},
	// 	ProgramStudi: "D4 Teknik Informatika",
	// }

	// var buf *bytes.Buffer
	// var fileName string
	// var err error

	// if data.ProgramStudi == "D4 Teknik Informatika" {
	// 	// Dummy SignatureData
	// 	signature := model.SignatureData{
	// 		PenandaTangan:   "Roni Andarsyah",
	// 		DocName:         "BAP IF123 Pemrograman Go.pdf",
	// 		PemilikDocument: "Universitas Logistik dan Bisnis Internasional",
	// 	}

	// 	// Create QR code link
	// 	docID := generateDocID(time.Now().String())
	// 	token := pdf.CreateToken(docID, "https://mrt.ulbi.ac.id/token/create", signature)
	// 	qrCodeLink := pdf.GenerateLink(token)

	// 	// Generate PDF with signature
	// 	buf, fileName, err = pdf.GenerateBAPPDF(data, qrCodeLink)
	// 	if err != nil {
	// 		fmt.Println("Error generating PDF:", err)
	// 		return
	// 	}
	// } else {
	// 	// Generate PDF without signature
	// 	buf, fileName, err = pdf.GenerateBAPPDFwithoutsignature(data)
	// 	if err != nil {
	// 		fmt.Println("Error generating PDF:", err)
	// 		return
	// 	}
	// }
	// Create dummy data
	data := model.RekapBimbingan{
		JudulProposal:      "Pengembangan Aplikasi Mobile",
		SesiBahasan:        "Desain User Interface",
		NIM:                "123456789",
		Mahasiswa:          "John Doe",
		PembimbingProposal: "Dr. Jane Smith",
		Percakapan:         "Diskusi mengenai desain dan fitur aplikasi.",
	}

	// Generate the PDF
	buf, fileName, err := pdf.GenerateBKD(data)
	if err != nil {
		log.Fatalf("Error generating PDF: %v", err)
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
