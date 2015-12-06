package server

import (
	"fmt"
	"net/http"
	"net/url"

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
		"DisplayError": displayError,
	})
	if err != nil {
		http.Error(rw, "Something broke", http.StatusInternalServerError)
		return
	}
	rw.Write([]byte(response))
}
