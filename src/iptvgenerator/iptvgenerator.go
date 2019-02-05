package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

func createLogin() string {
	login := ""
	for index := 0; index < 10; index++ {
		login += string('a' + rand.Int31n('z'-'a'))
	}
	return login
}

func elementIsPresent(by, value string) selenium.Condition {
	return func(wd selenium.WebDriver) (bool, error) {
		_, err := wd.FindElement(by, value)
		return err == nil, nil
	}
}

func elementIsDisplayed(by, value string) selenium.Condition {
	return func(wd selenium.WebDriver) (bool, error) {
		elem, err := wd.FindElement(by, value)
		if err != nil {
			return false, nil
		}

		displayed, err := elem.IsDisplayed()
		return displayed, nil
	}
}

func checkError(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func sendKeysToElem(wd selenium.WebDriver, selectorBy, selectorValue, sendValue string) {
	elem, err := wd.FindElement(selectorBy, selectorValue)
	checkError(err)
	err = elem.SendKeys(sendValue)
	checkError(err)
}

func clickElem(wd selenium.WebDriver, selectorBy, selectorValue string) {
	elem, err := wd.FindElement(selectorBy, selectorValue)
	checkError(err)
	err = elem.Click()
	checkError(err)
}
func main() {
	const (
		seleniumPath     = "/usr/bin/selenium-server-standalone.jar"
		chromeDriverPath = "/usr/bin/chromedriver"
		port             = 9393

		loginTimeout   = 5
		elementTimeout = 10 * time.Second
	)
	opts := []selenium.ServiceOption{
		selenium.StartFrameBuffer(), // Start an X frame buffer for the browser to run in.
		selenium.ChromeDriver(chromeDriverPath),
		selenium.Output(os.Stderr),
	}

	service, err := selenium.NewSeleniumService(seleniumPath, port, opts...)
	checkError(err)
	defer service.Stop()

	// Connect to the WebDriver instance running locally.
	caps := selenium.Capabilities{"browserName": "chrome"}
	caps.AddChrome(chrome.Capabilities{
		Args: []string{
			"--no-sandbox",
		},
	})
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	checkError(err)
	defer wd.Quit()

	for login := 0; login < loginTimeout; login++ {
		err = wd.Get("https://server.azsat.org/novo/usuario/registrar_iptv.php")
		checkError(err)

		user := createLogin()
		sendKeysToElem(wd, selenium.ByID, "login", user)
		sendKeysToElem(wd, selenium.ByID, "email", user+"@gmail.com")
		sendKeysToElem(wd, selenium.ByID, "senha", user)
		sendKeysToElem(wd, selenium.ByID, "senha2", user)

		clickElem(wd, selenium.ByCSSSelector, "button[type=\"submit\"]")

		err = wd.WaitWithTimeout(elementIsPresent(selenium.ByCSSSelector, "#btnMore a"), elementTimeout)
		if err != nil {
			present, err := elementIsPresent(selenium.ByID, "modal-window")(wd)
			checkError(err)
			if !present || (present && login == loginTimeout-1) {
				log.Panic(errors.New("Falha ao logar"))
			}
		} else {
			clickElem(wd, selenium.ByCSSSelector, "#btnMore a")
			break
		}
	}

	err = wd.WaitWithTimeout(elementIsDisplayed(selenium.ByID, "operadora2"), elementTimeout)
	checkError(err)
	sendKeysToElem(wd, selenium.ByID, "operadora", "I")
	sendKeysToElem(wd, selenium.ByID, "operadora2", "I")
	sendKeysToElem(wd, selenium.ByID, "operadora3", "I")

	clickElem(wd, selenium.ByCSSSelector, "button[type=\"submit\"]")

	err = wd.WaitWithTimeout(elementIsPresent(selenium.ByCSSSelector, "a.modal-btn.btn-green"), elementTimeout)
	checkError(err)
	clickElem(wd, selenium.ByCSSSelector, "a.modal-btn.btn-green")

	wd.Refresh()

	err = wd.WaitWithTimeout(elementIsPresent(selenium.ByName, "url"), elementTimeout)
	checkError(err)
	elem, err := wd.FindElement(selenium.ByName, "url")
	checkError(err)
	url, err := elem.GetAttribute("value")
	checkError(err)

	log.Printf("New url '%s'", url)

	err = ioutil.WriteFile("playlist-url", []byte(url), 0644)
	checkError(err)
}
