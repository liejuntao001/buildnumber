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
var mapping map[string]BN

// BN wrap the build number
type BN struct {
	id  string
	Seq int `json:"seq"`
	m   *sync.Mutex
}

// IsValidUUID check if the string is UUID format
func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

func getResult(bn *BN) BN {
	bn.m.Lock()
	defer bn.m.Unlock()

	filename := fmt.Sprintf("%s/%s", folder, bn.id)
	bn.Seq = readAndUpdate(filename)
	return *bn
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
		fmt.Println(err)
		if err == io.EOF {
			log.Println("Create new file", filename)
		} else {
			log.Printf("Reset content. Scan Failed %s: %v\n", filename, err)
		}
	}
	fd.Seek(0, 0)
	line++
	fmt.Fprintf(fd, "%d\n", line)
	return line
}

// GetBN handle the request
func GetBN(w http.ResponseWriter, r *http.Request) {
	// debug concurrency
	//time.Sleep(time.Second * 5)
	params := mux.Vars(r)
	id := params["uuid"]

	if !IsValidUUID(id) {
		http.Error(w, "invalid UUID format", http.StatusBadRequest)
		return
	}

	bn, ok := mapping[id]
	if !ok {
		//log.Println("adding new bn for id", id)
		bn = BN{id: id, m: &sync.Mutex{}}
		mapping[id] = bn
	}

	bn = getResult(&bn)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(bn); err != nil {
		panic(err)
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
		os.MkdirAll(folder, 0644)
	} else if !info.IsDir() {
		panic(fmt.Sprintf(
			"Error! Storage folder %s exist but it's a file\nPlease delete it\n",
			folder))
	}
}

func main() {
	folder = os.Getenv("STORAGE_DIR")
	mapping = make(map[string]BN)

	if folder = strings.TrimSpace(folder); len(folder) == 0 {
		folder = "data/"
	}
	checkFolder(folder)

	router := mux.NewRouter().StrictSlash(true)

	var handler http.Handler
	handler = http.HandlerFunc(GetBN)
	handler = logger(handler)
	router.Methods("POST").Path("/{uuid}").Handler(handler)
	log.Fatal(http.ListenAndServe(":8080", router))
}
