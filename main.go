package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Violent-Idiot/Go-Url_Shortner/controller"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type Env struct {
	db *sql.DB
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	interupt := make(chan os.Signal, 1)

	signal.Notify(interupt, os.Interrupt)
	go func() {
		log.Print("Server Started")
		router := mux.NewRouter()
		dbAccess := fmt.Sprintf("root:%s@/url", os.Getenv("DBPASS"))
		db, err := sql.Open("mysql", dbAccess)
		if err != nil {
			log.Fatalf("Db failed %s", err.Error())
		}
		env := &Env{
			db: db,
		}

		router.HandleFunc("/{shortUrl}", env.homeHandler).Methods(http.MethodGet)
		router.HandleFunc("/", env.storeHandler).Methods(http.MethodPost)
		router.HandleFunc("/", testRoute).Methods(http.MethodGet)
		server := &http.Server{
			Addr:    ":8080",
			Handler: router,
		}
		log.Fatal(server.ListenAndServe())
	}()
	<-interupt
	log.Print("Server Closed")

}

func testRoute(w http.ResponseWriter, r *http.Request) {
	// log.Print(time.Now())
	fmt.Fprint(w, time.Date(2020, time.June, 12, 0, 0, 0, 0, time.UTC).Format(time.RFC3339))
}

func (env *Env) storeHandler(w http.ResponseWriter, r *http.Request) {
	Name := r.FormValue("name")
	Url := r.FormValue("url")
	Date := r.FormValue("date")
	time.Now()
	log.Printf("%s %s", Name, Url)

	err := controller.SaveUrl(env.db, Name, Url, Date)
	if err != nil {
		// log.Fatal(err)
		fmt.Fprintf(w, err.Error())

	}
	fmt.Fprintf(w, "Saved")
}

func (env *Env) homeHandler(w http.ResponseWriter, r *http.Request) {

	url := mux.Vars(r)

	log.Print(url["shortUrl"])
	// var urls []controller.UrlPair
	urls := controller.FetchUrl(env.db, url["shortUrl"])
	log.Print(urls)
	http.Redirect(w, r, fmt.Sprintf("%s", urls), http.StatusFound)
}
