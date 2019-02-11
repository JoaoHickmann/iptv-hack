package main

import (
	"errors"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	signUpURL = "https://server.azsat.org/novo/usuario/registrar_iptv.php"
	dashURL   = "https://server.azsat.org/"

	fg = "2e72627a07cac7284cf3bfa072b76a9b"

	playlistFile = "/data/playlist-url"
)

func checkError(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func createUsername() (username string) {
	rand.Seed(time.Now().UnixNano())

	for index := 0; index < 16; index++ {
		username += string('a' + rand.Int31n('z'-'a'))
	}
	return
}

func createPostData(csrf, user string) (postData url.Values) {
	postData = url.Values{
		"step":        {"2"},
		"fg":          {fg},
		"csrf":        {csrf},
		"login":       {user},
		"email":       {user + "@gmail.com"},
		"senha":       {user},
		"senha2":      {user},
		"operadora[]": {"IPTV", "IPTV", "IPTV"},
	}
	return
}

func getCSRF() (csrf string, err error) {
	var response *http.Response
	response, err = http.Get(signUpURL)
	if err != nil {
		return
	}
	defer response.Body.Close()

	var document *goquery.Document
	document, err = goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return
	}

	csrf = document.Find("input[name=\"csrf\"]").Nodes[0].Attr[2].Val
	return
}

func signUpAndGetURL() (playlistURL string, err error) {
	var csrf string
	csrf, err = getCSRF()
	if err != nil {
		return
	}

	user := createUsername()
	postData := createPostData(csrf, user)

	var response *http.Response
	response, err = http.PostForm(signUpURL, postData)
	if err != nil {
		return
	}
	defer response.Body.Close()

	var bodyBytes []byte
	bodyBytes, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	if !strings.Contains(string(bodyBytes), "ALERTA DE SUCESSO") {
		err = errors.New(string(bodyBytes))
		return
	}

	playlistURL, err = getPlaylistURL()
	return
}

func getPlaylistURL() (playlistURL string, err error) {
	var response *http.Response
	response, err = http.Get(dashURL)
	if err != nil {
		return
	}
	defer response.Body.Close()

	response, err = http.Get(dashURL)
	if err != nil {
		return
	}
	defer response.Body.Close()

	var document *goquery.Document
	document, err = goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return
	}

	playlistURL = document.Find("input[name=\"url\"]").Nodes[0].Attr[3].Val
	return
}

func main() {
	log.SetOutput(os.Stdout)
	http.DefaultClient.Jar, _ = cookiejar.New(nil)

	playlistURL, err := signUpAndGetURL()
	checkError(err)

	log.Print(playlistURL)

	err = ioutil.WriteFile(playlistFile, []byte(playlistURL), 0644)
	checkError(err)
}
