package main

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
)

const (
	playlistFile = "/data/playlist-url"
)

func handleError(w *http.ResponseWriter, statusCode int, err error) {
	log.Print(err)
	(*w).WriteHeader(statusCode)
	(*w).Write([]byte("Error: " + err.Error()))
}

func redirect(w http.ResponseWriter, r *http.Request) {
	url, err := ioutil.ReadFile(playlistFile)
	if err != nil {
		handleError(&w, 500, err)
		return
	}

	if string(url) == "" {
		err = errors.New("Playlist URL not found!\nRun /update")
		handleError(&w, 404, err)
		return
	}

	log.Print(string(url))
	http.Redirect(w, r, string(url), 307)
}

func update(w http.ResponseWriter, r *http.Request) {
	output, err := exec.Command("iptvgenerator").Output()
	if err != nil {
		handleError(&w, 500, err)
		return
	}

	log.Print(string(output))
	w.Write(output)
}

func main() {
	log.SetOutput(os.Stdout)
	log.SetPrefix("PROXY: ")

	http.HandleFunc("/", redirect)
	http.HandleFunc("/update", update)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Panic(err)
	}
}
