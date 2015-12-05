package server

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/siim-/siil/entity/session"
	"github.com/siim-/siil/entity/site"
)

func handleAddSiteForm(rw http.ResponseWriter, rq *http.Request) {
	if userLoggedIn(rq) {
		addForm(rw, rq, false)
	} else {
		http.Redirect(rw, rq, "/api/signin/a1s2d34", http.StatusFound)
	}
}

func handleAddSiteFormFailed(rw http.ResponseWriter, rq *http.Request) {
	if userLoggedIn(rq) {
		addForm(rw, rq, true)
	} else {
		http.Redirect(rw, rq, "/api/signin/a1s2d34", http.StatusFound)
	}
}

func handleAddSiteSuccess(rw http.ResponseWriter, rq *http.Request) {
	if t, err := templates["addsite-success.hbs"].Exec(map[string]interface{}{}); err != nil {
		http.Error(rw, "Something broke", http.StatusInternalServerError)
	} else {
		rw.Write([]byte(t))
	}
}

func handleAddSiteRequest(rw http.ResponseWriter, rq *http.Request) {
	var queryParams url.Values = rq.URL.Query()

	owner, err := getOwnerFromSession(rq)
	if err != nil {
		fmt.Println(err)
		http.Redirect(rw, rq, "/addsite/fail", http.StatusFound)
		return
	}

	if len(queryParams) != 0 {
		entry := site.Entry{
			Owner:       owner,
			Name:        queryParams.Get("site-name"),
			Domain:      queryParams.Get("domain-name"),
			CallbackURL: queryParams.Get("callback-url"),
			CancelURL:   queryParams.Get("cancel-url"),
		}

		_, err := site.NewSite(&entry)
		if err == nil {
			http.Redirect(rw, rq, "/addsite/success", http.StatusFound)
			return
		}
	}

	http.Redirect(rw, rq, "/addsite/fail", http.StatusFound)
}

func addForm(rw http.ResponseWriter, rq *http.Request, displayError bool) {
	response, err := templates["addsite.hbs"].Exec(map[string]interface{}{
		"displayError": displayError,
	})
	if err != nil {
		http.Error(rw, "Something broke", http.StatusInternalServerError)
		return
	}
	rw.Write([]byte(response))
}

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
