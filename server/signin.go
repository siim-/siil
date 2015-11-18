package server

import (
	"fmt"
	"net/http"

	"github.com/siim-/siil/cert"
)

type signin struct{}

//Handler for the authentication endpoint
func (s signin) ServeHTTP(rw http.ResponseWriter, rq *http.Request) {

	cert, err := cert.NewCertFromRequest(rq)
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
	rw.Write([]byte(response))
}
