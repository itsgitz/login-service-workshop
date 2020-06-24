package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/itsgitz/ory-kratos-workshop/service/kratos_api_login_system/utils"
	kratos_model "github.com/ory/kratos-client-go/models"
)

// LoginRequest is model that used for parsed json body from the login request
type LoginRequest struct {
	UserName string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// ParsedLoginRequestFormAPI for parsed json request body (login request from)
func ParsedLoginRequestFormAPI(c *gin.Context) *LoginRequest {
	lr := LoginRequest{}

	// get body from request (raw request as json)
	body, err := c.GetRawData()
	if err != nil {
		return nil
	}

	// parsing json
	err = json.Unmarshal(body, &lr)
	if err != nil {
		return nil
	}

	return &lr
}

// LoginSubmitRequestWithPostForm function
func LoginSubmitRequestWithPostForm(c *gin.Context, data map[string]string, kratosLoginConfig *kratos_model.LoginRequestMethodConfig) {
	requestNameField := make(map[string]string)
	requestValue := make(map[string]interface{})

	for _, v := range kratosLoginConfig.Fields {
		log.Println(*v.Name, v.Value)
		requestNameField[*v.Name] = *v.Name
		requestValue[*v.Name] = v.Value
	}

	requestBody, err := json.Marshal(map[string]string{
		requestNameField["csrf_token"]: requestValue["csrf_token"].(string),
		requestNameField["identifier"]: data["username"],
		requestNameField["password"]:   data["password"],
	})

	if err != nil {
		utils.ErrorResponse(c, err)
	}

	resp, err := http.Post(*kratosLoginConfig.Action, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	fmt.Fprintf(c.Writer, string(body))
}

// LoginSubmitRequestAPIURLEncoded function for submit login request
func LoginSubmitRequestAPIURLEncoded(c *gin.Context, data map[string]string, kratosLoginConfig *kratos_model.LoginRequestMethodConfig) {
	// adminEndPoint := utils.URLModifyForAdminEndpoint(*kratosLoginConfig.Action)
	getCookies := c.Request.Cookies()
	var cookies []*http.Cookie
	jar, err := cookiejar.New(nil)
	utils.ErrorResponse(c, err)

	requestNameField := make(map[string]string)
	requestValue := make(map[string]interface{})

	for _, v := range kratosLoginConfig.Fields {
		log.Println(*v.Name, v.Value)
		requestNameField[*v.Name] = *v.Name
		requestValue[*v.Name] = v.Value
	}
	u, err := url.Parse(*kratosLoginConfig.Action)
	utils.ErrorResponse(c, err)

	v := url.Values{}
	v.Set("csrf_token", requestValue["csrf_token"].(string))
	v.Set("identifier", data["username"])
	v.Set("password", data["password"])
	postData := strings.NewReader(v.Encode())
	// postData := bytes.NewBufferString(v.Encode())

	for _, c := range getCookies {
		cookies := append(cookies, &http.Cookie{
			Name:       c.Name,
			Value:      c.Value,
			Path:       c.Path,
			Domain:     c.Domain,
			Expires:    c.Expires,
			RawExpires: c.RawExpires,
			MaxAge:     c.MaxAge,
			Secure:     c.Secure,
			HttpOnly:   c.HttpOnly,
			SameSite:   c.SameSite,
			Raw:        c.Raw,
			Unparsed:   c.Unparsed,
		})
		jar.SetCookies(u, cookies)
	}

	// define the http client
	client := &http.Client{
		Jar: jar,
	}
	req, err := http.NewRequest("POST", *kratosLoginConfig.Action, postData)
	// req, err := http.NewRequest("POST", *kratosLoginConfig.Action, bytes.NewBufferString(v.Encode()))
	// req, err := http.NewRequest("POST", adminEndPoint.String(), bytes.NewBufferString(v.Encode()))
	utils.ErrorResponse(c, err)

	// req.RemoteAddr = c.Request.RemoteAddr
	// req.Host = "127.0.0.1:9080"
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	// req.Header.Set("Access-Control-Allow-Origin", "*")
	req.Header.Set("Access-Control-Allow-Credentials", "true")
	// req.Header.Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	// req.Header.Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
	// req.Header.Set("User-Agent", c.Request.Header["User-Agent"][0])

	// just in case csrf_token is missing:
	// req.AddCookie(&http.Cookie{
	// 	Name:  "csrf_token",
	// 	Value: requestValue["csrf_token"].(string),
	// 	Path:  "/",
	// })

	res, err := client.Do(req)
	if err != nil {
		utils.ErrorResponse(c, err)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		utils.ErrorResponse(c, err)
	}

	log.Println("response header:", res.Header)
	log.Println("jar cookies:", jar.Cookies(u))
	log.Println("http headers:", c.Request.Header)
	log.Println("http request data:", data)
	log.Println("request header sent:", req.Header)
	log.Println("remote address:", req.RemoteAddr, c.Request.RemoteAddr)
	log.Println("url action:", *kratosLoginConfig.Action)
	log.Println("posted data:", postData)
	fmt.Fprintf(c.Writer, string(body))

	return

	// get current cookies
	// cookies := c.Request.Cookies()
	// log.Println("url encoded", v.Encode())
	// log.Println("url action", adminEndPoint.String())
	// log.Println("cookies from login submit:", cookies)
	// log.Println("cookies from request login submit", req.Cookies())
	// log.Println("requestValue", requestValue)

	// add cookie for the next request
	// for _, c := range cookies {
	// 	req.AddCookie(&http.Cookie{
	// 		Name:    c.Name,
	// 		Value:   c.Value,
	// 		Path:    c.Path,
	// 		Expires: c.Expires,
	// 	})
	// }
}

// LoginSubmitRequestAPI2 function for submit login request
func LoginSubmitRequestAPI2(c *gin.Context, data map[string]string, kratosLoginConfig *kratos_model.LoginRequestMethodConfig) {
	requestNameField := make(map[string]string)
	requestValue := make(map[string]interface{})

	for _, v := range kratosLoginConfig.Fields {
		log.Println(*v.Name, v.Value)
		requestNameField[*v.Name] = *v.Name
		requestValue[*v.Name] = v.Value
	}

	requestBody, err := json.Marshal(map[string]string{
		requestNameField["csrf_token"]: requestValue["csrf_token"].(string),
		requestNameField["identifier"]: data["username"],
		requestNameField["password"]:   data["password"],
	})

	if err != nil {
		log.Println(err.Error())
	}

	// get current cookies
	cookies := c.Request.Cookies()

	// define the http client
	client := &http.Client{}
	req, err := http.NewRequest("POST", *kratosLoginConfig.Action, bytes.NewBuffer(requestBody))
	if err != nil {
		utils.ErrorResponse(c, err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.RemoteAddr = c.Request.RemoteAddr

	for k, v := range c.Request.Header {
		for _, j := range v {
			req.Header.Set(k, j)
		}
	}

	// add cookie for the next request
	for _, c := range cookies {
		req.AddCookie(&http.Cookie{
			Name:    c.Name,
			Value:   c.Value,
			Path:    c.Path,
			Expires: c.Expires,
		})
	}

	// just in case csrf_token is missing:
	req.AddCookie(&http.Cookie{
		Name:  "csrf_token",
		Value: requestValue["csrf_token"].(string),
		Path:  "/",
	})

	res, err := client.Do(req)
	if err != nil {
		utils.ErrorResponse(c, err)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		utils.ErrorResponse(c, err)
	}

	log.Println(*kratosLoginConfig.Action)
	log.Println(string(body))
	log.Println("cookies from login submit:", cookies)
	log.Println(string(requestBody))
	log.Println("cookies from request login submit", req.Cookies())
	fmt.Fprintf(c.Writer, string(body))

	return
}

// LoginSubmitRequestAPI function for submit login request
func LoginSubmitRequestAPI(c *gin.Context, login *LoginRequest, kratosLoginConfig *kratos_model.LoginRequestMethodConfig) {
	requestNameField := make(map[string]string)
	requestValue := make(map[string]interface{})

	for _, v := range kratosLoginConfig.Fields {
		log.Println(*v.Name, v.Value)
		requestNameField[*v.Name] = *v.Name
		requestValue[*v.Name] = v.Value
	}

	log.Println(requestNameField)
	log.Println(requestValue["csrf_token"])

	// define request body as json
	requestBody, err := json.Marshal(map[string]string{
		requestNameField["csrf_token"]: requestValue["csrf_token"].(string),
		requestNameField["identifier"]: login.UserName,
		requestNameField["password"]:   login.Password,
	})

	if err != nil {
		utils.ErrorResponse(c, err)
	}

	// get current cookies
	cookies := c.Request.Cookies()

	// define the http client
	client := &http.Client{}
	req, err := http.NewRequest("POST", *kratosLoginConfig.Action, bytes.NewBuffer(requestBody))
	if err != nil {
		utils.ErrorResponse(c, err)
	}

	req.Header.Set("Content-Type", "application/json")

	// add cookie for the next request
	for _, c := range cookies {
		req.AddCookie(&http.Cookie{
			Name:    c.Name,
			Value:   c.Value,
			Path:    c.Path,
			Expires: c.Expires,
		})
	}

	res, err := client.Do(req)
	if err != nil {
		utils.ErrorResponse(c, err)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		utils.ErrorResponse(c, err)
	}

	log.Println(*kratosLoginConfig.Action)
	log.Println(string(body))
	log.Println("cookies from login submit:", cookies)

	return
}
