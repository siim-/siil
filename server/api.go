package server

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/siim-/siil/entity/session"
	"github.com/siim-/siil/entity/site"
	"github.com/siim-/siil/entity/user"
)

//API response containing information about the session and user
type apiResponse struct {
	Token     string       `json:"token"`
	ExpiresAt string       `json:"expires_at"`
	CreatedAt string       `json:"created_at"`
	User      *user.Entity `json:"user,omitempty"`
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
							timeFormat := "2006-01-02 15:04:05T-0700"
							response := apiResponse{
								Token:     token,
								ExpiresAt: sess.ExpiresAt.Format(timeFormat),
								CreatedAt: sess.CreatedAt.Format(timeFormat),
								User:      usr,
							}

							enc := json.NewEncoder(rw)
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

}
