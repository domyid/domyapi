package domyApi

import (
	"fmt"
	"net/http"
)

func GetRefreshToken(client *http.Client, token string) (string, error) {
	homeURL := "https://siakad.ulbi.ac.id/siakad/data_bimbingan/detail/112"

	req, err := http.NewRequest("GET", homeURL, nil)
	if err != nil {
		return "", nil
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Referer", "https://siakad.ulbi.ac.id/gate/menu")
	req.Header.Set("Cookie", fmt.Sprintf("SIAKAD_CLOUD_ACCESS=%s", token))

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	tokenStr := resp.Header.Get("Sx-Session")

	if tokenStr == "" {
		return "", fmt.Errorf("no token found")
	}

	return tokenStr, nil
}
