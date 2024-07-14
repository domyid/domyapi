package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

func getPdfUrl(fileName string) (string, error) {
	// Buat context baru dengan timeout
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// URL halaman repository
	url := "https://repo.ulbi.ac.id/buktiajar/#2023-2"

	// Variabel untuk menyimpan hasil
	var res string

	// Jalankan chromedp
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Sleep(2*time.Second), // Tambahkan delay untuk memastikan halaman termuat
		chromedp.Evaluate(`document.documentElement.outerHTML`, &res),
	)
	if err != nil {
		return "", err
	}

	// Cari URL yang sesuai
	startIndex := strings.Index(res, fileName)

	if startIndex == -1 {
		return "", fmt.Errorf("file not found")
	}
	hrefStart := strings.LastIndex(res[:startIndex], `href="`) + 6
	hrefEnd := strings.Index(res[hrefStart:], `"`) + hrefStart

	pdfURL := res[hrefStart:hrefEnd]

	if pdfURL == "" {
		return "", fmt.Errorf("failed to find PDF URL on repository page")
	}

	return pdfURL, nil
}

func main() {
	// Contoh penggunaan fungsi
	pdfURL, err := getPdfUrl("BAP-Kecerdasan_Buatan_Artificial_Intelligence-51.pdf")
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("PDF URL:", pdfURL)
	}
}
