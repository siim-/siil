package server

import (
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/mux"
	"github.com/siim-/siil/cert"
	"github.com/siim-/siil/entity/session"
	"github.com/siim-/siil/entity/site"
	"github.com/siim-/siil/entity/user"
)

//Invalidate the session
func handleSignoutRequest(rw http.ResponseWriter, rq *http.Request) {
	if rq.Method != "POST" {
		http.Error(rw, "Invalid method", http.StatusMethodNotAllowed)
	} else {
		reqVars := mux.Vars(rq)
		if token, ok := reqVars["token"]; !ok || len(token) != session.TOKEN_LENGTH {
			http.Error(rw, "Bad token provided", http.StatusBadRequest)
		} else {
			if sess, err := session.GetSession(token); err != nil {
				http.Error(rw, "Session not found", http.StatusUnauthorized)
			} else {
				if !cert.ClientVerified(rq) {
					http.Error(rw, "Cert not provided", http.StatusBadRequest)
				} else {
					if userCert, err := cert.NewCertFromRequest(rq); err != nil {
						log.Println(err)
						http.Error(rw, "Failed to parse your client cert", http.StatusBadRequest)
					} else {
						if usr, err := user.Find(userCert); err != nil {
							log.Println(err)
							http.Error(rw, "We don't know you", http.StatusUnauthorized)
						} else {
							if usr.Id == sess.UserId {
								if err := sess.Delete(); err != nil {
									http.Error(rw, "Failed to end session", http.StatusInternalServerError)
								} else {
									wanted := site.Entity{ClientId: sess.SiteId}
									if err := wanted.Load(); err != nil {
										log.Println(err)
										if strings.Contains(err.Error(), "no rows") {
											http.Error(rw, "Site not found", http.StatusNotFound)
										} else {
											http.Error(rw, "Something broke", http.StatusInternalServerError)
										}
										return
									}
									if callback, err := url.Parse(wanted.CallbackURL); err != nil {
										http.Error(rw, "Invalid callback URL provided", http.StatusInternalServerError)
									} else {
										//Indicate signin action with GET parameter
										q := callback.Query()
										q.Set("siil_action", "signout")
										callback.RawQuery = q.Encode()
										if t, err := templates["success.hbs"].Exec(map[string]string{"token": sess.Token, "callback": callback.String()}); err != nil {
											log.Println(err)
											http.Error(rw, "Something broke", http.StatusInternalServerError)
										} else {
											rw.Write([]byte(t))
										}
									}
								}
							} else {
								http.Error(rw, "Session and user don't match", http.StatusBadRequest)
							}
						}
					}
				}
			}
		}
	}
}
