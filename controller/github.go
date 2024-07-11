package domyApi

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"

	config "github.com/domyid/domyapi/config"
	helper "github.com/domyid/domyapi/helper"
	ghupload "github.com/domyid/domyapi/helper/ghupload"
	model "github.com/domyid/domyapi/model"
	"go.mongodb.org/mongo-driver/bson"

	atdb "github.com/domyid/domyapi/helper/atdb"
)

func GetGithubFiles(w http.ResponseWriter, r *http.Request) {
	var respn model.Response

	// Mendapatkan kredensial GitHub dari database
	gh, err := atdb.GetOneDoc[model.Ghcreates](config.Mongoconn, "github", bson.M{})
	if err != nil {
		respn.Info = helper.GetSecretFromHeader(r)
		respn.Response = err.Error()
		helper.WriteJSON(w, http.StatusConflict, respn)
		return
	}

	// Mendapatkan daftar file dari repositori GitHub
	content, err := ghupload.GithubListFiles(gh.GitHubAccessToken, "alittifaq", "cdn", "img")
	if err != nil {
		respn.Response = err.Error()
		helper.WriteJSON(w, http.StatusInternalServerError, respn)
		return
	}

	// Tambahkan logging untuk melihat data yang dikembalikan
	fmt.Printf("GetGithubFiles: %v\n", content)

	// Mengonversi content ke JSON string
	contentJSON, err := json.Marshal(content)
	if err != nil {
		respn.Response = err.Error()
		helper.WriteJSON(w, http.StatusInternalServerError, respn)
		return
	}

	respn.Info = "Files retrieved successfully"
	respn.Response = string(contentJSON)
	helper.WriteJSON(w, http.StatusOK, respn)
}

func PostUploadGithub(w http.ResponseWriter, r *http.Request) {
	var respn model.Response
	// _, err := watoken.Decode(config.PublicKeyWhatsAuth, helper.GetLoginFromHeader(req))
	// if err != nil {
	// 	respn.Info = helper.GetSecretFromHeader(req)
	// 	respn.Response = err.Error()
	// 	helper.WriteJSON(respw, http.StatusForbidden, respn)
	// 	return
	// }
	// Parse the form file
	_, header, err := r.FormFile("image")
	if err != nil {

		respn.Response = err.Error()
		helper.WriteJSON(w, http.StatusBadRequest, respn)
		return
	}

	//folder := ctx.Params("folder")
	folder := helper.GetParam(r)
	var pathFile string
	if folder != "" {
		pathFile = folder + "/" + header.Filename
	} else {
		pathFile = header.Filename
	}

	// save to github
	gh, err := atdb.GetOneDoc[model.Ghcreates](config.Mongoconn, "github", bson.M{})
	if err != nil {
		respn.Info = helper.GetSecretFromHeader(r)
		respn.Response = err.Error()
		helper.WriteJSON(w, http.StatusConflict, respn)
		return
	}

	content, _, err := ghupload.GithubUpload(gh.GitHubAccessToken, gh.GitHubAuthorName, gh.GitHubAuthorEmail, header, "alittifaq", "cdn", pathFile, false)
	if err != nil {
		respn.Info = "gagal upload github"
		respn.Response = err.Error()
		helper.WriteJSON(w, http.StatusEarlyHints, content)
		return
	}
	respn.Info = *content.Content.Name
	respn.Response = *content.Content.Path
	helper.WriteJSON(w, http.StatusOK, respn)

}

func UpdateGithubFile(w http.ResponseWriter, r *http.Request) {
	var respn model.Response

	// Parse multipart form
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		respn.Response = err.Error()
		helper.WriteJSON(w, http.StatusBadRequest, respn)
		return
	}

	// Get the uploaded file
	file, handler, err := r.FormFile("file")
	if err != nil {
		respn.Response = err.Error()
		helper.WriteJSON(w, http.StatusBadRequest, respn)
		return
	}
	defer file.Close()

	// Get the file name from form
	fileName := r.FormValue("fileName")
	if fileName == "" {
		respn.Response = "File name is required"
		helper.WriteJSON(w, http.StatusBadRequest, respn)
		return
	}

	// Get GitHub credentials from the database
	gh, err := atdb.GetOneDoc[model.Ghcreates](config.Mongoconn, "github", bson.M{})
	if err != nil {
		respn.Info = helper.GetSecretFromHeader(r)
		respn.Response = err.Error()
		helper.WriteJSON(w, http.StatusConflict, respn)
		return
	}

	// Create a multipart.FileHeader from the uploaded file
	fileHeader := &multipart.FileHeader{
		Filename: handler.Filename,
		Header:   handler.Header,
		Size:     handler.Size,
	}

	// Update the file in GitHub
	content, _, err := ghupload.GithubUpdateFile(gh.GitHubAccessToken, gh.GitHubAuthorName, gh.GitHubAuthorEmail, fileHeader, "alittifaq", "cdn", fileName)
	if err != nil {
		respn.Info = "Failed to update GitHub file"
		respn.Response = err.Error()
		helper.WriteJSON(w, http.StatusInternalServerError, respn)
		return
	}

	respn.Info = "File updated successfully"
	respn.Response = *content.Content.Path
	helper.WriteJSON(w, http.StatusOK, respn)
}

func DeleteGithubFile(w http.ResponseWriter, r *http.Request) {
	var respn model.Response
	var deleteRequest struct {
		FileName string `json:"fileName"`
	}

	if err := json.NewDecoder(r.Body).Decode(&deleteRequest); err != nil {
		respn.Response = err.Error()
		helper.WriteJSON(w, http.StatusBadRequest, respn)
		return
	}

	gh, err := atdb.GetOneDoc[model.Ghcreates](config.Mongoconn, "github", bson.M{})
	if err != nil {
		respn.Info = helper.GetSecretFromHeader(r)
		respn.Response = err.Error()
		helper.WriteJSON(w, http.StatusConflict, respn)
		return
	}

	_, _, err = ghupload.GithubDeleteFile(gh.GitHubAccessToken, gh.GitHubAuthorName, gh.GitHubAuthorEmail, "alittifaq", "cdn", deleteRequest.FileName)
	if err != nil {
		respn.Info = "Failed to delete GitHub file"
		respn.Response = err.Error()
		helper.WriteJSON(w, http.StatusInternalServerError, respn)
		return
	}

	respn.Info = "File deleted successfully"
	respn.Response = deleteRequest.FileName
	helper.WriteJSON(w, http.StatusOK, respn)
}
