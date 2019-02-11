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

func redirect(w http.ResponseWriter, r *http.Request) {
	url, err := ioutil.ReadFile(playlistFile)
	if err != nil {
		log.Print(err)
		w.WriteHeader(500)
		w.Write([]byte("Error: " + err.Error()))
		return
	}

	if string(url) == "" {
		err = errors.New("Playlist URL not found!\nRun /update")
		log.Print(err)
		w.WriteHeader(404)
		w.Write([]byte("Error: " + err.Error()))
		return
	}

	log.Print(string(url))
	http.Redirect(w, r, string(url), 307)
}

func update(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Running...\n"))
	output, err := exec.Command("iptvgenerator").Output()
	if err != nil {
		log.Print(err)
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
	} else {
		log.SetPrefix("")
		log.Print(output)
		log.SetPrefix("Proxy")
		w.Write(output)
	}
}

func main() {
	log.SetPrefix("PROXY: ")
	log.SetOutput(os.Stdout)

	http.HandleFunc("/", redirect)
	http.HandleFunc("/update", update)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Panic(err)
	}
}
