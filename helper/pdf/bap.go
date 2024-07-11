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

	color := []int{255, 255, 153}
	align := []string{"J", "C", "J"}

	pdf := CreateHeaderBAP(Text, 90)
	pdf = ImageCustomize(pdf, "./ulbi.png", InfoImageURL, 28, 11, 35, 12, 100, 100, 0.3)

	pdf.CellFormat(0, 10, fmt.Sprintf("Kode Matakuliah/Nama Matakuliah: %s/%s", data.Kode, data.MataKuliah), "", 1, "", false, 0, "")
	pdf.CellFormat(0, 10, fmt.Sprintf("Kelas: %s", data.Kelas), "", 1, "", false, 0, "")
	pdf.CellFormat(0, 10, fmt.Sprintf("Semester/SKS: %s/%s SKS", data.SMT, data.SKS), "", 1, "", false, 0, "")

	// Add Riwayat Mengajar table
	pdf.Ln(10)
	pdf = SetMergedCell(pdf, "Tabel Log Aktivitas", "J", 150, color)
	headers := []string{"Pertemuan", "Tanggal", "Jam", "Rencana Materi", "Realisasi Materi", "Pengajar", "Ruang", "Hadir", "Persentase"}
	widths := []float64{20, 30, 30, 50, 50, 40, 20, 10, 20}
	alignPertemuan := []string{"C", "C", "C", "L", "L", "L", "C", "C", "C"}
	pdf = SetHeaderTable(pdf, headers, widths, []int{135, 206, 235})
	for _, item := range data.RiwayatMengajar {
		row := []string{
			item.Pertemuan,
			item.Tanggal,
			item.Jam,
			item.RencanaMateri,
			truncateToThreeWords(item.RealisasiMateri),
			item.Pengajar,
			item.Ruang,
			item.Hadir,
			item.Persentase,
		}
		pdf = SetTableContent(pdf, [][]string{row}, widths, alignPertemuan)
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
	headers = []string{"No", "NIM", "Nama", "Hadir", "ATS", "AAS", "Nilai", "Grade", "Lulus", "Keterangan"}
	widths = []float64{10, 20, 40, 20, 20, 20, 20, 10, 10, 20}
	align = []string{"C", "C", "L", "C", "C", "C", "C", "C", "C", "C"}
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
			fmt.Sprintf("%t", item.Lulus),
			item.Keterangan,
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
