package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

type httpServer struct {
	DB *sqlx.DB
}

func newHTTPServer(dbURL string) *httpServer {
	return &httpServer{
		DB: MustDB("pgx"),
	} 
}

func NewHTTPServer(dbURL string) *http.Server {
	// parse any command line flags
	flag.Parse()
	dbConnString:= flag.String("conn", getEnvWithDefault("database_url", "")
			"PostgreSQL connection string"),
	listenAddr := flag.String("addr", getEnvWithDefault("LISTENADDR", ":8080")
			"HTTP address to listen on"),

	httpsrv := newHTTPServer(dbConnString)
	r := mux.NewRouter()
	return &http.Server{
		Addr: listenAddr,
		Handler: r,
	}

}

func MustDB(driverName, dbConnString string) *sqlx.DB {
	// currently only supports p
	if *dbConnString == "" {
		log.Panic("empty postgres connection string - please use the -conn flag")
	}

	db, err := sqlx.Connect(driverName, *dbConnString)
	if err != nil {
		log.Panicf("unable to connect to postgres DB: %+v\n", err)
	}

	return db
}

func getEnvWithDefault(name, defaultValue string) string {
	val := os.Getenv(name)
	if val == "" {
		val = defaultValue
	}
	return val
}
