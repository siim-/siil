
type Cert struct{
	Verified bool
	UserData string
}

func getCertFromHttpRequest(rq *http.Request) Cert {
	cert := Cert()
	if cert.Verified = clientVerified(rq); cert.Verified {
		cert.UserData = getUserData(rq)
	}
	return cert
}