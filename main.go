package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)


var MyMap = map[string]string{}

var m = sync.RWMutex{} 

type Form struct{
	Key string `json: "Key"`
	Value string `json: "Value"` 
}

func addValue(w http.ResponseWriter, r *http.Request){
	
	var form Form

	form.Key = mux.Vars(r)["key"]
	reqBody, err :=ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	
	json.Unmarshal(reqBody, &form)
	
	m.Lock()
	fmt.Printf("for key: %s adding value:  %s\n", form.Key, form.Value)
	MyMap[form.Key] = form.Value
	fmt.Printf("for key: %s value added: %s\n", form.Key, form.Value)
	m.Unlock()
	
	w.WriteHeader(http.StatusCreated)
}

func deleteValue(w http.ResponseWriter, r *http.Request){
	var form Form
	form.Key = mux.Vars(r)["key"]

	m.Lock()
	
	fmt.Printf("for key: %s deleting\n", form.Key)
	
	_, ok:=MyMap[form.Key]
	if ok{
		delete(MyMap, form.Key)
	}
	fmt.Printf("for key: %s value deleted: %t\n", form.Key, ok)
	
	m.Unlock()

}

func getValue(w http.ResponseWriter, r *http.Request){
	m.RLock()
	key:= mux.Vars(r)["key"]
	fmt.Printf("Iniating get for: %s\n", key)
	if val,ok:=MyMap[key]; ok {
		
		fmt.Printf("getting %s : %s\n", key, val)
	
	}else{
		fmt.Printf("Value not found for: %s \n", key)
	}
	m.RUnlock()
	
}


func main(){
	router:= mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/key/{key}", getValue).Methods("GET")
	router.HandleFunc("/key/{key}", addValue).Methods("PUT")
	router.HandleFunc("/key/{key}", deleteValue).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))
}
