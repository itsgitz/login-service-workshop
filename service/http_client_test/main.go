package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

var selfURL = "/.ory/kratos/public/"

func main() {
	sendPost(
		"http://127.0.0.1:9080/.ory/kratos/public/self-service/browser/flows/login/strategies/password?request=c96ed26e-192b-493c-b555-55086750715a",
		"api",
		"P4ssword123**",
		"Wx15qRb25OOu9IaMSIM23IVPhg45PG/eT+rC2FyECJfwmBtD5q7mHl8fiHSksshaUgOCyJGK69gKpAzn9nwvcg==",
	)
}

func modifyURL() *url.URL {
	requrl := "http://127.0.0.1:9080/.ory/kratos/public/self-service/browser/flows/login/strategies/password?request=09faf96b-b91d-4260-856c-9d92927e8009"
	parsed, err := url.Parse(requrl)
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	selfURLLength := len(selfURL)
	newPath := parsed.Path[selfURLLength-1:]

	parsed.Host = "127.0.0.1:4434"
	parsed.Path = newPath

	log.Println(newPath)
	log.Println(selfURLLength)
	// log.Println(parsed.String())

	return parsed
}

func sendGet() {
	url := "http://127.0.0.1:4434/self-service/browser/flows/login"
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err.Error())
		return
	}

	res, err := client.Do(req)
	if err != nil {
		log.Println(err.Error())
		return
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err.Error())
		return
	}

	log.Println("response body:", string(body))
}

func sendPost(urlAction string, identifier string, password string, csrf string) {
	// HTTP Client
	client := &http.Client{}

	// Define the request body as x-www-form-urlcencoded
	v := url.Values{}
	v.Set("identifier", identifier)
	v.Set("password", password)
	v.Set("csrf_token", csrf)

	// HTTP Request
	req, err := http.NewRequest("POST", urlAction, bytes.NewBufferString(v.Encode()))
	if err != nil {
		log.Println(err.Error())
		return
	}

	// Set header as "application/x-www.form-urlencoded"
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	if err != nil {
		log.Println(err.Error())
		return
	}

	res, err := client.Do(req)
	if err != nil {
		log.Println(err.Error())
		return
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err.Error())
		return
	}

	log.Println("headers sent:", req.Header)
	log.Println("response body:", string(body))
}
