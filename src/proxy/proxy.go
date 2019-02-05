package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func redirect(w http.ResponseWriter, r *http.Request) {
	filedir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Print(err)
	}
	filename := "playlist-url"

	url, err := ioutil.ReadFile(filepath.Join(filedir, filename))
	if err != nil {
		log.Print(err)
	}

	if string(url) == "" {
		url = []byte("https://www.google.com/")
	}

	log.Print(string(url))
	http.Redirect(w, r, string(url), 307)
}

func main() {
	http.HandleFunc("/", redirect)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
