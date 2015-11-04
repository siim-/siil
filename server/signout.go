package server

import (
	"net/http"
)

type signout struct {}

//Handler for the authentication endpoint
func (s signout) ServeHTTP (rw http.ResponseWriter, rq *http.Request) {
	rw.Write([]byte("Authentication endpoint"))
}