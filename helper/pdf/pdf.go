package domyApi

import (
	"github.com/jung-kurt/gofpdf"
)

func AddHeadText(pdf *gofpdf.Fpdf, spacing, x float64, align, text string) *gofpdf.Fpdf {
	pdf.SetFont("Times", "B", 9)

	// Get the current Y position
	currentY := pdf.GetY()

	// Set the X position
	pdf.SetX(x)
	// Add the text
	pdf.CellFormat(0, 10, text, "0", 1, align, false, 0, "")
	//pdf.Text(147, 140, "Juru Bayar")

	// Adjust the Y position to create spacing
	pdf.SetY(currentY + spacing)

	return pdf
}

func AddNameText(pdf *gofpdf.Fpdf, Text string, spacing, x, size float64) *gofpdf.Fpdf {

	pdf.SetFont("Times", "B", size)
	//pdf.Text(137, 138, Text)
	pdf.SetX(x)
	pdf.CellFormat(0, 10, Text, "0", 0, "C", false, 0, "")
	pdf.Ln(0.5 * size)

	currentY := pdf.GetY()

	pdf.SetY(currentY + spacing)

	return pdf
}

func SetMergedCell(pdf *gofpdf.Fpdf, text, align string, width float64, rgb []int) *gofpdf.Fpdf {
	pdf.SetFont("Times", "B", 10)
	pdf.SetFillColor(rgb[0], rgb[1], rgb[2])
	totalWidth := 0.0
	totalWidth += width

	// Calculate the X-coordinate to center the table on the page
	pageWidth, _ := pdf.GetPageSize()
	x := (pageWidth - totalWidth) / 2

	// Set the X-coordinate
	pdf.SetX(x)

	// Create 6 cells that make up the merged cell
	pdf.CellFormat(width, 7, text, "1", 0, align, true, 0, "")

	// Move to the next line
	pdf.Ln(-1)
	return pdf
}

func SetMergedCellSkyBlue(pdf *gofpdf.Fpdf, text string, width float64) *gofpdf.Fpdf {
	pdf.SetFont("Times", "B", 10)
	pdf.SetFillColor(135, 206, 235)
	totalWidth := 0.0
	totalWidth += width

	// Calculate the X-coordinate to center the table on the page
	pageWidth, _ := pdf.GetPageSize()
	x := (pageWidth - totalWidth) / 2

	// Set the X-coordinate
	pdf.SetX(x)

	// Create 6 cells that make up the merged cell
	pdf.CellFormat(width, 7, text, "1", 0, "L", true, 0, "")

	// Move to the next line
	pdf.Ln(-1)
	return pdf
}

func SetHeaderTable(pdf *gofpdf.Fpdf, hdr []string, widths []float64, rgb []int) *gofpdf.Fpdf {
	pdf.SetFont("Times", "B", 8)
	pdf.SetFillColor(rgb[0], rgb[1], rgb[2])
	// Calculate the total width of the table
	totalWidth := 0.0
	for _, width := range widths {
		totalWidth += width
	}

	// Calculate the X-coordinate to center the table on the page
	pageWidth, _ := pdf.GetPageSize()
	x := (pageWidth - totalWidth) / 2

	// Set the X-coordinate
	pdf.SetX(x)
	for i, str := range hdr {
		// The `CellFormat()` method takes a couple of parameters to format
		// the cell. We make use of this to create a visible border around
		// the cell, and to enable the background fill.
		pdf.CellFormat(widths[i], 7, str, "1", 0, "C", true, 0, "")
	}

	// Passing `-1` to `Ln()` uses the height of the last printed cell as
	// the line height.
	pdf.Ln(-1)
	return pdf
}

func SetTableContent(pdf *gofpdf.Fpdf, tbl [][]string, widths []float64, align []string) *gofpdf.Fpdf {
	pdf.SetFont("Times", "", 8)
	pdf.SetFillColor(255, 255, 255)

	for _, line := range tbl {
		// Calculate the total width of the table
		totalWidth := 0.0
		for _, width := range widths {
			totalWidth += width
		}

		// Calculate the X-coordinate to center the table on the page
		pageWidth, _ := pdf.GetPageSize()
		x := (pageWidth - totalWidth) / 2

		// Set the X-coordinate
		pdf.SetX(x)
		for i, str := range line {
			pdf.CellFormat(widths[i], 7, str, "1", 0, align[i], true, 0, "")
		}
		pdf.Ln(-1)
	}
	return pdf
}

func ImagePdf(pdf *gofpdf.Fpdf, filename, urlimage string) *gofpdf.Fpdf {
	if !FileExists(filename) {
		DownloadFile(filename, urlimage)
	}
	pdf.ImageOptions(filename, 12, 16, 20, 10, false, gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}, 0, "")
	return pdf
}

func SavePDF(pdf *gofpdf.Fpdf, path string) error {
	return pdf.OutputFileAndClose(path)
}

func SignatureImage(pdf *gofpdf.Fpdf, filename string, x, spacing float64, textlines []string, textYOffset float64) *gofpdf.Fpdf {
	currentY := pdf.GetY()
	y := currentY + spacing

	pdf.ImageOptions(filename, x, y, 30, 30, false, gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}, 0, "")
	pdf.Ln(-1)

	textX := x - 42          // Use the same x position as the image
	textY := y + textYOffset // Adjust the Y position based on the offset

	pdf.SetFont("Times", "B", 9)
	pdf.SetXY(textX-2, textY-6)

	// Add two lines of text using CellFormat
	pdf.CellFormat(0, 10, textlines[0], "0", 0, "C", false, 0, "")

	// Add the second line of text using Text
	pdf.Text(textX+56, textY+5, textlines[1])

	return pdf
}

func AddText(pdf *gofpdf.Fpdf, x, y float64, text string) *gofpdf.Fpdf {
	pdf.SetFont("Times", "", 9)
	pdf.Text(x, y, text)
	return pdf
}

func ImageCustomize(pdf *gofpdf.Fpdf, filename, urlimage string, x, y, w, h, wimg, himg, borderWidth float64) *gofpdf.Fpdf {
	if !FileExists(filename) {
		DownloadFile(filename, urlimage)
	}

	// Draw the image
	pdf.ImageOptions(filename, x, y, w, h, false, gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}, 0, "")

	return pdf
}

func AddTextCustomSize(pdf *gofpdf.Fpdf, x, y, size float64, text string) *gofpdf.Fpdf {
	pdf.SetFont("Times", "", size)
	pdf.Text(x, y, text)
	return pdf
}

func SetTableContentCustomY(pdf *gofpdf.Fpdf, tbl [][]string, widths []float64, align []string, customY []float64) *gofpdf.Fpdf {
	pdf.SetFont("Times", "", 10)
	pdf.SetFillColor(255, 255, 255)

	for i, line := range tbl {
		// Calculate the total width of the table
		totalWidth := 0.0
		for _, width := range widths {
			totalWidth += width
		}

		// Calculate the X-coordinate to center the table on the page

		x := 30.0

		// Set the X-coordinate and custom Y-coordinate
		pdf.SetXY(x, customY[i])

		for j, str := range line {
			pdf.CellFormat(widths[j], 7, str, "0", 0, align[j], true, 0, "")
		}
		pdf.Ln(-1)
	}
	return pdf
}
