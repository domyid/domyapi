package main

import (
	"bytes"
	"fmt"
	"os"
	"time"

	"github.com/jung-kurt/gofpdf"
)

// CreateHeaderBAP generates the header for the BAP PDF
const InfoImageURL = "https://home.ulbi.ac.id/ulbi.png"
const SourceURL = "https://siakad.ulbi.ac.id/siakad/rep_perkuliahan"

// CreateHeaderBAP generates the header for the BAP PDF
func CreateHeaderBAP(Text []string) *gofpdf.Fpdf {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Times", "B", 12)

	// Set timezone to Asia/Jakarta
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		loc = time.FixedZone("WIB", 7*3600) // Default to WIB if timezone load fails
	}
	timestamp := time.Now().In(loc).Format("2006-01-02 15:04:05")

	// Add timestamp at top left
	pdf.SetFont("Times", "", 10)
	pdf.SetXY(10, 10) // X: 10 mm from left, Y: 10 mm from top
	pdf.CellFormat(0, 10, timestamp, "", 0, "L", false, 0, "")

	// Add source URL at top right
	pdf.SetXY(150, 10) // X: 150 mm from left, Y: 10 mm from top (A4 width is 210mm, right-aligned with some margin)
	pdf.CellFormat(0, 10, SourceURL, "", 0, "R", false, 0, "")

	// Set header text below the timestamp and source URL
	pdf.SetXY(70, 20) // Centered text (A4 width is 210mm)
	pdf.CellFormat(70, 10, Text[0], "0", 0, "C", false, 0, "")
	pdf.Ln(5)
	pdf.SetX(70)
	pdf.CellFormat(70, 10, Text[1], "0", 0, "C", false, 0, "")
	pdf.Ln(5)

	pdf.SetY(30)
	return pdf
}

func main() {
	Text := []string{
		"UNIVERSITAS LOGISTIK DAN BISNIS INTERNASIONAL",
		"Berita Acara Perkuliahan dan Absensi Perkuliahan",
	}

	pdf := CreateHeaderBAP(Text)

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		fmt.Println("Error generating PDF:", err)
		return
	}

	fileName := "header_bap.pdf"
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error creating PDF file:", err)
		return
	}
	defer file.Close()

	_, err = file.Write(buf.Bytes())
	if err != nil {
		fmt.Println("Error writing to PDF file:", err)
		return
	}

	fmt.Printf("PDF successfully created: %s\n", fileName)
}
