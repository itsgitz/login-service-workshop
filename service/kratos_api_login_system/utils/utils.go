package utils

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

// ParsedJSON function for parsing json data into object
// func ParsedJSON(data []byte) models.ServiceFlowsRequest {
// 	var serviceFlowsRequest models.ServiceFlowsRequest
// 	err := json.Unmarshal(data, &serviceFlowsRequest)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		return models.ServiceFlowsRequest{}
// 	}

// 	return serviceFlowsRequest
// }

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

// GetCookies for get current cookies on request
func GetCookies(c *gin.Context) []*http.Cookie {
	return c.Request.Cookies()
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

// URLModifyForAdminEndpoint function
func URLModifyForAdminEndpoint(requrl string) *url.URL {
	parsed, err := url.Parse(requrl)
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	selfAPILengthString := len(SelfPublicAPI)
	newPath := parsed.Path[selfAPILengthString-1:]

	parsed.Host = "127.0.0.1:4434"
	parsed.Path = newPath

	return parsed
}
