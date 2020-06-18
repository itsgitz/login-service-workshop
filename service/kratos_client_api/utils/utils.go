package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/itsgitz/ory-kratos-workshop/service/kratos_client_api/models"
)

var (
	AdminAPI      = "http://127.0.0.1:4434"
	PublicAPI     = "http://127.0.0.1:4433"
	ClientAPI     = "http://127.0.0.1:9080"
	SelfPublicAPI = "/.ory/kratos/public/"
)

// GenerateNewUUID for generate new UUID4
func GenerateNewUUID() interface{} {
	// create new uiid
	userID, err := uuid.NewUUID()
	if err != nil {
		log.Fatal(err.Error())
	}

	return userID
}

// HTTPClient function
func HTTPClient(c *gin.Context, url string, method string, payload io.Reader) []byte {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return nil
	}

	res, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return nil
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return nil
	}

	return body
}

// HTTPClientForGetSession function
func HTTPClientForGetSession(c *gin.Context, url string, method string, dataCookie []*http.Cookie) []byte {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return nil
	}

	// set cookie in request
	for _, cookie := range dataCookie {
		data := make(map[string]string)
		data["name"] = cookie.Name
		data["value"] = cookie.Value
		data["path"] = "/"

		AddCookie(c, data)
	}

	res, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return nil
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return nil
	}

	return body
}

// PostRequest for making post request to server
// func PostRequest(c *gin.Context, models models.ServiceFlowsRequest, requestValue string) {
// 	requestURL := fmt.Sprintf("%s/self-service/browser/flows/registration/strategies/password?request=%s", PublicAPI, requestValue)

// 	var err error
// 	var client = &http.Client{}
// 	var param = url.Values{}
// 	csrfToken := models.Methods.Password.Config.Fields[0].Value

// 	param.Set("csrf_token", csrfToken)
// 	payload := bytes.NewBufferString(param.Encode())

// 	request, err := http.NewRequest("POST", requestURL, payload)
// 	if err != nil {
// 		log.Println(err.Error())
// 	}

// 	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
// 	request.Header.Set("X-CSRF-TOKEN", csrfToken)
// 	response, err := client.Do(request)
// 	if err != nil {
// 		log.Println(err.Error())
// 	}

// 	defer response.Body.Close()

// 	body, err := ioutil.ReadAll(response.Body)
// 	if err != nil {
// 		log.Println(err.Error())
// 	}

// 	fmt.Fprintf(c.Writer, string(body))
// 	// fmt.Println(models.Methods.Password.Config)
// 	// for k, v := range models.Methods.Password.Config.Fields {
// 	// 	fmt.Println(k, v)
// 	// }

// 	// fmt.Println("First form is:", models.Methods.Password.Config.Fields[0])

// 	// requestBody, err := json.Marshal(map[string]string{
// 	// 	models.Methods.Password.Config.Fields[0].Name: models.Methods.Password.Config.Fields[0].Value,
// 	// })

// 	// if err != nil {
// 	// 	log.Println(err.Error())
// 	// }

// 	// resp, err := http.Post(requestURL, "application/json", bytes.NewBuffer(requestBody))
// 	// if err != nil {
// 	// 	log.Println(err.Error())
// 	// }

// 	// defer resp.Body.Close()

// 	// body, err := ioutil.ReadAll(resp.Body)
// 	// if err != nil {
// 	// 	log.Println(err.Error())
// 	// }

// 	// fmt.Fprintf(c.Writer, string(body))
// }

// ParsedJSON function for parsing json data into object
func ParsedJSON(data []byte) models.ServiceFlowsRequest {
	var serviceFlowsRequest models.ServiceFlowsRequest
	err := json.Unmarshal(data, &serviceFlowsRequest)
	if err != nil {
		fmt.Println(err.Error())
		return models.ServiceFlowsRequest{}
	}

	return serviceFlowsRequest
}

// AddCookie function for add cookie in request/response header
func AddCookie(c *gin.Context, data map[string]string) {
	expire := time.Now().Add(30 * time.Minute)
	cookie := http.Cookie{
		Name:     data["name"],
		Value:    data["value"],
		Path:     data["path"],
		Expires:  expire,
		Secure:   false,
		HttpOnly: false,
	}

	http.SetCookie(c.Writer, &cookie)
}

// RemoveCookie function for remove cookie from the system
func RemoveCookie(c *gin.Context, data map[string]string) {
	cookie := http.Cookie{
		Name:   data["name"],
		Value:  data["value"],
		Path:   data["path"],
		MaxAge: -1,
	}

	http.SetCookie(c.Writer, &cookie)
}

// ErrorResponse function for show error message in JSON format
func ErrorResponse(c *gin.Context, err error) {
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})

		return
	}
}
