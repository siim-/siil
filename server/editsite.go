package server

import (
	"errors"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/siim-/siil/entity/site"
)

func handleEditSiteForm(rw http.ResponseWriter, rq *http.Request) {
	reqVars := mux.Vars(rq)
	siteId, ok := reqVars["site"]

	if !ok || len(siteId) == 0 {
		http.Error(rw, "Site ID must be provided", http.StatusBadRequest)
		return
	}

	s, err := checkSiteAndUserConnetion(rq, siteId)
	if err != nil {
		log.Println(err)
		http.Redirect(rw, rq, "/signin/"+site.SIIL_SITE_ID, http.StatusFound)
		return
	}

	editForm(s, rw, rq, false)
}

func handleEditSiteFormFailed(rw http.ResponseWriter, rq *http.Request) {
	reqVars := mux.Vars(rq)
	siteId, ok := reqVars["site"]

	if !ok || len(siteId) == 0 {
		http.Error(rw, "Site ID must be provided", http.StatusBadRequest)
		return
	}

	s, err := checkSiteAndUserConnetion(rq, siteId)
	if err != nil {
		log.Println(err)
		http.Redirect(rw, rq, "/signin/"+site.SIIL_SITE_ID, http.StatusFound)
		return
	}

	editForm(s, rw, rq, true)
}

func handleEditSiteSuccess(rw http.ResponseWriter, rq *http.Request) {
	if t, err := templates["editsite-success.hbs"].Exec(map[string]interface{}{}); err != nil {
		http.Error(rw, "Something broke", http.StatusInternalServerError)
	} else {
		rw.Write([]byte(t))
	}
}

func handleEditSiteRequest(rw http.ResponseWriter, rq *http.Request) {
	reqVars := mux.Vars(rq)
	siteId, ok := reqVars["site"]

	if !ok || len(siteId) == 0 {
		http.Error(rw, "Site ID must be provided", http.StatusBadRequest)
		return
	}

	s, err := checkSiteAndUserConnetion(rq, siteId)
	if err != nil {
		log.Println(err)
		http.Redirect(rw, rq, "/editsite/"+siteId+"/fail", http.StatusFound)
		return
	}

	var params url.Values

	err = rq.ParseForm()
	if err != nil {
		log.Println(err)
		http.Redirect(rw, rq, "/editsite/"+siteId+"/fail", http.StatusFound)
		return
	} else {
		params = rq.PostForm
	}

	if len(params) != 0 {
		entry := site.Entry{
			ClientId:    s.ClientId,
			Owner:       s.Owner,
			Name:        params.Get("site-name"),
			Domain:      params.Get("domain-name"),
			CallbackURL: params.Get("callback-url"),
			CancelURL:   params.Get("cancel-url"),
		}

		_, err := site.EditSite(&entry)
		if err == nil {
			log.Println(err)
			http.Redirect(rw, rq, "/editsite/success", http.StatusFound)
			return
		}
	}

	http.Redirect(rw, rq, "/editsite/"+siteId+"/fail", http.StatusFound)
}

func editForm(s *site.Entity, rw http.ResponseWriter, rq *http.Request, displayError bool) {
	response, err := templates["editsite.hbs"].Exec(map[string]interface{}{
		"DisplayError": displayError,
		"ClientId":     s.ClientId,
		"Name":         s.Name,
		"Domain":       s.Domain,
		"CallbackURL":  s.CallbackURL,
		"CancelURL":    s.CancelURL,
	})
	if err != nil {
		http.Error(rw, "Something broke", http.StatusInternalServerError)
		return
	}
	rw.Write([]byte(response))
}

func checkSiteAndUserConnetion(rq *http.Request, siteId string) (*site.Entity, error) {
	owner, err := getOwnerFromSession(rq)
	if err != nil {
		return nil, err
	}

	s := site.Entity{ClientId: siteId}
	if err := s.Load(); err != nil {
		log.Println(err)
		return nil, errors.New("Failed to find site")
	}

	if s.Owner != owner {
		return nil, errors.New("Site owner doesn't match session owner.")
	}

	return &s, nil
}
