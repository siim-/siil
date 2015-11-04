package server

import (
	"fmt"
	"net/http"

	"github.com/siim-/siil/cert"
)

type signin struct{}

//Handler for the authentication endpoint
func (s signin) ServeHTTP(rw http.ResponseWriter, rq *http.Request) {

	var cert cert.Cert = cert.NewCertFromRequest(rq)

	if cert.Verified {
		response := fmt.Sprintf(
			"Hello, %s %s! You serial number is %s.",
			cert.FirstName,
			cert.LastName,
			cert.SerialNumber,
		)
		rw.Write([]byte(response))
	} else {
		rw.Write([]byte("Not authed!"))
	}
}
