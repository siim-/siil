package server

import (
	"fmt"
	"net/http"

	"github.com/codegangsta/cli"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/go-sql-driver/mysql"
)
var DB *sqlx.DB

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
		port       int = c.GlobalInt("port")
	)

	_, err := sqlx.Connect("mysql", c.GlobalString("mysql"))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Starting API server on port %d...\n", port)

	baseRouter = mux.NewRouter()

	//Root endpoint doesn't really do anything
	baseRouter.HandleFunc("/", handleRootRequest)

	//Authentication handler for new sessions
	baseRouter.HandleFunc("/signin/{site:[a-zA-Z0-9]*}", handleSigninRequest)

	//Invalidate sessions
	baseRouter.Handle("/signout", signout{})

	http.ListenAndServe(fmt.Sprintf(":%d", port), baseRouter)
}
