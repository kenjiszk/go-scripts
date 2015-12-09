package main

import (
	"net/http"
	"crypto/tls"
	"time"
	"fmt"
	"flag"

	slack "go-scripts/libs/slack"
	config "go-scripts/config"
)

func main() {
	flag.Parse()
	if flag.Arg(0) == "" || flag.Arg(0) == "test" {
		for _, target := range config.HttpTargets() {
			check(target, flag.Arg(0))
		}
	} else {
		fmt.Println("Invalid Argment", flag.Arg(0))
	}
}

func check(target config.HttpConfig, opt string) {
	errorNum := 0
	checkNum := 0
	fatalNum := 2
	errorMessage := ""
	for checkNum < fatalNum {
		checkNum += 1
		targetPath := target.Proto + "://" + target.Host + target.Path
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Timeout: 5 * time.Second, Transport: tr}
		req, err := http.NewRequest("GET", targetPath, nil)
		if err != nil {
			slack.Post(err.Error(), opt)
			return
		}

		req.Header.Add("Host", target.Domain)
		resp, err := client.Do(req)
		if err != nil {
			slack.Post(err.Error(), opt)
			return
		}

		if resp.StatusCode != 200 {
			errorNum += 1
			if (errorNum >= fatalNum) {
				errorMessage += targetPath + " [" + target.Name + "] " + "returns " + fmt.Sprint(resp.StatusCode) + "\n"
			}
		} else {
			fmt.Println(target.Domain, resp.StatusCode)
			break
		}

		defer resp.Body.Close()
	}

	if errorMessage != "" {
		fmt.Print(errorMessage)
		slack.Post(errorMessage, opt)
	}
}
