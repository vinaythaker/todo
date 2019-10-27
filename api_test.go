package todo

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

var a App
var jwtToken *http.Cookie

func TestInit(t *testing.T) {
	a = App{}
	a.Initialize()
	//a.Run()
}

func TestAPIMalformedRequestBodyAddToDo(t *testing.T) {

	payload := []byte(``)
	req, _ := http.NewRequest("POST", "/v1/todo/10", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Invalid request payload" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Invalid request payload'. Got '%s'", m["error"])
	}
}

func TestAPIMalformedRequestURLAddToDo(t *testing.T) {

	payload := []byte(`{"id":10,"name":"grrr"}`)
	req, _ := http.NewRequest("POST", "/v1/todo/8&6%434", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Invalid request" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Invalid todo ID'. Got '%s'", m["error"])
	}
}

func TestAPIAddToDo(t *testing.T) {

	payload := []byte(`{"id":1,"task":"do work"}`)
	req, _ := http.NewRequest("POST", "/v1/todo/1", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	var todo ToDo
	json.Unmarshal(response.Body.Bytes(), &todo)

	if todo.Task != "do work" {
		t.Errorf("Expected user name to be 'do work'. Got '%v'", todo.Task)
	}
}

func TestAPIMalformedRequestURLGetToDo(t *testing.T) {

	req, _ := http.NewRequest("GET", "/v1/todo/8&6%434", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Invalid request" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Invalid todo ID'. Got '%s'", m["error"])
	}
}

func TestAPIGetNonExistentToDo(t *testing.T) {
	req, _ := http.NewRequest("GET", "/v1/todo/45", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "ToDo not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'ToDo not found'. Got '%s'", m["error"])
	}
}

func TestAPIGetToDo(t *testing.T) {
	req, _ := http.NewRequest("GET", "/v1/todo/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestAPIGetToDos(t *testing.T) {
	req, _ := http.NewRequest("GET", "/v1/todos", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	var todos []ToDo
	json.Unmarshal(response.Body.Bytes(), &todos)
	if len(todos) != 1 {
		t.Errorf("Expected 1 todo got %d", len(todos))
	}
}

func TestAPIMalformedRequestBodyUpdateToDo(t *testing.T) {

	payload := []byte(``)
	req, _ := http.NewRequest("PUT", "/v1/todo/1", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Invalid request payload" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Invalid request payload'. Got '%s'", m["error"])
	}
}

func TestAPIMalformedRequestURLUpdateToDo(t *testing.T) {

	payload := []byte(`{"id":1,"task":"do work"}`)
	req, _ := http.NewRequest("PUT", "/v1/todo/%21%40#23%24%25%5E1234567890abcdefghijklmnopqrstuvwxyz", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Invalid request" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Invalid todo ID'. Got '%s'", m["error"])
	}
}

func TestAPIUpdateNonExistentToDo(t *testing.T) {

	payload := []byte(`{"id":45,"task":"do work"}`)

	req, _ := http.NewRequest("PUT", "/v1/todo/45", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "ToDo not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'ToDo not found'. Got '%s'", m["error"])
	}
}

func TestAPIUpdateToDo(t *testing.T) {

	payload := []byte(`{"id":1,"task":"do work"}`)
	req, _ := http.NewRequest("POST", "/v1/todo/1", bytes.NewBuffer(payload))
	response := executeRequest(req)

	payload = []byte(`{"id":1,"task":"do more work"}`)
	req, _ = http.NewRequest("PUT", "/v1/todo/1", bytes.NewBuffer(payload))
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var todo ToDo
	json.Unmarshal(response.Body.Bytes(), &todo)

	if todo.Task != "do more work" {
		t.Errorf("Expected the task to change from '%v' to '%v'. Got '%v'", "do work", "do more work", todo.Task)
	}
}

func TestAPIMalformedRequestURLDeleteToDo(t *testing.T) {

	req, _ := http.NewRequest("DELETE", "/v1/todo/%21%40#23%24%25%5E1234567890abcdefghijklmnopqrstuvwxyz", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Invalid request" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Invalid todo ID'. Got '%s'", m["error"])
	}
}

func TestAPIDeleteNonExistentToDo(t *testing.T) {

	req, _ := http.NewRequest("DELETE", "/v1/todo/45", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "ToDo not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'ToDo not found'. Got '%s'", m["error"])
	}
}

func TestAPIDeleteToDo(t *testing.T) {

	req, _ := http.NewRequest("DELETE", "/v1/todo/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/v1/todo/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	if jwtToken != nil {
		req.AddCookie(jwtToken)
	}
	a.router.ServeHTTP(rr, req)
	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}
