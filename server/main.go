package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/siim-/siil/entity"
	"github.com/siim-/siil/entity/session"
	"github.com/siim-/siil/entity/site"
	"github.com/siim-/siil/entity/user"

	"github.com/aymerick/raymond"
	"github.com/codegangsta/cli"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var templates map[string]*raymond.Template = make(map[string]*raymond.Template)

//Initialize the Siil API server
func StartAPIServer(c *cli.Context) {
	var (
		baseRouter *mux.Router
		port       int    = c.GlobalInt("port")
		tmplDir    string = fmt.Sprintf("%s/templates", c.GlobalString("wd"))
	)

	fmt.Printf("Starting API server on port %d...\n", port)
	if err := entity.CreateConnection(c.GlobalString("mysql")); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Loading templates...")
	if t, err := ioutil.ReadDir(tmplDir); err != nil {
		log.Fatal(err)
	} else {
		for _, file := range t {
			name := fmt.Sprintf("%s/%s", tmplDir, file.Name())
			if template, err := raymond.ParseFile(name); err != nil {
				log.Panic(err)
			} else {
				templates[file.Name()] = template
			}
		}
	}

	site.SIIL_SITE_ID = c.GlobalString("sid")

	baseRouter = mux.NewRouter()

	//Root endpoint doesn't really do anything
	baseRouter.HandleFunc("/", handleRootRequest)
	baseRouter.HandleFunc("/success", handleSuccessRequest)
	baseRouter.HandleFunc("/fail", handleCancelRequest)

	//User primer & authentication handlers
	baseRouter.HandleFunc("/signin/{site:[a-zA-Z0-9]*}", handleSigninRequest)
	baseRouter.HandleFunc("/id/{site:[a-zA-Z0-9]*}", handleSessionCreation)

	//New site creation
	baseRouter.HandleFunc("/addsite", handleAddSiteForm)
	baseRouter.HandleFunc("/addsite/fail", handleAddSiteFormFailed)
	baseRouter.HandleFunc("/addsite/success", handleAddSiteSuccess)
	baseRouter.HandleFunc("/addsite/request", handleAddSiteRequest)

	//Edit site
	baseRouter.HandleFunc("/editsite/success", handleEditSiteSuccess)
	baseRouter.HandleFunc("/editsite/{site:[a-zA-Z0-9]*}", handleEditSiteForm)
	baseRouter.HandleFunc("/editsite/{site:[a-zA-Z0-9]*}/fail", handleEditSiteFormFailed)
	baseRouter.HandleFunc("/editsite/{site:[a-zA-Z0-9]*}/request", handleEditSiteRequest)

	//Invalidate sessions
	baseRouter.HandleFunc("/signout/{token:[a-zA-Z0-9]*}", handleSignoutRequest)

	baseRouter.HandleFunc("/api/session", handleAPISessionRequest)
	baseRouter.HandleFunc("/api/me", handleAPIMeRequest)

	http.ListenAndServe(fmt.Sprintf(":%d", port), baseRouter)
}

//Handle the root request
func handleRootRequest(rw http.ResponseWriter, rq *http.Request) {
	//Detect if the user is authenticated with Siil
	authenticated := false
	token := ""
	if tokenCookie, err := rq.Cookie("token"); err == nil {
		token = tokenCookie.Value
		if sess, err := session.GetSession(token); err == nil {
			authenticated = sess.SiteId == site.SIIL_SITE_ID
		}
	}

	if authenticated {
		owner, err := getOwnerFromSession(rq)
		if err != nil {
			log.Println(err)
			http.Redirect(rw, rq, "/signin/"+site.SIIL_SITE_ID, http.StatusFound)
			return
		}

		usr, err := user.FindById(int(owner))
		if err != nil {
			log.Println(err)
			http.Error(rw, "Invalid user", http.StatusInternalServerError)
			return
		}

		var sections string = ""
		sites, err := site.GetUsersSites(owner)
		if err != nil {
			log.Println(err)
			return
		} else {
			for _, s := range sites {
				titleArea := surroundWithRow(title(s.Name) + button(s.ClientId))
				descriptionListArea := surroundWithRow(surroundWithColumn(getDescriptionList(&s), "twelve"))
				sections += surroundWithSection(titleArea + descriptionListArea)
			}
		}

		response, err := templates["index.hbs"].Exec(map[string]interface{}{
			"Sections":   sections,
			"authed":     authenticated,
			"first_name": usr.FirstName,
			"last_name":  usr.LastName,
			"site_id":    site.SIIL_SITE_ID,
			"token":      token,
		})
		if err != nil {
			http.Error(rw, "Something broke", http.StatusInternalServerError)
		} else {
			rw.Write([]byte(response))
		}
	} else {
		if t, err := templates["index.hbs"].Exec(map[string]interface{}{"authed": authenticated, "site_id": site.SIIL_SITE_ID, "token": token}); err != nil {
			http.Error(rw, "Something broke", http.StatusInternalServerError)
		} else {
			rw.Write([]byte(t))
		}
	}
}

func handleCancelRequest(rw http.ResponseWriter, rq *http.Request) {
	if t, err := templates["cancel.hbs"].Exec(map[string]interface{}{"site_id": site.SIIL_SITE_ID}); err != nil {
		http.Error(rw, "Something broke", http.StatusInternalServerError)
	} else {
		rw.Write([]byte(t))
	}
}
