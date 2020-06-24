package example

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func main() {

	url := "http://localhost:4433/identities"
	method := "POST"

	payload := strings.NewReader("{\n  \"addresses\": [\n    {\n      \"expires_at\": null,\n      \"id\": \"929asdn1022-kivmue\",\n      \"value\": \"anggit@isi.co.id\",\n      \"verified\": false,\n      \"verified_at\": null,\n      \"via\": \"email\"\n    }\n  ],\n  \"id\": \"f92kak012920\",\n  \"traits\": {\n  	\"email\": \"anggit@isi.co.id\"\n  },\n  \"traits_schema_id\": \"default\",\n  \"traits_schema_url\": \"http://127.0.0.1:4455/.ory/kratos/public/schemas/default\"\n}")

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	fmt.Println(string(body))
}
