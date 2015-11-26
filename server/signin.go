package server

import (
	//"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/siim-/siil/entity/site"
)

type signin struct{
	Site site.Entity
}

func handleSigninRequest(rw http.ResponseWriter, rq *http.Request) {
	reqVars := mux.Vars(rq)

	if siteId, ok := reqVars["site"]; !ok || len(siteId) == 0 {
		http.Error(rw, "Site ID must be provided", http.StatusBadRequest)
		return
	} else {
		site := site.Entity{ClientId: siteId}
		site.Load()
	}

	/*cert, err := cert.NewCertFromRequest(rq)
	if err != nil {
		fmt.Print(err)
		rw.Write([]byte("Not authed!"))
		return
	}

	response := fmt.Sprintf(
		"Hello, %s %s! You serial number is %s.",
		cert.FirstName,
		cert.LastName,
		cert.SerialNumber,
	)
	rw.Write([]byte(response))*/
}
