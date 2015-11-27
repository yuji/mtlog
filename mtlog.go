package main

import (
	"os"
	"github.com/codegangsta/cli"
	"encoding/json"
	"fmt"
	"log"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Token struct {
	AccessToken string `json:"accessToken"`
	SessionID   string `json:"sessionId"`
}

type Log struct {
	IP      string `json:"ip"`
	Date    string `json:"date"`
	Message string `json:"message"`
	Level   string `json:"level"`
}

type Logs struct {
	Total string `json:"totalResults"`
	Items []*Log
}

func authentication(ch chan string) {
	authURL := "http://localhost/cgi-bin/mt/mt-data-api.cgi/v3/authentication"
	authParam := url.Values{}
	authParam.Add("username", "takayama")
	authParam.Add("password", "password")
	authParam.Add("clientId", "mtlog")

	client := &http.Client{}

	req, err := http.NewRequest("POST", authURL, strings.NewReader(authParam.Encode()))
	if err != nil {
		log.Fatal(err)
	}

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Fatal(res)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var token Token
	err = json.Unmarshal(body, &token)
	if err != nil {
		log.Fatal(err)
	}

	ch <- token.AccessToken
}

func main() {
	app := cli.NewApp()
	app.Name = "mtlog"
	app.Usage = "print arguments"
	app.Version = "1.0.0"
	app.Action = func(c *cli.Context) {
		ch := make(chan string)
		go authentication(ch)
		
		token := <-ch

		logURL := "http://localhost/cgi-bin/mt/mt-data-api.cgi/v3/sites/1/logs?dateFrom=2015-11-24&dateTo=2015-11-25"

		client := &http.Client{}

		req, err := http.NewRequest("GET", logURL, nil)
		if err != nil {
			log.Fatal(err)
		}

		req.Header.Add("X-MT-Authorization", "MTAuth accessToken=" + token)
		res, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			log.Fatal(res)
		}

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}

		var logs Logs
		err = json.Unmarshal(body, &logs)
		if err != nil {
			log.Fatal(err)
		}

		for i := 0; i < len(logs.Items); i++ {
			log.Println(fmt.Sprintf("%s\t%s\t%s\t%s", logs.Items[i].Date, logs.Items[i].Level, logs.Items[i].IP, logs.Items[i].Message))
		}
	}
	app.Run(os.Args)
}

