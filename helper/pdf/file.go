package domyApi

import (
	"encoding/base64"
	"io"
	"net/http"
	"os"
)

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func PDFToBase64(filePath string) (string, error) {
	// Open the PDF file
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Read the PDF data
	fileInfo, _ := file.Stat()
	size := fileInfo.Size()
	pdfData := make([]byte, size)
	_, err = file.Read(pdfData)
	if err != nil {
		return "", err
	}

	// Encode the PDF data to Base64
	base64String := base64.StdEncoding.EncodeToString(pdfData)

	return base64String, nil
}
