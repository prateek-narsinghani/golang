package main

import (
	// "database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db, _ = gorm.Open(sqlite.Open("first-app.db"), &gorm.Config{})

var m = sync.RWMutex{}

type Form struct {
	Key       string `gorm:"primaryKey"`
	Value     string `json: "Value"`
	Timestamp time.Time `gorm:"autoUpdateTime" json: "Timestamp"`
}

func addValue(w http.ResponseWriter, r *http.Request) {

	var form Form
	form.Key = mux.Vars(r)["key"]
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	
	json.Unmarshal(reqBody, &form)
	m.Lock()
	err=db.Create(&Form{Key:form.Key,
		Value:form.Value}).Error
	if err!=nil{
		db.Save(&Form{Key:form.Key,
			Value:form.Value})
	}
	m.Unlock()
	w.WriteHeader(http.StatusCreated)
}

func deleteValue(w http.ResponseWriter, r *http.Request) {
	var form Form
	form.Key = mux.Vars(r)["key"]
	
	m.Lock()
	
	db.Where("key = ?",form.Key).Delete(&Form{})
	m.Unlock()

}

func getValue(w http.ResponseWriter, r *http.Request) {
	var key string
	var form Form
	key = mux.Vars(r)["key"]
	m.RLock()
	err := db.Where("key=?",key).First(&form).Error
	if err != nil{
		fmt.Fprintf(w,"value not found")
	}else{
		fmt.Fprintf(w,"%s %s %s\n", form.Key, form.Value, form.Timestamp.String())
	}
	m.RUnlock()
}

func main() {
	db.AutoMigrate(&Form{})
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/key/{key}", getValue).Methods("GET")
	router.HandleFunc("/key/{key}", addValue).Methods("PUT")
	router.HandleFunc("/key/{key}", deleteValue).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))
}
