package cert

import (
	"errors"
	"net/http"
	"strings"
)

func NewCertFromRequest(rq *http.Request) (*Cert, error) {
	if ClientVerified(rq) {
		var userData string = getUserData(rq)
		s := strings.Split(userData, ",")

		cert := &Cert{}
		// Serial number (EE - isikukood)
		ss := strings.Split(s[0], "=")
		cert.SerialNumber = ss[1]
		// First name
		ss = strings.Split(s[1], "=")
		cert.FirstName = capitalize(ss[1])
		// Last name
		ss = strings.Split(s[2], "=")
		cert.LastName = capitalize(ss[1])
		return cert, nil
	} else {
		return nil, errors.New("Client not verified.")
	}
}

func ClientVerified(rq *http.Request) bool {
	if verStatus := rq.Header["SSL_CLIENT_VERIFY"]; len(verStatus) > 0 {
		if verStatus[0] == "SUCCESS" {
			return true
		}
	}
	return false
}

func getUserData(rq *http.Request) string {
	return rq.Header["SSL_CLIENT_S_DN"][0]
}

func capitalize(word string) string {
	return strings.Title(strings.ToLower(word))
}
