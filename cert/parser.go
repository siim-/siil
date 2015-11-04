package cert

import (
	"net/http"
)

func NewCertFromRequest(rq *http.Request) Cert {
	cert := Cert{}
	if cert.Verified = clientVerified(rq); cert.Verified {
		cert.UserData = getUserData(rq)
	}
	return cert
}

func clientVerified(rq *http.Request) bool {
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