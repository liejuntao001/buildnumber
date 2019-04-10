package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var folder string
var mapping sync.Map

// Request wrap the build number
type Request struct {
	id string
	m  *sync.Mutex
}

type Response struct {
	Bn int `json:"bn"`
}

// IsValidUUID check if the string is UUID format
func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

func getResult(r *Request) Response {
	r.m.Lock()
	defer r.m.Unlock()

	filename := fmt.Sprintf("%s/%s", folder, r.id)
	res := Response{readAndUpdate(filename)}
	return res
}

func readAndUpdate(filename string) int {
	fd, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer fd.Close()

	var line int
	_, err = fmt.Fscanf(fd, "%d\n", &line)
	if err != nil {
		log.Println(err)
		if err == io.EOF {
			log.Println("Create new file", filename)
		} else {
			log.Printf("Reset content. Scan Failed %s\n", filename)
		}
	}
	fd.Seek(0, 0)
	line++
	fmt.Fprintf(fd, "%d\n", line)
	return line
}

// myHandler handle the request
func myHandler(w http.ResponseWriter, r *http.Request) {
	// debug concurrency
	//time.Sleep(time.Second * 5)
	params := mux.Vars(r)
	id := params["uuid"]

	if !IsValidUUID(id) {
		http.Error(w, "invalid UUID format", http.StatusBadRequest)
		return
	}

	a, ok := mapping.Load(id)
	if !ok {
		//log.Println("adding new bn for id", id)
		a = Request{id: id, m: &sync.Mutex{}}
		mapping.Store(id, a)
	}
	req := a.(Request)

	res := getResult(&req)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Printf("Error sending response %v", err)
	}
}

// logger log the access
func logger(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf(
			"%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			time.Since(start),
		)
	})
}

func checkFolder(folder string) {
	log.Println("checking storage folder", folder)

	if info, err := os.Stat(folder); os.IsNotExist(err) {
		log.Printf("Stroage folder %s doesn't exist, creating it\n", folder)
		os.MkdirAll(folder, 0777)
	} else if !info.IsDir() {
		panic(fmt.Sprintf(
			"Error! Storage folder %s exist but it's a file\nPlease delete it\n",
			folder))
	}
}

func main() {
	folder = os.Getenv("STORAGE_DIR")
	if folder = strings.TrimSpace(folder); len(folder) == 0 {
		folder = "data/"
	}
	checkFolder(folder)

	router := mux.NewRouter().StrictSlash(true)

	var handler http.Handler
	handler = http.HandlerFunc(myHandler)
	handler = logger(handler)
	router.Methods("POST").Path("/{uuid}").Handler(handler)
	log.Fatal(http.ListenAndServe(":8080", router))
}
