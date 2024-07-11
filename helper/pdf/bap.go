package domyApi

import (
	"fmt"
	"strconv"

	model "github.com/domyid/domyapi/model"
	"github.com/jung-kurt/gofpdf"
)

const InfoImageURL = "https://home.ulbi.ac.id/ulbi.png"

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

func GenerateBAPPDF(data model.BAP) (string, error) {
	Text := []string{
		"UNIVERSITAS LOGISTIK DAN BISNIS INTERNASIONAL",
		"Berita Acara Perkuliahan dan Absensi Perkuliahan",
	}

	pdf := CreateHeaderBAP(Text, 90)
	pdf = ImageCustomize(pdf, "./ulbi.png", InfoImageURL, 28, 11, 35, 12, 100, 100, 0.3)

	// Header Information
	pdf.SetFont("Times", "", 12)
	pdf.CellFormat(0, 10, fmt.Sprintf("Kode Matakuliah/Nama Matakuliah: %s/%s", data.Kode, data.MataKuliah), "", 1, "", false, 0, "")
	pdf.CellFormat(0, 10, fmt.Sprintf("Kelas: %s", data.Kelas), "", 1, "", false, 0, "")
	pdf.CellFormat(0, 10, fmt.Sprintf("Semester/SKS: %s/%s SKS", data.SMT, data.SKS), "", 1, "", false, 0, "")
	pdf.Ln(10)

	// Tabel Log Aktivitas
	pdf.SetFont("Times", "B", 12)
	pdf.CellFormat(0, 10, "Tabel Log Aktivitas", "0", 1, "C", false, 0, "")
	pdf.Ln(5)

	logHeaders := []string{"Pertemuan", "Tanggal", "Jam", "Rencana Materi", "Realisasi Materi", "Pengajar", "Ruang", "Hadir", "Persentase"}
	widths := []float64{20, 30, 30, 60, 60, 50, 30, 20, 20}
	align := []string{"C", "C", "C", "C", "C", "C", "C", "C", "C"}

	pdf.SetFont("Times", "B", 10)
	for i, header := range logHeaders {
		pdf.CellFormat(widths[i], 10, header, "1", 0, align[i], false, 0, "")
	}
	pdf.Ln(-1)

	pdf.SetFont("Times", "", 10)
	for _, item := range data.RiwayatMengajar {
		pdf.CellFormat(widths[0], 10, item.Pertemuan, "1", 0, "C", false, 0, "")
		pdf.CellFormat(widths[1], 10, item.Tanggal, "1", 0, "C", false, 0, "")
		pdf.CellFormat(widths[2], 10, item.Jam, "1", 0, "C", false, 0, "")
		pdf.CellFormat(widths[3], 10, item.RencanaMateri, "1", 0, "C", false, 0, "")
		pdf.CellFormat(widths[4], 10, item.RealisasiMateri, "1", 0, "C", false, 0, "")
		pdf.CellFormat(widths[5], 10, item.Pengajar, "1", 0, "C", false, 0, "")
		pdf.CellFormat(widths[6], 10, item.Ruang, "1", 0, "C", false, 0, "")
		pdf.CellFormat(widths[7], 10, item.Hadir, "1", 0, "C", false, 0, "")
		pdf.CellFormat(widths[8], 10, item.Persentase, "1", 0, "C", false, 0, "")
		pdf.Ln(-1)
	}

	// Tabel Presensi
	pdf.Ln(10)
	pdf.SetFont("Times", "B", 12)
	pdf.CellFormat(0, 10, "Tabel Presensi", "0", 1, "C", false, 0, "")
	pdf.Ln(5)

	presensiHeaders := []string{"No", "NIM", "Nama", "Pertemuan", "Alfa", "Hadir", "Ijin", "Sakit", "Presentase"}
	widthsPresensi := []float64{10, 20, 40, 20, 10, 10, 10, 10, 20}
	alignPresensi := []string{"C", "C", "L", "C", "C", "C", "C", "C", "C"}

	pdf.SetFont("Times", "B", 10)
	for i, header := range presensiHeaders {
		pdf.CellFormat(widthsPresensi[i], 10, header, "1", 0, alignPresensi[i], false, 0, "")
	}
	pdf.Ln(-1)

	pdf.SetFont("Times", "", 10)
	for _, item := range data.AbsensiKelas {
		pdf.CellFormat(widthsPresensi[0], 10, item.No, "1", 0, "C", false, 0, "")
		pdf.CellFormat(widthsPresensi[1], 10, item.NIM, "1", 0, "C", false, 0, "")
		pdf.CellFormat(widthsPresensi[2], 10, item.Nama, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widthsPresensi[3], 10, item.Pertemuan, "1", 0, "C", false, 0, "")
		pdf.CellFormat(widthsPresensi[4], 10, item.Alfa, "1", 0, "C", false, 0, "")
		pdf.CellFormat(widthsPresensi[5], 10, item.Hadir, "1", 0, "C", false, 0, "")
		pdf.CellFormat(widthsPresensi[6], 10, item.Ijin, "1", 0, "C", false, 0, "")
		pdf.CellFormat(widthsPresensi[7], 10, item.Sakit, "1", 0, "C", false, 0, "")
		pdf.CellFormat(widthsPresensi[8], 10, item.Presentase, "1", 0, "C", false, 0, "")
		pdf.Ln(-1)
	}

	// Tabel Nilai Akhir
	pdf.Ln(10)
	pdf.SetFont("Times", "B", 12)
	pdf.CellFormat(0, 10, "Tabel Nilai Akhir", "0", 1, "C", false, 0, "")
	pdf.Ln(5)

	nilaiHeaders := []string{"No", "NIM", "Nama", "Hadir", "ATS", "AAS", "Nilai", "Grade", "Lulus", "Keterangan"}
	widthsNilai := []float64{10, 20, 40, 20, 20, 20, 20, 10, 10, 20}
	alignNilai := []string{"C", "C", "L", "C", "C", "C", "C", "C", "C", "L"}

	pdf.SetFont("Times", "B", 10)
	for i, header := range nilaiHeaders {
		pdf.CellFormat(widthsNilai[i], 10, header, "1", 0, alignNilai[i], false, 0, "")
	}
	pdf.Ln(-1)

	pdf.SetFont("Times", "", 10)
	for _, item := range data.ListNilai {
		hadir, _ := strconv.ParseFloat(item.Hadir, 64)
		ats, _ := strconv.ParseFloat(item.ATS, 64)
		aas, _ := strconv.ParseFloat(item.AAS, 64)
		nilai, _ := strconv.ParseFloat(item.Nilai, 64)

		pdf.CellFormat(widthsNilai[0], 10, item.No, "1", 0, "C", false, 0, "")
		pdf.CellFormat(widthsNilai[1], 10, item.NIM, "1", 0, "C", false, 0, "")
		pdf.CellFormat(widthsNilai[2], 10, item.Nama, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widthsNilai[3], 10, fmt.Sprintf("%.2f", hadir), "1", 0, "C", false, 0, "")
		pdf.CellFormat(widthsNilai[4], 10, fmt.Sprintf("%.2f", ats), "1", 0, "C", false, 0, "")
		pdf.CellFormat(widthsNilai[5], 10, fmt.Sprintf("%.2f", aas), "1", 0, "C", false, 0, "")
		pdf.CellFormat(widthsNilai[6], 10, fmt.Sprintf("%.2f", nilai), "1", 0, "C", false, 0, "")
		pdf.CellFormat(widthsNilai[7], 10, item.Grade, "1", 0, "C", false, 0, "")
		pdf.CellFormat(widthsNilai[8], 10, fmt.Sprintf("%t", item.Lulus), "1", 0, "C", false, 0, "")
		pdf.CellFormat(widthsNilai[9], 10, item.Keterangan, "1", 0, "L", false, 0, "")
		pdf.Ln(-1)
	}

	// Save the PDF to a file
	filePath := "bap.pdf"
	err := SavePDF(pdf, filePath)
	if err != nil {
		return "", err
	}
	return filePath, nil
}
