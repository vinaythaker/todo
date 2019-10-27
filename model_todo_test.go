package todo

import (
	"database/sql"
	"fmt"
	"log"
	"testing"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

var db *sql.DB

func TestMainModel(t *testing.T) {
	var config = "dev-config"

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

	db, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
}

func TestCreateDB(t *testing.T) {
	todo := ToDo{}
	err := todo.createDB(db)
	if err != nil {
		log.Println(err)
	}
}

func TestAddToDo(t *testing.T) {
	todo := ToDo{
		Task: "todo-1",
	}
	err := todo.addToDo(db)
	if err != nil {
		log.Println(err)
	}
	if todo.ID != 0 {
		log.Printf("Expected todo id to be 0 instead got %d", todo.ID)
	}
}

func TestAddErrorToDo(t *testing.T) {
	todo := ToDo{
		Task: "123456789012345678901234567890123456789012345678901234567890",
	}

	err := todo.addToDo(db)
	if err == nil {
		log.Println("Expected error...but didn't")
	}
}

func TestGetToDo(t *testing.T) {
	todo := ToDo{
		ID: 1,
	}
	err := todo.getToDo(db)

	if err != nil {
		log.Println(err)
	}

	if todo.Task != "todo-1" {
		log.Printf("Expected todo task to be 'todo-1' instead got %s", todo.Task)

	}
}

func TestGetToDos(t *testing.T) {
	todo := ToDo{}
	var todos []ToDo
	var err error
	todos, err = todo.getToDos(db, 0, 10)

	if err != nil {
		log.Println(err)
	}

	if todos[0].Task != "todo-1" {
		log.Printf("Expected todo task to be 'todo-1' instead got %s", todo.Task)

	}
}

func TestUpdateToDo(t *testing.T) {
	todo := ToDo{
		ID:   1,
		Task: "todo-2",
	}

	err := todo.updateToDo(db)
	if err != nil {
		log.Println(err)
	}

	if todo.Task != "todo-2" {
		log.Printf("Expected todo task to be 'todo-2' instead got %s", todo.Task)

	}
}

func TestUpdateNonExistentToDo(t *testing.T) {
	todo := ToDo{
		ID: 2,
	}

	err := todo.updateToDo(db)
	if err != sql.ErrNoRows {
		log.Println(err)
	}
}

func TestUpdateErrorToDo(t *testing.T) {
	todo := ToDo{
		ID:   1,
		Task: "123456789012345678901234567890123456789012345678901234567890",
	}

	err := todo.updateToDo(db)
	if err == nil {
		log.Println("Expected error...but didn't")
	}
}

func TestDeleteToDo(t *testing.T) {
	todo := ToDo{
		ID: 1,
	}

	err := todo.deleteToDo(db)
	if err != nil {
		log.Println(err)
	}

	err = todo.getToDo(db)
	if err != sql.ErrNoRows {
		log.Println(err)
	}
}

func TestDeleteNonExistentToDo(t *testing.T) {
	todo := ToDo{
		ID: 1,
	}

	err := todo.deleteToDo(db)
	if err != sql.ErrNoRows {
		log.Println(err)
	}
}
