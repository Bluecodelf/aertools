package main

import (
	"database/sql"
	"github.com/Bluecodelf/rets"
	"log"
	"net/http"
	"regexp"
)

var db *sql.DB
var cfg *Configuration

func CheckError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	var err error
	log.Println("Reading configuration...")
	// REVIEW: config path is hardcoded... this is probably bad practice.
	cfg, err = ReadConfiguration("/etc/aertools/config.json")
	CheckError(err)

	log.Println("Connecting to database...")
	db, err = rets.OpenDatabase(cfg.SQL)
	defer db.Close()
	CheckError(err)

	router := rets.NewRouter()

	// Connect
	router.AddRoute(rets.Route{"GET", regexp.MustCompile("^\\/connect$"),
		HandlerGETConnect})

	// Lockers
	router.AddRoute(rets.Route{"POST", regexp.MustCompile("^\\/lockers$"),
		HandlerPOSTLockers})
	router.AddRoute(rets.Route{"PUT", regexp.MustCompile("^\\/locker/[0-9]+$"),
		HandlerPUTLocker})
	router.AddRoute(rets.Route{"GET", regexp.MustCompile("^\\/lockers$"),
		HandlerGETLockers})

	log.Println("AERTools Server is up and running. Thanks for all the fish!")
	http.ListenAndServe(":"+cfg.Port, router)
}
