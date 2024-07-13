package domyApi

import (
	"bytes"
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

func GenerateBAPPDF(data model.BAP) (*bytes.Buffer, string, error) {
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

	// Save the PDF to a buffer
	fileName := fmt.Sprintf("BAP-%s-%s.pdf", sanitizeFileName(data.MataKuliah), sanitizeFileName(data.Kelas))
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, "", err
	}

	return &buf, fileName, nil
}
