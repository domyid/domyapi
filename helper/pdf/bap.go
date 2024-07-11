package domyApi

import (
	"fmt"
	"strconv"

	model "github.com/domyid/domyapi/model"
	"github.com/jung-kurt/gofpdf"
)

const InfoImageURL = "https://home.ulbi.ac.id/ulbi.png"

// CreateHeaderBAP generates the header for the BAP PDF
func CreateHeaderBAP(Text []string, x float64) *gofpdf.Fpdf {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Times", "B", 12)
	pdf.SetX(x)
	pdf.CellFormat(70, 10, Text[0], "0", 0, "C", false, 0, "")
	pdf.Ln(5)
	pdf.SetX(x)
	pdf.CellFormat(70, 10, Text[1], "0", 0, "C", false, 0, "")
	pdf.Ln(5)
	pdf.SetY(20)
	return pdf
}

// GenerateBAPPDF generates the BAP PDF
func GenerateBAPPDF(data model.BAP) (string, error) {
	Text := []string{
		"UNIVERSITAS LOGISTIK DAN BISNIS INTERNASIONAL",
		"Berita Acara Perkuliahan dan Absensi Perkuliahan",
	}

	width := []float64{60, 5, 70}
	color := []int{255, 255, 153}
	align := []string{"J", "C", "J"}
	yCoordinates := []float64{40, 45, 50}

	pdf := CreateHeaderBAP(Text, 90)
	pdf = ImageCustomize(pdf, "./ulbi.png", InfoImageURL, 28, 11, 35, 12, 100, 100, 0.3)

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
	widths := []float64{18, 25, 22, 40, 40}
	align = []string{"C", "C", "C", "C", "C"}
	pdf = SetHeaderTable(pdf, headers, widths, []int{135, 206, 235})
	for _, item := range data.RiwayatMengajar {
		pdf.CellFormat(widths[0], 10, item.Pertemuan, "1", 0, align[0], false, 0, "")
		pdf.CellFormat(widths[1], 10, item.Tanggal, "1", 0, align[1], false, 0, "")
		pdf.CellFormat(widths[2], 10, item.Jam, "1", 0, align[2], false, 0, "")
		// Use MultiCell for "Rencana Materi" and "Realisasi Materi"
		x := pdf.GetX()
		y := pdf.GetY()
		pdf.MultiCell(widths[3], 10, item.RencanaMateri, "1", align[3], false)
		pdf.SetXY(x+widths[3], y)
		pdf.MultiCell(widths[4], 10, item.RealisasiMateri, "1", align[4], false)
		pdf.Ln(-1)
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

	// Save the PDF to a file
	filePath := "bap.pdf"
	err := SavePDF(pdf, filePath)
	if err != nil {
		return "", err
	}
	return filePath, nil
}
