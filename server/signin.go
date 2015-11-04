package server

import (
	"fmt"
	"net/http"
)

type signin struct {}

//Handler for the authentication endpoint
func (s signin) ServeHTTP (rw http.ResponseWriter, rq *http.Request) {
	fmt.Println(rq.Header["SSL_CLIENT_VERIFY"])
	fmt.Println(rq.Header["SSL_CLIENT_S_DN"])
	rw.Write([]byte("Authentication endpoint"))
}