package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

var (
	response   string
	port       = "8080"
	bufferSize = 100000
	logs       chan *http.Request
	disableCH  = false
)

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "%s", response)
	if !disableCH {
		logs <- r
	}
}

func load(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "%d", len(logs))
}

func main() {
	if portEnv, present := os.LookupEnv("PORT"); present {
		if v, err := strconv.Atoi(portEnv); err == nil {
			if v > 0 && v < 65536 {
				port = portEnv
			}
		}
	}

	if bufferSizeEnv, present := os.LookupEnv("BUFFER"); present {
		if v, err := strconv.Atoi(bufferSizeEnv); err == nil {
			if v > 0 {
				bufferSize = v
			}
		}
	}

	if v, present := os.LookupEnv("DISABLE_CH"); present && v == "true" {
		disableCH = true
	}

	logs = make(chan *http.Request, bufferSize)

	response = os.Getenv("RESPONSE")

	if !disableCH {
		go logger(logs)
	}

	router := httprouter.New()
	router.GET("/", index)
	router.GET("/load", load)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
