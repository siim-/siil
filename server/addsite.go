package server

import (
	"log"
	"net/http"
	"net/url"

	"github.com/siim-/siil/entity/site"
)

func handleAddSiteForm(rw http.ResponseWriter, rq *http.Request) {
	if userLoggedIn(rq) {
		addForm(rw, rq, false)
	} else {
		http.Redirect(rw, rq, "/api/signin/" + site.SIIL_SITE_ID, http.StatusFound)
	}
}

func handleAddSiteFormFailed(rw http.ResponseWriter, rq *http.Request) {
	if userLoggedIn(rq) {
		addForm(rw, rq, true)
	} else {
		http.Redirect(rw, rq, "/api/signin/" + site.SIIL_SITE_ID, http.StatusFound)
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
	var params url.Values

	err := rq.ParseForm()
	if err != nil {
		log.Println(err)
		http.Redirect(rw, rq, "/addsite/fail", http.StatusFound)
		return
	} else {
		params = rq.PostForm
	}

	owner, err := getOwnerFromSession(rq)
	if err != nil {
		log.Println(err)
		http.Redirect(rw, rq, "/addsite/fail", http.StatusFound)
		return
	}

	if len(params) != 0 {
		entry := site.Entry{
			Owner:       owner,
			Name:        params.Get("site-name"),
			Domain:      params.Get("domain-name"),
			CallbackURL: params.Get("callback-url"),
			CancelURL:   params.Get("cancel-url"),
		}

		_, err := site.NewSite(&entry)
		if err == nil {
			log.Println(err)
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
