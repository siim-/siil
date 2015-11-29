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

//Handle the root request
func handleRootRequest(rw http.ResponseWriter, rq *http.Request) {
	rw.Write([]byte("Hello from Siil!"))
}

func handleOptionsRequest(rw http.ResponseWriter, rq *http.Request) {
	rw.Write([]byte("Hello from Siil options!"))
}

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

	//User primer & authentication handler
	baseRouter.HandleFunc("/signin/{site:[a-zA-Z0-9]*}", handleSigninRequest)
	baseRouter.HandleFunc("/id/{site:[a-zA-Z0-9]*}", handleSessionCreation)

	//Invalidate sessions
	baseRouter.Handle("/signout", signout{})

	http.ListenAndServe(fmt.Sprintf(":%d", port), baseRouter)
}
