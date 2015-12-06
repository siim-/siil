package server

import (
	"log"
	"net/http"

	"github.com/siim-/siil/entity/site"
)

func handleSites(rw http.ResponseWriter, rq *http.Request) {
	owner, err := getOwnerFromSession(rq)
	if err != nil {
		log.Println(err)
		respondWithSections(rw, "")
		return
	}

	sites, err := site.GetUsersSites(owner)
	if err != nil {
		log.Println(err)
		respondWithSections(rw, "")
		return
	}

	var sections string = ""
	for _, s := range sites {
		var html string = ""
		html += surrondWithH3(s.Name)
		html += dataMarkup("Domain", s.Domain)
		html += dataMarkup("Client ID", s.ClientId)
		html += dataMarkup("Private Key", s.PrivateKey)
		html += dataMarkup("Callback URL", s.CallbackURL)
		html += dataMarkup("Cancel URL", s.CancelURL)
		sections += surroundWithDl(html)
	}

	respondWithSections(rw, sections)
}

func respondWithSections(rw http.ResponseWriter, sections string) {
	response, err := templates["sites.hbs"].Exec(map[string]interface{}{
		"Sections": sections,
	})
	if err != nil {
		http.Error(rw, "Something broke", http.StatusInternalServerError)
	} else {
		rw.Write([]byte(response))
	}
}

func dataMarkup(dt, dd string) string {
	var str string = ""
	str += surroundWithDt(dt)
	str += surroundWithDd(dd)
	return str
}

func surrondWithH3(str string) string {
	return "<h3>" + str + "</h3>"
}

func surroundWithDl(str string) string {
	return "<dl>" + str + "</dl>"
}

func surroundWithDt(str string) string {
	return "<dt>" + str + "</dt>"
}

func surroundWithDd(str string) string {
	return "<dd>" + str + "</dd>"
}
