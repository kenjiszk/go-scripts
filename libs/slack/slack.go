package slack

import (
	"log"
	"net/http"
	"net/url"
	"strings"
	"encoding/json"
	"time"

	config "go-scripts/config"
)

type CheckConfig struct {
	Http []HttpConfig `json:"http"`
}

type HttpConfig struct {
	Name   string `json:"name"`
	Host   string `json:"host"`
	Path   string `json:"path"`
	Proto  string `json:"proto"`
	Domain string `json:"domain"`
}

type Slack struct {
	Text     string `json:"text"`
	Username string `json:"username"`
}

func Post(msg string, opt string) {
	params, _ := json.Marshal(Slack{
		msg,
		"MonitoringBot",
	})

	client := &http.Client{Timeout: 5 * time.Second}
	values := url.Values{"payload": {string(params)}}
	postUrl := config.IncomingUrl
	if opt == "test" {
		postUrl = config.IncomingUrlTest
	}
	req, err := http.NewRequest("POST", postUrl, strings.NewReader(values.Encode()))
	if err != nil {
		log.Print(err)
		return
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		log.Print(err)
		return
	}
	defer resp.Body.Close()
}
