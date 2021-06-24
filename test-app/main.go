package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

var wg = sync.WaitGroup{}

type Form struct{
	Value string `json: "Value"` 
}

func main(){
	client := &http.Client{}
	
	wg.Add(5)
	//PUT request
	go func(){
		json, _ := json.Marshal(Form{
			Value: "MP",
		})
		req, _ := http.NewRequest(http.MethodPut, "http://localhost:8080/key/Indore", bytes.NewBuffer(json))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		resp, _ := client.Do(req)
		fmt.Println(resp.StatusCode)
		wg.Done()
	}()	
	
	//PUT request
	go func(){
		json, _ := json.Marshal(Form{
			Value: "Karnatka",
		})
		req, _ := http.NewRequest(http.MethodPut, "http://localhost:8080/key/Bengaluru", bytes.NewBuffer(json))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		resp, _ := client.Do(req)
		fmt.Println(resp.StatusCode)
		wg.Done()
	}()	

	//PUT request
	go func(){
		json, _ := json.Marshal(Form{
			Value: "Goa",
		})
		req, _ := http.NewRequest(http.MethodPut, "http://localhost:8080/key/Panaji", bytes.NewBuffer(json))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		resp, _ := client.Do(req)
		fmt.Println(resp.StatusCode)
		wg.Done()
	}()

	//GET request
	go func(){
		req, _ := http.NewRequest("GET", "http://localhost:8080/key/Panaji", nil)
		resp, _ := client.Do(req)
		body, _ := ioutil.ReadAll(resp.Body)
		sb := string(body)
   		fmt.Printf("%s \n", sb)
		wg.Done()
	}()


	//DELETE request
	go func(){
		req, _ := http.NewRequest("DELETE", "http://localhost:8080/key/Panaji", nil)
		resp, _ := client.Do(req)
		body, _ := ioutil.ReadAll(resp.Body)
		sb := string(body)
   		fmt.Printf("%s \n", sb)
		wg.Done()
	}()
	wg.Wait()
}

/*
for key: Panaji deleting
for key: Panaji value deleted: false
Iniating get for: Panaji
Value not found for: Panaji
for key: Bengaluru adding value:  Karnatka
for key: Bengaluru value added: Karnatka
for key: Indore adding value:  MP
for key: Indore value added: MP
for key: Panaji adding value:  Goa
for key: Panaji value added: Goa

****************Running Test again**************************

Iniating get for: Panaji
getting Panaji : Goa
for key: Panaji deleting
for key: Panaji value deleted: true
for key: Indore adding value:  MP
for key: Indore value added: MP
for key: Bengaluru adding value:  Karnatka
for key: Bengaluru value added: Karnatka
for key: Panaji adding value:  Goa
for key: Panaji value added: Goa

**************************************************************
Since for every request that request is completed before 
initiang the next req there test passed.
**************************************************************
*/
