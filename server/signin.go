package server

import (
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/siim-/siil/cert"
	"github.com/siim-/siil/entity/session"
	"github.com/siim-/siil/entity/site"
	"github.com/siim-/siil/entity/user"
)

func handleSigninRequest(rw http.ResponseWriter, rq *http.Request) {
	reqVars := mux.Vars(rq)

	if siteId, ok := reqVars["site"]; !ok || len(siteId) == 0 {
		http.Error(rw, "Site ID must be provided", http.StatusBadRequest)
		return
	} else {
		wanted := site.Entity{ClientId: siteId}
		if err := wanted.Load(); err != nil {
			log.Println(err)
			if strings.Contains(err.Error(), "no rows") {
				http.Error(rw, "Site not found", http.StatusNotFound)
			} else {
				http.Error(rw, "Something broke", http.StatusInternalServerError)
			}
			return
		}

		if result, err := templates["signin.hbs"].Exec(wanted); err != nil {
			log.Println(err)
			http.Error(rw, "Something broke", http.StatusInternalServerError)
			return
		} else {
			rw.Write([]byte(result))
		}
	}
}

func handleSessionCreation(rw http.ResponseWriter, rq *http.Request) {
	reqVars := mux.Vars(rq)

	if rq.Method != "POST" {
		http.Error(rw, "Method not allowed", http.StatusMethodNotAllowed)
	} else if siteId, ok := reqVars["site"]; !ok {
		http.Error(rw, "Site ID must be provided", http.StatusBadRequest)
	} else {
		wanted := site.Entity{ClientId: siteId}
		if err := wanted.Load(); err != nil {
			log.Println(err)
			if strings.Contains(err.Error(), "no rows") {
				http.Error(rw, "Site not found", http.StatusNotFound)
			} else {
				http.Error(rw, "Something broke", http.StatusInternalServerError)
			}
			return
		}

		if !cert.ClientVerified(rq) {
			http.Error(rw, "Client certificate not provided. Please restart your browser to provide it.", http.StatusBadRequest)
		} else {
			if userCert, err := cert.NewCertFromRequest(rq); err != nil {
				log.Println(err)
				http.Error(rw, "Failed to parse your client cert", http.StatusBadRequest)
			} else {
				if userEntity, err := user.FindOrCreate(userCert); err != nil {
					log.Println(err)
					http.Error(rw, "Something broke", http.StatusInternalServerError)
				} else {
					if sess, err := session.NewSession(&wanted, userEntity); err != nil {
						log.Println(err)
						http.Error(rw, "Something broke", http.StatusInternalServerError)
					} else {
						if callback, err := url.Parse(wanted.CallbackURL); err != nil {
							http.Error(rw, "Invalid callback URL provided", http.StatusInternalServerError)
						} else {
							//Indicate signin action with GET parameter
							q := callback.Query()
							q.Set("siil_action", "signin")
							callback.RawQuery = q.Encode()
							if t, err := templates["success.hbs"].Exec(map[string]string{"token": sess.Token, "callback": callback.String()}); err != nil {
								log.Println(err)
								http.Error(rw, "Something broke", http.StatusInternalServerError)
							} else {
								rw.Write([]byte(t))
							}
						}
					}
				}
			}
		}
	}
}

func handleSuccessRequest(rw http.ResponseWriter, rq *http.Request) {
	if rq.Method != "POST" {
		http.Error(rw, "Method not allowed", http.StatusMethodNotAllowed)
	} else {
		switch rq.FormValue("siil_action") {
		case "signin":
			if token := rq.FormValue("token"); len(token) != session.TOKEN_LENGTH {
				http.Error(rw, "Invalid token provided", http.StatusBadRequest)
			} else {
				if _, err := session.GetSession(token); err != nil {
					log.Fatal(err)
					http.Error(rw, "No session", http.StatusUnauthorized)
				} else {
					cookie := http.Cookie{
						Name:  "token",
						Value: token,
					}
					http.SetCookie(rw, &cookie)
					http.Redirect(rw, rq, "/", http.StatusFound)
				}
			}
		case "signout":
			//Just clear the token cookie
			cookie := http.Cookie{
				Name:    "token",
				Value:   "lel",
				Expires: time.Now().UTC().Add(time.Minute * -2),
			}
			http.SetCookie(rw, &cookie)
			http.Redirect(rw, rq, "/", http.StatusFound)
		default:
			http.Error(rw, "Invalid action provided", http.StatusBadRequest)
		}

	}
}
