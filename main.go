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
// var wg = sync.WaitGroup{}

type Value struct{
	Value string `json: "Value"`
}

func addValue(w http.ResponseWriter, r *http.Request){
	m.Lock()
	var value Value
	key := mux.Vars(r)["key"]
	reqBody, err :=ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter details")
	}
	json.Unmarshal(reqBody, &value)
	MyMap[key] = value.Value
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, "added value")
	m.Unlock()
}

func deleteValue(w http.ResponseWriter, r *http.Request){
	m.Lock()
	key := mux.Vars(r)["key"]
	delete(MyMap, key)
	fmt.Fprint(w, "deleted value")
	m.Unlock()
}

func getValue(w http.ResponseWriter, r *http.Request){
	m.RLock()
	key:= mux.Vars(r)["key"]
	if val,ok:=MyMap[key];ok{
		fmt.Fprint(w,val)
	}else{
		fmt.Fprint(w,"value not found")
	}
	m.RUnlock()
}

// func test(w http.ResponseWriter, r *http.Request){
// // 	wg.Add(2)
// // 	//sending put req for {go:lang}
// 	go func(){
// 		client := &http.Client{}
// 		json,err := json.Marshal(Value{Value: "lang"})
// 		if err!=nil{
// 			panic(err)
// 		}
// 		req,err := http.NewRequest(http.MethodPut,"http://localhost:8080/key/go", bytes.NewBuffer(json))
// 		if err != nil {
// 			panic(err)
// 		}
// 		req.Header.Set("Content-Type", "application/json; charset=utf-8")
// 		resp, err := client.Do(req)
// 		if err != nil {
// 			panic(err)
// 		}
// 		// fmt.Println(resp.StatusCode)
// 		body, _ := ioutil.ReadAll(resp.Body)
// 		sb := string(body)
// 		fmt.Fprint(w, sb)
// 		// wg.Done()
// 	}()

// 	//sending put req for {go:language}
// 	go func(){
// 		client := &http.Client{}
// 		json,err := json.Marshal(Value{Value: "language"})
// 		if err!=nil{
// 			panic(err)
// 		}
// 		req,err := http.NewRequest(http.MethodPut,"http://localhost:8080/key/go", bytes.NewBuffer(json))
// 		if err != nil {
// 			panic(err)
// 		}
// 		req.Header.Set("Content-Type", "application/json; charset=utf-8")
// 		resp , err := client.Do(req)
// 		if err != nil {
// 			panic(err)
// 		}
// 		fmt.Println(resp.StatusCode)
// 		wg.Done()
// 	}()
// 	// //sending get req for go
// 	// go func(){
// 	// 	resp, err := http.Get("http://localhost:8080/key/go")
// 	// 	if err != nil {
// 	// 		panic(err)
// 	// 	}
// 	// 	body, err := ioutil.ReadAll(resp.Body)
// 	// 	if err != nil {
// 	// 		panic(err)
// 	// 	}
// 	// 	sb := string(body)
// 	// 	fmt.Fprint(w,sb)
// 	// 	}()
// }

func main(){
	router:= mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/key/{key}", getValue).Methods("GET")
	router.HandleFunc("/key/{key}", addValue).Methods("PUT")
	router.HandleFunc("/key/{key}", deleteValue).Methods("DELETE")
	// router.HandleFunc("/test", test).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))
	// fmt.Print("hi there")
	// wg.Wait()
}

