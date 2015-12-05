package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/siim-/siil/entity"

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

	baseRouter = mux.NewRouter()

	//Root endpoint doesn't really do anything
	baseRouter.HandleFunc("/", handleRootRequest)
	baseRouter.HandleFunc("/success", handleSuccessRequest)

	//User primer & authentication handlers
	baseRouter.HandleFunc("/signin/{site:[a-zA-Z0-9]*}", handleSigninRequest)
	baseRouter.HandleFunc("/id/{site:[a-zA-Z0-9]*}", handleSessionCreation)

	//Invalidate sessions
	baseRouter.HandleFunc("/signout/{token:[a-zA-Z0-9]*}", handleSignoutRequest)

	baseRouter.HandleFunc("/api/session/{token:[a-zA-Z0-9]*}", handleAPISessionRequest)
	baseRouter.HandleFunc("/api/me", handleAPIMeRequest)

	http.ListenAndServe(fmt.Sprintf(":%d", port), baseRouter)
}

//Handle the root request
func handleRootRequest(rw http.ResponseWriter, rq *http.Request) {
	if t, err := templates["index.hbs"].Exec(map[string]interface{}{}); err != nil {
		http.Error(rw, "Something broke", http.StatusInternalServerError)
	} else {
		rw.Write([]byte(t))
	}
}
