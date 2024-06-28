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
		at.WriteJSON(w, http.StatusForbidden, "No valid login found")
		return
	}

	token, err := helper.GetRefreshToken(client, login)
	if err != nil {
		if errors.Is(err, errors.New("no token found")) {
			at.WriteJSON(w, http.StatusForbidden, "token is invalid")
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		at.WriteJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	result := &model.ResponseAct{
		Login:     true,
		SxSession: token,
	}

	at.WriteJSON(w, http.StatusOK, result)
}
