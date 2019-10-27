package todo

import (
	"database/sql"
	"fmt"
	"time"
)

//ToDo ...
type ToDo struct {
	ID      int64     `json:"id"`
	Task    string    `json:"task"`
	Created time.Time `json:"created_date,omitempty"`
	Updated time.Time `json:"updated_date,omitempty"`
}

func (t *ToDo) addToDo(db *sql.DB) error {
	time := time.Now().UTC()
	result, err := db.Exec("INSERT INTO todos(task, created_ts, updated_ts) VALUES($1, $2, $3)", t.Task, time, time)
	if err != nil {
		return err
	}

	t.ID, _ = result.LastInsertId()

	return nil
}

func (t *ToDo) getToDo(db *sql.DB) error {
	row := db.QueryRow("SELECT id, task, created_ts, updated_ts FROM todos WHERE id = $1", t.ID)

	if err := row.Scan(&t.ID, &t.Task, &t.Created, &t.Updated); err != nil {
		return err
	}

	return nil
}

func (t *ToDo) getToDos(db *sql.DB, start, count int) ([]ToDo, error) {
	statement := fmt.Sprintf("SELECT id, task, created_ts, updated_ts FROM todos LIMIT %d OFFSET %d", count, start)
	rows, err := db.Query(statement)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	todos := []ToDo{}

	for rows.Next() {
		var t ToDo
		if err := rows.Scan(&t.ID, &t.Task, &t.Created, &t.Updated); err != nil {
			return nil, err
		}
		todos = append(todos, t)
	}

	return todos, nil
}

func (t *ToDo) updateToDo(db *sql.DB) error {
	time := time.Now().UTC()

	result, err := db.Exec("UPDATE todos SET task=$1, updated_ts=$2 WHERE id=$3", t.Task, time, t.ID)
	if err != nil {
		return err
	}
	if count, _ := result.RowsAffected(); count == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (t *ToDo) deleteToDo(db *sql.DB) error {
	statement := fmt.Sprintf("DELETE FROM todos WHERE id=%d", t.ID)
	result, err := db.Exec(statement)
	if err != nil {
		return err
	}

	if count, _ := result.RowsAffected(); count == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (t *ToDo) createDB(db *sql.DB) error {

	const tableDropQuery = `DROP TABLE IF EXISTS todos`

	const tableCreationQuery = `
	CREATE TABLE IF NOT EXISTS todos
	(
	    id SERIAL PRIMARY KEY,
	    task VARCHAR(50) NOT NULL,
			created_ts timestamptz NOT NULL,
			updated_ts timestamptz NOT NULL
	)`

	db.Exec(tableDropQuery)

	_, err := db.Exec(tableCreationQuery)
	return err
}
