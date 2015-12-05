package server

import (
	"net/http"

	"github.com/siim-/siil/entity/session"
)

func getOwnerFromSession(rq *http.Request) (uint, error) {
	cookie, err := rq.Cookie("token")
	if err != nil {
		return 0, err
	}

	session, err := session.GetSession(cookie.Value)
	if err != nil {
		return 0, err
	}

	return uint(session.UserId), nil
}

func userLoggedIn(rq *http.Request) bool {
	_, err := getOwnerFromSession(rq)
	if err != nil {
		return false
	} else {
		return true
	}
}
