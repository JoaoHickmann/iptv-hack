package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
)

func redirect(w http.ResponseWriter, r *http.Request) {
	url, err := ioutil.ReadFile("/data/playlist-url")
	if err != nil {
		log.Print(err)
	}

	if string(url) == "" {
		url = []byte("https://www.google.com/")
	}

	log.Print(string(url))
	http.Redirect(w, r, string(url), 307)
}

func update(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Running...\n"))
	output, err := exec.Command("iptvgenerator").Output()
	if err != nil {
		log.Print(err)
		w.Write([]byte(err.Error()))
	} else {
		w.Write(output)
	}
}

func main() {
	http.HandleFunc("/", redirect)
	http.HandleFunc("/update", update)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
