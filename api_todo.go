package todo

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// API ...
type API struct {
	db *sql.DB
}

func (a *API) getToDos(w http.ResponseWriter, r *http.Request) {

	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count > 10 || count < 1 {
		count = 10
	}

	if start < 0 {
		start = 0
	}

	t := ToDo{}
	todos, _ := t.getToDos(a.db, start, count)

	respondWithJSON(w, http.StatusOK, todos)
}

func (a *API) getToDoByID(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["todoId"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	t := ToDo{ID: int64(id)}
	err = t.getToDo(a.db)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "ToDo not found")
			return
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	respondWithJSON(w, http.StatusOK, t)
}

func (a *API) addToDo(w http.ResponseWriter, r *http.Request) {
	var t ToDo

	vars := mux.Vars(r)
	_, err := strconv.Atoi(vars["todoId"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&t); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	defer r.Body.Close()

	if err := t.addToDo(a.db); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, t)
	return

}

func (a *API) deleteToDo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["todoId"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	t := ToDo{ID: int64(id)}
	if err := t.deleteToDo(a.db); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "ToDo not found")
			return
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	respondWithJSON(w, http.StatusOK, nil)
	return
}

func (a *API) updateToDo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	_, err := strconv.Atoi(vars["todoId"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	var t ToDo
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&t); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := t.updateToDo(a.db); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "ToDo not found")
			return
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	respondWithJSON(w, http.StatusOK, t)
	return
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code)
	w.Write(response)
}
