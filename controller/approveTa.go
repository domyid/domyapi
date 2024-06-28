package domyApi

import (
	"errors"
	"net/http"
	"net/http/cookiejar"

	at "github.com/domyid/domyapi/helper/at"
	helper "github.com/domyid/domyapi/helper/atapi"
	model "github.com/domyid/domyapi/model"
)

func SaveTokenString(w http.ResponseWriter, reg *http.Request) {
	jar, _ := cookiejar.New(nil)

	// Create a new HTTP client with the cookie jar
	client := &http.Client{
		Jar: jar,
	}

	login := at.GetLoginFromHeader(reg)
	if login == "" {
		http.Error(w, "No valid login found", http.StatusForbidden)
		return
	}

	token, err := helper.GetRefreshToken(client, login)
	if err != nil {
		if errors.Is(err, errors.New("no token found")) {
			http.Error(w, "token is invalid", http.StatusForbidden)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result := &model.ResponseAct{
		Login:     true,
		SxSession: token,
	}

	at.WriteJSON(w, http.StatusOK, result)
}
