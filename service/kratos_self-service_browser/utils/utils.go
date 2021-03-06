package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/itsgitz/login-service-workshop/service/kratos_self-service_browser/models"
)

var (
	AdminAPI      string
	PublicAPI     string
	ClientAPI     string
	SelfPublicAPI string
)

// var (
// 	AdminAPI      = "http://127.0.0.1:4434"
// 	PublicAPI     = "http://127.0.0.1:4433"
// 	ClientAPI     = "http://127.0.0.1:9080"
// 	SelfPublicAPI = "/.ory/kratos/public/"
// )

func init() {
	AdminAPI = os.Getenv("ADMIN_API")
	PublicAPI = os.Getenv("PUBLIC_API")
	ClientAPI = os.Getenv("CLIENT_API")
	SelfPublicAPI = os.Getenv("SELF_PUBLIC_API_PATH")
}

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

// GetCurrentSession for retrieve current logged in session
func GetCurrentSession(c *gin.Context) (*models.CurrentKratosSession, error) {
	cookies := c.Request.Cookies()

	cl := http.Client{}
	req, err := http.NewRequest("GET", "http://127.0.0.1:9080/.ory/kratos/public/sessions/whoami", nil)
	if err != nil {
		return nil, err
	}

	// Add session cookie
	for _, cookie := range cookies {
		req.AddCookie(&http.Cookie{
			Name:    cookie.Name,
			Value:   cookie.Value,
			Path:    cookie.Path,
			Expires: cookie.Expires,
		})
	}

	resp, err := cl.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	session := &models.CurrentKratosSession{}

	err = json.Unmarshal(body, session)
	if err != nil {
		return nil, err
	}

	return session, nil
}
