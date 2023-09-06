package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
)

var counter int
var lock = &sync.Mutex{}

func incrementCounter(w http.ResponseWriter, r *http.Request) {
	lock.Lock()
	counter++
	fmt.Fprintf(w, "Your current increment: %v", strconv.Itoa(counter))
	lock.Unlock()
}

func echoString(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "'ello, govnah")
}

func MethodHandler(w http.ResponseWriter, r *http.Request) {
	// if r.URL.Path != "/" {
	// 	http.Error(w, "404 Not Found", http.StatusNotFound)
	// 	return
	// }

	objects := Objects{
		Object{ID: 1, Name: "Object 1", Description: "Obj 1 desc", Active: true},
		Object{ID: 2, Name: "Object 2", Description: "Obj 2 desc", Active: false},
		Object{ID: 3, Name: "Object 3", Description: "Obj 3 desc", Active: true},
		Object{ID: 4, Name: "Object 4", Description: "Obj 4 desc", Active: false},
	}

	switch r.Method {
	case "GET":
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(objects)
	case "POST":
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(objects); err != nil {
			panic(err)
		}
		json.NewEncoder(w).Encode(objects)
	default:
		fmt.Fprintf(w, "The request used a non-GET or POST method")
	}
}

func main() {

	http.HandleFunc("/", echoString)

	http.HandleFunc("/increment", incrementCounter)

	http.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hi there")
	})

	http.HandleFunc("/methods", MethodHandler)

	log.Fatal(http.ListenAndServe(":8081", nil))
}

type Object struct {
	ID          int    `json: "id"`
	Name        string `json: "name"`
	Description string `json: "description"`
	Active      bool   `json: "active"`
}

type Objects []Object
