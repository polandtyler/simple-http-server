package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
)

// you can declare a series of `var` and `const` in this manner
var (
	counter int
	lock    = &sync.Mutex{}
)

func incrementCounter(w http.ResponseWriter, _ *http.Request) { // can use blank identifier in params to note an argument is required by not used
	lock.Lock()
	// defer is a hand construct... here is some info https://gobyexample.com/defer
	defer lock.Unlock()
	counter++
	// use blank identifier to avoid lint error of unused error
	_, _ = fmt.Fprintf(w, "Your current increment: %v", strconv.Itoa(counter))
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

	// Just a note: the below is good to understand in terms of stdlin net/http
	// you'll "graduate" to a proper router, like chi, which makes this stuff a lot cleaner.
	switch r.Method {
	// net/http has constants for the methods
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(objects) // Great
	// net/http has constants for the methods
	case http.MethodPost:
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(objects); err != nil {
			// calling panic in your code is very rare, the only real times I call panic directly is during
			// application startup (e.g., if config is missing then panic because there is no point)
			// panic(err)

			// instead...
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		//_ = json.NewEncoder(w).Encode(objects)
	default:
		// alternative to fmt.Fprintf...
		// w.Write([]byte("The request used a non-GET or POST method"))
		// prefer w.Write to fmt.Fprintf
		fmt.Fprintf(w, "The request used a non-GET or POST method")
	}
}

func main() {

	http.HandleFunc("/", echoString)

	http.HandleFunc("/increment", incrementCounter)

	http.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
		// prefer w.Write to fmt.Fprintf
		fmt.Fprintf(w, "hi there")
	})

	http.HandleFunc("/methods", MethodHandler)

	// fatal is a log messsage which calls panic.
	// in the code below, even if the server exits without error, you will see a panic in your logs.
	//log.Fatal(http.ListenAndServe(":8081", nil))

	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal(err)
	}
}

type Object struct {
	// no space between colon and "
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Active      bool   `json:"active"`
}

type Objects []Object
