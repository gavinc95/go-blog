package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gavinc95/go-blog/db"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type App struct {
	BlogStore db.BlogStore
	Addr      string
	Router    *mux.Router
}

func NewApp(addr string, idManager db.IDManager) *App {
	pg := MustDB()
	app := &App{
		BlogStore: db.NewBlogStore(pg, idManager),
		Addr:      addr,
		Router:    mux.NewRouter(),
	}

	app.Router.HandleFunc("/users", app.HandleGetUser).Methods("GET")
	app.Router.HandleFunc("/users", app.HandleCreateUser).Methods("POST")
	app.Router.HandleFunc("/users", app.HandleUpdateUser).Methods("PUT")
	app.Router.HandleFunc("/users", app.HandleDeleteUser).Methods("DELETE")

	app.Router.HandleFunc("/posts", app.HandleGetPost).Methods("GET")
	app.Router.HandleFunc("/posts/all", app.HandleGetAllPosts).Methods("GET")
	app.Router.HandleFunc("/posts", app.HandleCreatePost).Methods("POST")
	app.Router.HandleFunc("/posts", app.HandleUpdatePost).Methods("PUT")
	app.Router.HandleFunc("/posts", app.HandleDeletePost).Methods("DELETE")
	return app
}

const (
	usersTableCreationQuery = `CREATE TABLE IF NOT EXISTS users
	(
		id UUID NOT NULL,
		name varchar,
		email varchar,

		PRIMARY KEY (id),
		UNIQUE (email)
	)
	`

	postsTableCreationQuery = `CREATE TABLE IF NOT EXISTS posts
	(
		id UUID NOT NULL,
		user_id UUID NOT NULL, 
		title varchar NOT NULL,
	 	content TEXT,

		PRIMARY KEY (id),
		FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE
	);

	CREATE INDEX IF NOT EXISTS idx_user_id ON posts(user_id);
	`
)

func (a *App) ensureTablesExists() {
	log.Printf("creating Users table")
	if _, err := a.BlogStore.GetDB().Exec(usersTableCreationQuery); err != nil {
		log.Fatal(err)
	}

	log.Printf("creating Post table")
	if _, err := a.BlogStore.GetDB().Exec(postsTableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func (a *App) Run() {
	defer a.Close()

	// create the relevant DB tables
	a.ensureTablesExists()

	// start the HTTP server
	log.Printf("HTTP server listening on port: %s", a.Addr)
	log.Fatal(http.ListenAndServe(a.Addr, a.Router))
}

func (a *App) Close() error {
	if _, err := a.BlogStore.GetDB().Exec("DROP TABLE posts;"); err != nil {
		return err
	}

	if _, err := a.BlogStore.GetDB().Exec("DROP TABLE users;"); err != nil {
		return err
	}

	if err := a.BlogStore.GetDB().Close(); err != nil {
		return err
	}

	return nil
}

func MustDB() *sql.DB {
	user := getEnvWithDefault("POSTGRES_USER", "postgres")
	password := getEnvWithDefault("POSTGRES_PASSWORD", "password")
	dbname := getEnvWithDefault("APP_DB_NAME", "postgres")
	connectionString :=
		fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Panicf("failed to open postgres: %+v", err)
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
