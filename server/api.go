package server

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/siim-/siil/cert"
	"github.com/siim-/siil/entity/session"
	"github.com/siim-/siil/entity/site"
	"github.com/siim-/siil/entity/user"
)

//API response containing information about the session and user
type apiResponse struct {
	Site         string       `json:"site_id"`
	Token        string       `json:"token"`
	ExpiresAt    string       `json:"expires_at"`
	CreatedAt    string       `json:"created_at"`
	User         *user.Entity `json:"user,omitempty"`
	Verification string       `json:"verification,omitempty"`
}

//Provide JSON encoded information about the user related to this session
func handleAPISessionRequest(rw http.ResponseWriter, rq *http.Request) {
	if rq.Method != "GET" {
		http.Error(rw, "Method not allowed", http.StatusMethodNotAllowed)
	} else {
		if token, clientId := rq.FormValue("token"), rq.FormValue("client_id"); len(token) != session.TOKEN_LENGTH || len(clientId) != site.CLIENT_ID_LENGTH {
			http.Error(rw, "Invalid request", http.StatusBadRequest)
		} else {
			if sess, err := session.GetSession(token); err != nil {
				log.Println(err)
				http.Error(rw, "Not found", http.StatusNotFound)
			} else {
				//Check the client id's
				if sess.SiteId != clientId {
					http.Error(rw, "Unauthorized", http.StatusUnauthorized)
					return
				}
				s := site.Entity{ClientId: sess.SiteId}
				if err := s.Load(); err != nil {
					log.Println(err)
					http.Error(rw, "Not found", http.StatusNotFound)
					return
				}

				//Lets check the X-Authorization header
				if sig := rq.Header.Get("X-Authorization"); len(sig) == 0 {
					http.Error(rw, "Unauthorized (Missing signature in X-Authorization)", http.StatusUnauthorized)
				} else {
					sigHMAC, err := base64.StdEncoding.DecodeString(sig)
					if err != nil {
						log.Println(err)
						http.Error(rw, "Invalid signature provided", http.StatusUnauthorized)
						return
					}
					ourMAC := hmac.New(sha512.New, []byte(s.PrivateKey))
					ourMAC.Write([]byte(fmt.Sprintf("%s\t%s", clientId, token)))

					if !hmac.Equal(sigHMAC, ourMAC.Sum(nil)) {
						http.Error(rw, "Invalid signature provided", http.StatusUnauthorized)
					} else {
						if usr, err := user.FindById(sess.UserId); err != nil {
							log.Println(err)
							http.Error(rw, "Not found", http.StatusNotFound)
						} else {
							response := apiResponse{
								Token:     token,
								ExpiresAt: sess.ExpiresAt.Format(time.RFC3339),
								CreatedAt: sess.CreatedAt.Format(time.RFC3339),
								User:      usr,
								Site:      clientId,
							}

							//Generate verification HMAC
							verHMAC := hmac.New(sha512.New, []byte(s.PrivateKey))

							if m, err := json.Marshal(response); err != nil {
								log.Println(err)
								http.Error(rw, "Failed to sign response", http.StatusInternalServerError)
							} else {
								verHMAC.Write(m)
							}

							response.Verification = base64.StdEncoding.EncodeToString(verHMAC.Sum(nil))
							enc := json.NewEncoder(rw)
							rw.Header().Set("Content-Type", "application/json")
							if err := enc.Encode(response); err != nil {
								log.Println(err)
								http.Error(rw, "Failed to compose response", http.StatusInternalServerError)
							}
						}
					}
				}
			}
		}
	}
}

//Provide information about the active id card user
func handleAPIMeRequest(rw http.ResponseWriter, rq *http.Request) {
	switch rq.Method {
	case "GET":
		if !cert.ClientVerified(rq) {
			http.Error(rw, "Certificate not provided", http.StatusBadRequest)
		} else {
			if c, err := cert.NewCertFromRequest(rq); err != nil {
				http.Error(rw, "Certificate not provided", http.StatusBadRequest)
			} else {
				if clientId := rq.FormValue("client_id"); len(clientId) == 0 {
					http.Error(rw, "Invalid client_id provided", http.StatusBadRequest)
				} else {
					//Check origin header validity
					if origin := rq.Header.Get("Origin"); len(origin) != 0 {
						if u, err := url.Parse(origin); err != nil {
							http.Error(rw, "Invalid origin provided", http.StatusBadRequest)
							return
						} else {
							s := site.Entity{Domain: u.Host}
							if err := s.Load(); err == nil && s.ClientId == clientId {
								rw.Header().Set("Access-Control-Allow-Origin", origin)
								rw.Header().Set("Access-Control-Allow-Methods", "GET")
								rw.Header().Set("Access-Control-Allow-Credentials", "true")
							} else {
								http.Error(rw, "Origin not allowed", http.StatusUnauthorized)
								return
							}
						}
					}
					s := site.Entity{ClientId: clientId}
					if err := s.Load(); err != nil {
						http.Error(rw, "Invalid client_id provided", http.StatusBadRequest)
						return
					}
					if usr, err := user.Find(c); err != nil {
						http.Error(rw, "Looks like we don't know you", http.StatusUnauthorized)
					} else {
						if s.HasActiveSessionFor(usr) {
							enc := json.NewEncoder(rw)
							rw.Header().Set("Content-Type", "application/json")
							if err := enc.Encode(usr); err != nil {
								log.Println(err)
								http.Error(rw, "Failed to compose response", http.StatusInternalServerError)
							}
						} else {
							http.Error(rw, "Computer says no", http.StatusUnauthorized)
						}
					}
				}
			}
		}
	default:
		http.Error(rw, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
