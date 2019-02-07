package main

import (
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
)

func checkError(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func createLogin() string {
	rand.Seed(time.Now().UnixNano())
	login := ""
	for index := 0; index < 16; index++ {
		login += string('a' + rand.Int31n('z'-'a'))
	}
	return login
}

func createPostData(csrf, user string) url.Values {
	data := url.Values{
		"step":        {"2"},
		"fg":          {fg},
		"csrf":        {csrf},
		"login":       {user},
		"email":       {user + "@gmail.com"},
		"senha":       {user},
		"senha2":      {user},
		"operadora[]": {"IPTV", "IPTV", "IPTV"},
	}
	return data
}

func getCSRF() string {
	response, err := http.Get(signUpURL)
	checkError(err)
	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)
	checkError(err)

	return document.Find("input[name=\"csrf\"]").Nodes[0].Attr[2].Val
}

func signUpNewUser() bool {
	csrf := getCSRF()

	user := createLogin()
	data := createPostData(csrf, user)

	response, err := http.PostForm(signUpURL, data)
	checkError(err)
	defer response.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(response.Body)

	return strings.Contains(string(bodyBytes), "ALERTA DE SUCESSO")
}

func getPlaylistURL() string {
	response, err := http.Get(dashURL)
	checkError(err)
	defer response.Body.Close()

	response, err = http.Get(dashURL)
	checkError(err)
	defer response.Body.Close()

	doc, err := goquery.NewDocumentFromReader(response.Body)
	checkError(err)

	return doc.Find("input[name=\"url\"]").Nodes[0].Attr[3].Val
}

func main() {
	log.SetOutput(os.Stdout)
	http.DefaultClient.Jar, _ = cookiejar.New(nil)

	if signUpNewUser() {
		playlistURL := getPlaylistURL()
		log.Print(playlistURL)

		err := ioutil.WriteFile("/data/playlist-url", []byte(playlistURL), 0644)
		checkError(err)
	} else {
		log.Print("FAILED!")
	}
}
