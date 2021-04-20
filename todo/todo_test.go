package todo

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

var a App

func TestMain(m *testing.M) {
	a.Initialize(Host, Port, User, Password, Dbname)
	ensureTableExists()
	code := m.Run()
	clearTable()
	os.Exit(code)
}

func ensureTableExists() {
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	a.DB.Exec("DELETE FROM tasks")
	a.DB.Exec("ALTER SEQUENCE tasks_id_seq RESTART WITH 1")
}

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS tasks
(
    Id	SERIAL PRIMARY KEY,
    Description    varchar(40),
	Completed int
)`


func TestEmptyTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/tasks", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}


func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}


func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}


func TestGetNonExistentTask(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/task/11", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Task not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Task not found'. Got '%s'", m["error"])
	}
}


func TestCreateTask(t *testing.T) {

	clearTable()

	var jsonStr = []byte(`{"description":"testowy task dla testu", "completed": 0}`)
	req, _ := http.NewRequest("POST", "/task", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["description"] != "testowy task dla testu" {
		t.Errorf("Expected task name to be 'test task'. Got '%v'", m["description"])
	}

	if m["completed"] != 0.0 {
		t.Errorf("Expected task price to be '0'. Got '%v'", m["completed"])
	}

	// the id is compared to 1.0 because JSON unmarshaling converts numbers to
	// floats, when the target is a map[string]interface{}
	if m["id"] != 1.0 {
		t.Errorf("Expected task ID to be '1'. Got '%v'", m["id"])
	}
}


func TestGetTask(t *testing.T) {
	clearTable()
	addTasks(1)

	req, _ := http.NewRequest("GET", "/task/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func addTasks(count int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		a.DB.Exec("INSERT INTO tasks(DESCRIPTION, COMPLETED) VALUES($1, $2)", "Task "+ strconv.Itoa(i), 2)
	}
}


func TestUpdateTask(t *testing.T) {

	clearTable()
	addTasks(1)

	req, _ := http.NewRequest("GET", "/task/1", nil)
	response := executeRequest(req)
	var originalTask map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalTask)

	var jsonStr = []byte(`{"description":"test task2", "completed": 1}`)
	req, _ = http.NewRequest("PUT", "/task/1", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)


	if m["description"] != "test task2" {
		t.Errorf("Expected task name to be 'test task2'. Got '%v'", m["description"])
	}

	// the id is compared to 1.0 because JSON unmarshaling converts numbers to
	// floats, when the target is a map[string]interface{}
	if m["completed"] != 1.0 {
			t.Errorf("Expected completed to be '1'. Got '%v'", m["completed"])
	}
}


func TestDeleteTask(t *testing.T) {
	//clearTable()
	addTasks(1)

	req, _ := http.NewRequest("GET", "/task/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/task/1", nil)
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/task/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}