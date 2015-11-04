package server

import (
	"net/http"
)

type signin struct{}

//Handler for the authentication endpoint
func (s signin) ServeHTTP(rw http.ResponseWriter, rq *http.Request) {

	var client Client getCertFromHttpRequest(rq)

	if client.Verified {
		rw.Write([]byte(client.UserData))
	} else {
		rw.Write([]byte("Not authed!"))
	}

	// if verStatus := rq.Header["SSL_CLIENT_VERIFY"]; len(verStatus) > 0 {
	// 	if verStatus[0] == "SUCCESS" {
	// 		rw.Write([]byte(rq.Header["SSL_CLIENT_S_DN"][0]))
	// 	} else {
	// 		rw.Write([]byte("Not authed!"))
	// 	}
	// }
}