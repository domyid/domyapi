package domyApi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"path/filepath"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	route "github.com/domyid/domyapi/route"
)

func init() {
	functions.HTTP("WebHook", route.URL)
}

func WebHook(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		Periode string `json:"periode"`
		Kelas   string `json:"kelas"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil || requestData.Periode == "" || requestData.Kelas == "" {
		http.Error(w, "Invalid request body or periode/kelas not provided", http.StatusBadRequest)
		return
	}

	fileName := fmt.Sprintf("BAP-%s-%s.pdf", requestData.Kelas, requestData.Periode)
	pdfURL, err := getPdfUrl(fileName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(pdfURL))
}

func getPdfUrl(fileName string) (string, error) {
	cmd := exec.Command("node", filepath.Join("nodejs", "index.js"), fileName)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}
