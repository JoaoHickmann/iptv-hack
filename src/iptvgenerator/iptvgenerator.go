package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

const (
	signUpURL = "https://server.azsat.org/novo/usuario/registrar_iptv.php"
	dashURL   = "https://server.azsat.org/"

	info = "dXNlcl9hZ2VudD0hPU1vemlsbGEvNS4wIChYMTE7IExpbnV4IHg4Nl82NCkgQXBwbGVXZWJLaXQvNTM3LjM2IChLSFRNTCwgbGlrZSBHZWNrbykgQ2hyb21lLzcyLjAuMzYyNi45NiBTYWZhcmkvNTM3LjM2QEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAfHx8fHx8fHx8fHx8fHxsYW5ndWFnZT0hPXB0LUJSQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAfHx8fHx8fHx8fHx8fHxjb2xvcl9kZXB0aD0hPTI0QEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAfHx8fHx8fHx8fHx8fHxwaXhlbF9yYXRpbz0hPTFAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEB8fHx8fHx8fHx8fHx8fGhhcmR3YXJlX2NvbmN1cnJlbmN5PSE9NEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQHx8fHx8fHx8fHx8fHx8cmVzb2x1dGlvbj0hPTE5MjAsMTA4MEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQHx8fHx8fHx8fHx8fHx8YXZhaWxhYmxlX3Jlc29sdXRpb249IT0xOTIwLDEwNDBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEB8fHx8fHx8fHx8fHx8fHRpbWV6b25lX29mZnNldD0hPTEyMEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQHx8fHx8fHx8fHx8fHx8c2Vzc2lvbl9zdG9yYWdlPSE9MUBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQHx8fHx8fHx8fHx8fHx8bG9jYWxfc3RvcmFnZT0hPTFAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEB8fHx8fHx8fHx8fHx8fGluZGV4ZWRfZGI9IT0xQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAfHx8fHx8fHx8fHx8fHxvcGVuX2RhdGFiYXNlPSE9MUBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQHx8fHx8fHx8fHx8fHx8Y3B1X2NsYXNzPSE9dW5rbm93bkBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQHx8fHx8fHx8fHx8fHx8bmF2aWdhdG9yX3BsYXRmb3JtPSE9TGludXggeDg2XzY0QEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAfHx8fHx8fHx8fHx8fHxkb19ub3RfdHJhY2s9IT0xQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAfHx8fHx8fHx8fHx8fHxyZWd1bGFyX3BsdWdpbnM9IT1DaHJvbWUgUERGIFBsdWdpbjo6UG9ydGFibGUgRG9jdW1lbnQgRm9ybWF0OjphcHBsaWNhdGlvbi94LWdvb2dsZS1jaHJvbWUtcGRmfnBkZixDaHJvbWUgUERGIFZpZXdlcjo6OjphcHBsaWNhdGlvbi9wZGZ"
	fg   = "ab2fa79efa2c077b8ab7232cc2646770"
)

func checkError(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func checkResponseStatus(response *http.Response) {
	if response.StatusCode != 200 {
		log.Panicf("status code error: %d %s", response.StatusCode, response.Status)
	}
}

func createLogin() string {
	login := ""
	for index := 0; index < 10; index++ {
		login += string('a' + rand.Int31n('z'-'a'))
	}
	return login
}

func createPostData(csrf, user string) url.Values {
	data := url.Values{
		"step":        []string{"2"},
		"info":        []string{info},
		"fg":          []string{fg},
		"csrf":        []string{csrf},
		"login":       []string{user},
		"email":       []string{user + "@gmail.com"},
		"senha":       []string{user},
		"senha2":      []string{user},
		"operadora[]": []string{"IPTV", "IPTV", "IPTV"},
	}
	return data
}

func getCSRF() string {
	response, err := http.Get(signUpURL)
	checkError(err)
	defer response.Body.Close()
	checkResponseStatus(response)

	document, err := goquery.NewDocumentFromReader(response.Body)
	checkError(err)

	return document.Find("input[name=\"csrf\"]").Nodes[0].Attr[2].Val
}

func signUpNewUser(csrf string) {
	user := createLogin()

	data := createPostData(csrf, user)

	response, err := http.PostForm(signUpURL, data)
	checkError(err)
	defer response.Body.Close()
	checkResponseStatus(response)

	bodyBytes, _ := ioutil.ReadAll(response.Body)
	bodyString := string(bodyBytes)
	fmt.Print(bodyString)
}

func getPlaylistURL() string {
	response, err := http.Get("http://server.azsat.org/")
	checkError(err)
	defer response.Body.Close()
	checkResponseStatus(response)

	response, err = http.Get("http://server.azsat.org/")
	checkError(err)
	defer response.Body.Close()
	checkResponseStatus(response)

	doc, err := goquery.NewDocumentFromReader(response.Body)
	checkError(err)

	return doc.Find("input[name=\"url\"]").Nodes[0].Attr[3].Val
}

func main() {
	http.DefaultClient.Jar, _ = cookiejar.New(nil)

	csrf := getCSRF()

	signUpNewUser(csrf)

	playlistURL := getPlaylistURL()
	log.Print(playlistURL)

	err := ioutil.WriteFile("/data/playlist-url", []byte(playlistURL), 0644)
	checkError(err)
}
