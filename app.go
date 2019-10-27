// app.go

package todo

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

var config = "dev-config"

// App ...
type App struct {
	router *mux.Router
	db     *sql.DB
	server *http.Server
	stop   chan os.Signal
	done   chan bool
}

// Route ...
type Route struct {
	name        string
	method      string
	pattern     string
	handlerFunc http.HandlerFunc
}

// Initialize ...
func (a *App) Initialize() {

	viper.SetConfigType("json")
	viper.SetConfigName(config)
	viper.AddConfigPath(".")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		log.Fatal(err)
	}

	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)

	connectionString := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=%s connect_timeout=%d",
		viper.GetString("db-server.host"),
		viper.GetString("db-server.port"),
		viper.GetString("db-server.user-id"),
		viper.GetString("db-server.secret"),
		viper.GetString("db-server.db-name"),
		viper.GetString("db-server.ssl-mode"),
		viper.GetInt("db-server.connect_timeout"),
	)

	log.Println("Connecting to database...")
	a.db, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Creating db...")
	err = a.createDB(a.db)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Creating router...")
	a.router = a.NewRouter()

	log.Println("Creating http server...")
	a.server = &http.Server{
		Handler:      a.router,
		Addr:         viper.GetString("http-server.address"),
		ReadTimeout:  viper.GetDuration("http-server.read-timeout") * time.Second,
		WriteTimeout: viper.GetDuration("http-server.write-timeout") * time.Second,
	}
}

// Run ...
func (a *App) Run() {
	// Start Server
	go func() {
		if err := a.server.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()
}

// WaitOnShutdown ...
func (a *App) WaitOnShutdown() {
	// Setting up signal capturing
	a.stop = make(chan os.Signal, 1)
	a.done = make(chan bool, 1)

	//signal.Notify(a.stop, os.Interrupt)
	signal.Notify(a.stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Waiting for signals
	<-a.stop

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := a.server.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}

func (a *App) createDB(db *sql.DB) error {

	const tableDropQuery = `DROP TABLE IF EXISTS todos`

	const tableCreationQuery = `
	CREATE TABLE IF NOT EXISTS todos
	(
	    id SERIAL PRIMARY KEY,
	    task VARCHAR(50) NOT NULL,
			completed boolean NULL,
			created_ts timestamptz NOT NULL,
			updated_ts timestamptz NOT NULL
	)`

	db.Exec(tableDropQuery)

	_, err := db.Exec(tableCreationQuery)
	return err
}

//NewRouter ...
func (a *App) NewRouter() *mux.Router {
	api := API{}
	api.db = a.db

	//Routes ...
	routes := []Route{

		Route{
			"addToDo",
			strings.ToUpper("Post"),
			"/v1/todo/{todoId}",
			api.addToDo,
		},

		Route{
			"deleteToDo",
			strings.ToUpper("Delete"),
			"/v1/todo/{todoId}",
			api.deleteToDo,
		},

		Route{
			"getToDos",
			strings.ToUpper("Get"),
			"/v1/todos",
			api.getToDos,
		},

		Route{
			"getToDoByID",
			strings.ToUpper("Get"),
			"/v1/todo/{todoId}",
			api.getToDoByID,
		},

		Route{
			"updateToDo",
			strings.ToUpper("Put"),
			"/v1/todo/{todoId}",
			api.updateToDo,
		},
	}

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.handlerFunc
		handler = Logger(handler, route.name)

		router.
			Methods(route.method).
			Path(route.pattern).
			Name(route.name).
			Handler(handler)
	}

	return router
}
