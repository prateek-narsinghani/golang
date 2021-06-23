package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)


var MyMap = map[string]string{}

type Value struct{
	Value string `json: "Value"`
}

func addValue(w http.ResponseWriter, r *http.Request){
	var value Value
	key := mux.Vars(r)["key"]
	reqBody, err :=ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter details")
	}
	json.Unmarshal(reqBody, &value)
	MyMap[key] = value.Value
	w.WriteHeader(http.StatusCreated)
}

func deleteValue(w http.ResponseWriter, r *http.Request){
	key := mux.Vars(r)["key"]
	delete(MyMap, key)
}

func getValue(w http.ResponseWriter, r *http.Request){
	key:= mux.Vars(r)["key"]
	fmt.Fprintf(w, MyMap[key])
}

func main(){
	router:= mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/key/{key}", getValue).Methods("GET")
	router.HandleFunc("/key/{key}", addValue).Methods("PUT")
	router.HandleFunc("/key/{key}", deleteValue).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))
}

