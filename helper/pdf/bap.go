package domyApi

import (
	"bytes"
	"fmt"
	"strconv"
	"time"

	model "github.com/domyid/domyapi/model"
	"github.com/jung-kurt/gofpdf"
)

// CreateHeaderBAP generates the header for the BAP PDF
const SourceURL = "https://siakad.ulbi.ac.id/siakad/rep_perkuliahan"
const InfoImageURL = "https://home.ulbi.ac.id/ulbi.png"

func CreateHeaderBAP(Text []string, x float64) *gofpdf.Fpdf {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Menambahkan timestamp di pojok kiri atas dengan font lebih kecil dan tipis
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		loc = time.FixedZone("WIB", 7*3600) // Default ke WIB jika timezone gagal dimuat
	}
	timestamp := time.Now().In(loc).Format("2006-01-02 15:04:05")
	pdf.SetFont("Times", "", 10) // Font normal, ukuran 8
	pdf.SetY(5)
	pdf.SetX(10)
	pdf.CellFormat(0, 10, timestamp, "", 0, "L", false, 0, "")

	// Menambahkan URL sumber di pojok kanan atas dengan font lebih kecil dan tipis
	pdf.SetY(5)
	pdf.SetX(-10 - pdf.GetStringWidth(SourceURL)) // Sesuaikan posisi X agar rata kanan
	pdf.CellFormat(0, 10, SourceURL, "", 0, "R", false, 0, "")

	// Menambahkan teks header di sebelah kanan gambar dengan sedikit spasi dan ke tengah
	pdf.SetFont("Times", "B", 12) // Kembali ke font bold ukuran 12 untuk header
	pdf.SetY(20)
	pdf.SetX(70) // Menambah nilai X agar lebih ke tengah dan ada spasi dengan gambar
	pdf.CellFormat(70, 10, Text[0], "0", 0, "L", false, 0, "")
	pdf.Ln(5)    // Menambah jarak antara baris header
	pdf.SetX(80) // Menambah nilai X agar lebih ke tengah dan ada spasi dengan gambar
	pdf.CellFormat(70, 10, Text[1], "0", 0, "L", false, 0, "")
	pdf.Ln(10)

	// Menggunakan ImageCustomize untuk menambahkan gambar di sebelah kiri teks header
	pdf = ImageCustomize(pdf, "./ulbi.png", InfoImageURL, 30, 20, 35, 12, 100, 100, 0.3)
	return pdf
}

func GenerateBAPPDFWithoutSignature(data model.BAP) (*bytes.Buffer, string, error) {
	Text := []string{
		"UNIVERSITAS LOGISTIK DAN BISNIS INTERNASIONAL",
		"Berita Acara Perkuliahan dan Absensi Perkuliahan",
	}

	width := []float64{60, 5, 70}
	color := []int{255, 255, 153}
	align := []string{"J", "C", "J"}
	yCoordinates := []float64{40, 45, 50}

	pdf := CreateHeaderBAP(Text, 90)

	// Header Information
	headerInfo := [][]string{
		{"Kode Matakuliah/Nama Matakuliah", ":", fmt.Sprintf("%s/%s", data.Kode, data.MataKuliah)},
		{"Kelas", ":", data.Kelas},
		{"Semester/SKS", ":", fmt.Sprintf("%s/%s SKS", data.SMT, data.SKS)},
	}

	pdf = SetTableContentCustomY(pdf, headerInfo, width, align, yCoordinates)

	// Add Riwayat Mengajar table
	pdf.Ln(5)
	pdf = SetMergedCell(pdf, "Tabel Log Aktivitas", "J", 150, color)
	headers := []string{"Pertemuan", "Tanggal", "Jam", "Rencana Materi", "Realisasi Materi"}
	widths := []float64{20, 30, 20, 40, 40}
	align = []string{"C", "C", "C", "C", "C"}
	pdf = SetHeaderTable(pdf, headers, widths, []int{135, 206, 235})
	for _, item := range data.RiwayatMengajar {
		row := []string{
			item.Pertemuan,
			item.Tanggal,
			item.Jam,
			truncateToThreeWords(item.RencanaMateri),
			truncateToThreeWords(item.RealisasiMateri),
		}
		pdf = SetTableContent(pdf, [][]string{row}, widths, align)
	}

	// Add Absensi Kelas table
	pdf.Ln(10)
	pdf = SetMergedCell(pdf, "Tabel Presensi", "J", 150, color)
	headers = []string{"No", "NIM", "Nama", "Pertemuan", "Alfa", "Hadir", "Ijin", "Sakit", "Presentase"}
	widths = []float64{10, 20, 40, 20, 10, 10, 10, 10, 20}
	align = []string{"C", "C", "L", "C", "C", "C", "C", "C", "C"}
	pdf = SetHeaderTable(pdf, headers, widths, []int{135, 206, 235})
	for _, item := range data.AbsensiKelas {
		row := []string{
			item.No,
			item.NIM,
			item.Nama,
			item.Pertemuan,
			item.Alfa,
			item.Hadir,
			item.Ijin,
			item.Sakit,
			item.Presentase,
		}
		pdf = SetTableContent(pdf, [][]string{row}, widths, align)
	}

	// Add List Nilai table
	pdf.Ln(10)
	pdf = SetMergedCell(pdf, "Tabel Nilai Akhir", "J", 150, color)
	headers = []string{"No", "NIM", "Nama", "Hadir", "ATS", "AAS", "Nilai", "Grade"}
	widths = []float64{10, 20, 40, 15, 15, 15, 15, 20}
	align = []string{"C", "C", "L", "C", "C", "C", "C", "C"}
	pdf = SetHeaderTable(pdf, headers, widths, []int{135, 206, 235})
	for _, item := range data.ListNilai {
		hadir, _ := strconv.ParseFloat(item.Hadir, 64)
		ats, _ := strconv.ParseFloat(item.ATS, 64)
		aas, _ := strconv.ParseFloat(item.AAS, 64)
		nilai, _ := strconv.ParseFloat(item.Nilai, 64)

		row := []string{
			item.No,
			item.NIM,
			item.Nama,
			fmt.Sprintf("%.2f", hadir),
			fmt.Sprintf("%.2f", ats),
			fmt.Sprintf("%.2f", aas),
			fmt.Sprintf("%.2f", nilai),
			item.Grade,
		}
		pdf = SetTableContent(pdf, [][]string{row}, widths, align)
	}

	// Menambahkan tempat tanda tangan tanpa QR code
	tanggalTerkini := "Bandung, " + getFormattedDate(time.Now())
	pdf.Ln(10)
	pdf.SetFont("Times", "", 10)
	pdf.SetX(-68) // Set posisi tanda tangan di sebelah kanan
	pdf.CellFormat(0, 5, tanggalTerkini, "", 1, "L", false, 0, "")
	pdf.SetX(-75)
	pdf.CellFormat(0, 5, "Ketua Prodi D4 Teknik Informatika", "", 1, "L", false, 0, "")
	pdf.Ln(30)
	pdf.SetX(-65)
	pdf.CellFormat(0, 5, "RONI ANDARSYAH", "", 1, "L", false, 0, "")
	pdf.SetX(-62)
	pdf.CellFormat(0, 5, "NIDN 0420058801", "", 1, "L", false, 0, "")

	// Save the PDF to a buffer without signature
	fileName := fmt.Sprintf("BAP-%s-%s.pdf", sanitizeFileName(data.MataKuliah), sanitizeFileName(data.Kelas))
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, "", err
	}

	return &buf, fileName, nil
}
