package server

import (
	"net/http"

	"github.com/siim-/siil/cert"
)

type signin struct{}

//Handler for the authentication endpoint
func (s signin) ServeHTTP(rw http.ResponseWriter, rq *http.Request) {

	var cert cert.Cert = cert.NewCertFromRequest(rq)

	if cert.Verified {
		rw.Write([]byte(cert.UserData))
	} else {
		rw.Write([]byte("Not authed!"))
	}
	
}
