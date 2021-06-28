package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

var db, _ = sql.Open("sqlite3", "./first-app.db")

var m = sync.RWMutex{}

type Form struct {
	Key       string `json: "Key"`
	Value     string `json: "Value"`
	Timestamp string `json: "Timestamp"`
}

func addValue(w http.ResponseWriter, r *http.Request) {

	var form Form
	form.Key = mux.Vars(r)["key"]
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	stmt, _ := db.Prepare("INSERT INTO myTable(key, value) VALUES (?,?)")
	json.Unmarshal(reqBody, &form)
	m.Lock()
	stmt.Exec(form.Key, form.Value)
	m.Unlock()
	w.WriteHeader(http.StatusCreated)
}

func deleteValue(w http.ResponseWriter, r *http.Request) {
	var form Form
	form.Key = mux.Vars(r)["key"]
	stmt, _ := db.Prepare("DELETE FROM myTable WHERE key = ? ")
	m.Lock()
	stmt.Exec(form.Key)
	m.Unlock()

}

func getValue(w http.ResponseWriter, r *http.Request) {
	var key, value, time string
	key = mux.Vars(r)["key"]
	m.RLock()
	row := db.QueryRow(`SELECT key, value, timestamp from myTable where key = ?`, key)
	err:= row.Scan(&key, &value, &time)
	if err == sql.ErrNoRows{
		fmt.Fprintf(w,"value not found")
	}else{
		fmt.Fprintf(w,"%s %s %s\n", key, value, time)
	}
	m.RUnlock()
}

func main() {
	stmt, _ := db.Prepare("CREATE TABLE IF NOT EXISTS myTable (key TEXT PRIMARY KEY, value TEXT, timestamp DATETIME DEFAULT CURRENT_TIMESTAMP)")
	stmt.Exec()
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/key/{key}", getValue).Methods("GET")
	router.HandleFunc("/key/{key}", addValue).Methods("PUT")
	router.HandleFunc("/key/{key}", deleteValue).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))
}
