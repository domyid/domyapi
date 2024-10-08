package domyApi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	client "github.com/domyid/domyapi/client"
	model "github.com/domyid/domyapi/model"
	"github.com/jung-kurt/gofpdf"
	"github.com/skip2/go-qrcode"
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
	pdf.SetFont("Times", "", 8) // Font normal, ukuran 8
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

func createQRCode(link string, filename string) error {
	// Generate QR code
	err := qrcode.WriteFile(link, qrcode.Highest, 256, filename)
	if err != nil {
		return err
	}
	fmt.Printf("QR code generated and saved to %s\n", filename)
	return nil
}

func addSignature(pdf *gofpdf.Fpdf, qrFileName string) *gofpdf.Fpdf {
	// Menambahkan tempat tanda tangan dan QR code
	tanggalTerkini := "Bandung, " + getFormattedDate(time.Now())
	pdf.Ln(10)
	pdf.SetFont("Times", "", 10)
	pdf.SetX(-68) // Set posisi tanda tangan di sebelah kanan
	pdf.CellFormat(0, 5, tanggalTerkini, "", 1, "L", false, 0, "")
	pdf.SetX(-75)
	pdf.CellFormat(0, 5, "Ketua Prodi D4 Teknik Informatika", "", 1, "L", false, 0, "")

	// Set Y-coordinate explicitly for QR code
	pdf.SetY(pdf.GetY() + 5) // Adjust this value as needed for spacing

	// Menambahkan QR code
	pdf.SetX(-63)
	pdf.Image(qrFileName, pdf.GetX(), pdf.GetY(), 30, 30, false, "", 0, "")

	// Set Y-coordinate explicitly for signature text
	pdf.SetY(pdf.GetY() + 35) // Adjust this value as needed for spacing

	pdf.SetX(-65)
	pdf.CellFormat(0, 5, "RONI ANDARSYAH", "", 1, "L", false, 0, "")
	pdf.SetX(-62)
	pdf.CellFormat(0, 5, "NIDN 0420058801", "", 1, "L", false, 0, "")

	return pdf
}

func GenerateBAPPDF(data model.BAP, qrCodeLink string) (*bytes.Buffer, string, error) {
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

	// Generate QR code
	qrFileName := "signature_qrcode.png"
	err := createQRCode(qrCodeLink, qrFileName)
	if err != nil {
		return nil, "", err
	}

	// Add the signature section with QR code
	pdf = addSignature(pdf, qrFileName)

	// Save the PDF to a buffer
	fileName := fmt.Sprintf("BAP-%s-%s.pdf", sanitizeFileName(data.MataKuliah), sanitizeFileName(data.Kelas))
	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		return nil, "", err
	}

	return &buf, fileName, nil
}

func GenerateBAPPDFwithoutsignature(data model.BAP) (*bytes.Buffer, string, error) {
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

	// Save the PDF to a buffer
	fileName := fmt.Sprintf("BAP-%s-%s.pdf", sanitizeFileName(data.MataKuliah), sanitizeFileName(data.Kelas))
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, "", err
	}

	return &buf, fileName, nil
}

func GenerateBKD(data model.RekapBimbingan) (*bytes.Buffer, string, error) {
	// Text for the header
	Text := []string{
		"UNIVERSITAS LOGISTIK DAN BISNIS INTERNASIONAL",
		"Jl. Sari Asih No.54, Sarijadi, Kec. Sukasari, Kota Bandung, Jawa Barat 40151",
		"Website : www.ulbi.ac.id/ e-Mail :info@ulbi.ac.id / Telepon :081311110194",
	}

	// Initial setup for the PDF
	pdf := CreateHeaderBAP(Text, 90)

	// Adding more details (this can be customized further)
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 10, fmt.Sprintf("Judul Proposal: %s", data.JudulProposal))
	pdf.Ln(7)
	pdf.Cell(0, 10, fmt.Sprintf("Sesi / Bahasan: %s", data.SesiBahasan))
	pdf.Ln(7)

	// Splitting Mahasiswa into NIM and Name
	mahasiswaDetails := fmt.Sprintf("%s - %s", data.NIM, data.Mahasiswa)
	pdf.Cell(0, 10, fmt.Sprintf("Mahasiswa: %s", mahasiswaDetails))
	pdf.Ln(7)

	pdf.Cell(0, 10, fmt.Sprintf("Pembimbing Proposal: %s", data.PembimbingProposal))
	pdf.Ln(10)

	if data.Percakapan != "" {
		pdf.Cell(0, 10, "Percakapan:")
		pdf.Ln(7)
		pdf.MultiCell(0, 10, data.Percakapan, "", "", false)
	} else {
		pdf.Cell(0, 10, "Tidak ada data percakapan")
	}
	pdf.Ln(10)

	// Saving the PDF to a buffer
	fileName := fmt.Sprintf("BAP-%s.pdf", sanitizeFileName(data.NIM))
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, "", err
	}

	return &buf, fileName, nil
}

func CreateToken(docID, url string, data model.SignatureData) string {
	resp := new(model.TokenResp)
	body := new(model.RequestData)
	body.Id = docID
	body.Data = data

	res, err := client.CreateRequestHTTP().
		SetBody(body).
		Post(url)

	if err != nil {
		return "error ni kakak sistem akademiknya silahkan hubungi admin yaaaaa........"
	}

	defer res.Body.Close()
	_ = json.Unmarshal(res.Bytes(), &resp)

	return resp.Token
}

func GenerateLink(token string) string {
	return fmt.Sprintf("https://mrt.ulbi.ac.id/token/get?token=%s", token)
}
