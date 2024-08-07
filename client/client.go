package domyApi

import (
	"encoding/json"

	"github.com/imroc/req/v3"
)

func CreateClientHTTP() *req.Client {
	return req.
		C().
		SetJsonUnmarshal(json.Unmarshal).
		SetJsonMarshal(json.Marshal)
}

func CreateRequestHTTP() *req.Request {
	return CreateClientHTTP().
		R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept-Language", "en-US,en;q=0.8,id;q=0.6")
}
