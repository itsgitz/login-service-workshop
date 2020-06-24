package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/itsgitz/ory-kratos-workshop/service/kratos_api_login_system/middleware"
	"github.com/itsgitz/ory-kratos-workshop/service/kratos_api_login_system/models"
	"github.com/itsgitz/ory-kratos-workshop/service/kratos_api_login_system/utils"
	kratos "github.com/ory/kratos-client-go/client"
	"github.com/ory/kratos-client-go/client/common"
	"github.com/ory/kratos-client-go/client/public"
)

// Data type string interface
type Data string

var adminKratosClient *kratos.OryKratos
var publicKratosClient *kratos.OryKratos

var requestBody Data = ""

var store = sessions.NewCookieStore([]byte("secret"))

func init() {
	adminHost, err := url.Parse(utils.AdminAPI)
	if err != nil {
		log.Println(err.Error())
	}

	clientHost, err := url.Parse(utils.PublicAPI)
	if err != nil {
		log.Println(err.Error())
	}

	adminKratosClient = kratos.NewHTTPClientWithConfig(
		nil,
		&kratos.TransportConfig{
			Schemes:  []string{adminHost.Scheme},
			Host:     adminHost.Host,
			BasePath: adminHost.Path,
		},
	)

	publicKratosClient = kratos.NewHTTPClientWithConfig(
		nil,
		&kratos.TransportConfig{
			Schemes:  []string{clientHost.Scheme},
			Host:     clientHost.Host,
			BasePath: clientHost.Path,
		},
	)
}

type MyTransport struct {
	RT http.RoundTripper
}

func (t *MyTransport) Transport() http.RoundTripper {
	if t.RT != nil {
		return t.RT
	}

	return http.DefaultTransport
}

func (t *MyTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	return t.RT.RoundTrip(r)
}

func main() {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.ProxyRequestMiddleware())

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to API Login system",
		})
	})

	// registration endpoint, will redirect to /auth/registration
	r.GET("/registration", func(c *gin.Context) {
		redirectURL := fmt.Sprintf("%s%sself-service/browser/flows/registration", utils.ClientAPI, utils.SelfPublicAPI)
		c.Redirect(http.StatusSeeOther, redirectURL)
	})

	r.GET("/login", func(c *gin.Context) {
		s, err := store.Get(c.Request, "request-body")
		utils.ErrorResponse(c, err)

		jsonParsed := models.ParsedLoginRequestFormAPI(c)

		s.Values["username"] = jsonParsed.UserName
		s.Values["password"] = jsonParsed.Password

		err = s.Save(c.Request, c.Writer)
		utils.ErrorResponse(c, err)

		log.Println("parsed login request", models.ParsedLoginRequestFormAPI(c))
		log.Println("remote address:", c.Request.RemoteAddr)

		redirectURL := fmt.Sprintf("%s%sself-service/browser/flows/login", utils.ClientAPI, utils.SelfPublicAPI)
		c.Redirect(http.StatusSeeOther, redirectURL)
	})

	r.GET("/auth/login", func(c *gin.Context) {
		data := make(map[string]string)

		s, err := store.Get(c.Request, "request-body")
		utils.ErrorResponse(c, err)

		getRequest := c.Request.URL.Query().Get("request")

		// get kratos params using request id
		params := common.NewGetSelfServiceBrowserLoginRequestParams()
		params.WithRequest(getRequest)

		resp, err := adminKratosClient.Common.GetSelfServiceBrowserLoginRequest(params)
		utils.ErrorResponse(c, err)

		config := resp.GetPayload().Methods["password"].Config

		log.Println("config:", *config.Action)
		log.Println("csrf_token:", config.Fields[2].Value.(string))
		log.Println("sessions values:", s.Values)
		log.Println("get request id:", getRequest)

		data["username"] = s.Values["username"].(string)
		data["password"] = s.Values["password"].(string)

		// uv := url.Values{}
		// uv.Set("csrf_token", config.Fields[2].Value.(string))
		// uv.Set("identifier", s.Values["username"].(string))
		// uv.Set("password", s.Values["password"].(string))

		// clientRequest := &http.Request{
		// 	Form: uv,
		// }

		// clientTransport := &MyTransport{}
		// clientTransport.RT.RoundTrip(clientRequest)
		// utils.ErrorResponse(c, err)

		// client := &http.Client{}
		// client.Transport = clientTransport

		// completeParams := public.NewCompleteSelfServiceBrowserSettingsPasswordStrategyFlowParams()
		// completeParams.SetHTTPClient(client)

		// err = publicKratosClient.Public.CompleteSelfServiceBrowserSettingsPasswordStrategyFlow(completeParams)
		// utils.ErrorResponse(c, err)

		// log.Println("params:", params)
		// log.Println("complete params:", completeParams)

		models.LoginSubmitRequestAPIURLEncoded(c, data, config)
		return
	})

	// r.GET("/login", func(c *gin.Context) {
	// 	redirectURL := fmt.Sprintf("%s%sself-service/browser/flows/login", utils.ClientAPI, utils.SelfPublicAPI)
	// 	// client := &http.Client{}
	// 	// req, err := http.NewRequest("GET", redirectURL, nil)
	// 	// utils.ErrorResponse(c, err)

	// 	// res, err := client.Do(req)
	// 	// utils.ErrorResponse(c, err)

	// 	// defer res.Body.Close()

	// 	// return

	// 	// // body, err := ioutil.ReadAll(res.Body)
	// 	// // utils.ErrorResponse(c, err)

	// 	// // // fmt.Fprintf(c.Writer, string(body))
	// 	c.Redirect(http.StatusSeeOther, redirectURL)
	// })

	// r.GET("/auth/login", func(c *gin.Context) {
	// 	log.Println("cookies login", c.Request.Cookies())
	// 	getRequest := c.Request.URL.Query().Get("request")
	// 	requrl := fmt.Sprintf("%s%sself-service/browser/flows/requests/login?request=%s", utils.ClientAPI, utils.SelfPublicAPI, getRequest)

	// 	client := &http.Client{}
	// 	req, err := http.NewRequest("GET", requrl, nil)
	// 	utils.ErrorResponse(c, err)

	// 	res, err := client.Do(req)
	// 	utils.ErrorResponse(c, err)

	// 	defer res.Body.Close()

	// 	body, err := ioutil.ReadAll(res.Body)
	// 	utils.ErrorResponse(c, err)

	// 	c.Writer.Header().Set("Content-Type", "application/json")
	// 	fmt.Fprintf(c.Writer, string(body))
	// })

	// r.GET("/auth/login", func(c *gin.Context) {
	// 	getRequest := c.Request.URL.Query().Get("request")
	// 	data := make(map[string]string)

	// 	s, err := store.Get(c.Request, "request-body")
	// 	utils.ErrorResponse(c, err)

	// 	data["username"] = s.Values["username"].(string)
	// 	data["password"] = s.Values["password"].(string)

	// 	// get request context
	// 	client := &http.Client{}
	// 	requrl := fmt.Sprintf("http://127.0.0.1:9080/.ory/kratos/public/self-service/browser/flows/requests/login?request=%s", getRequest)
	// 	req, err := http.NewRequest("GET", requrl, nil)
	// 	utils.ErrorResponse(c, err)

	// 	res, err := client.Do(req)
	// 	utils.ErrorResponse(c, err)

	// 	defer res.Body.Close()

	// 	body, err := ioutil.ReadAll(res.Body)
	// 	utils.ErrorResponse(c, err)

	// 	// models.LoginSubmitRequestAPIURLEncoded(c, data, config)
	// 	// models.LoginSubmitRequestAPI2(c, data, config)

	// 	c.Writer.Header().Set("Content-Type", "application/json")
	// 	fmt.Fprintf(c.Writer, string(body))
	// })

	// registration proccess for ory kratos
	r.GET("/auth/registration", func(c *gin.Context) {
		getRequest := c.Request.URL.Query().Get("request")

		params := common.NewGetSelfServiceBrowserRegistrationRequestParams()
		params.WithRequest(getRequest)

		resp, err := adminKratosClient.Common.GetSelfServiceBrowserRegistrationRequest(params)
		utils.ErrorResponse(c, err)

		config := resp.GetPayload().Methods["password"].Config

		// define template
		tpl := template.Must(template.ParseFiles("service/kratos_client_api/views/registration.html"))
		tpl.Execute(c.Writer, config)
	})

	r.GET("/logout", func(c *gin.Context) {
		redirectURL := fmt.Sprintf("%s%sself-service/browser/flows/logout", utils.ClientAPI, utils.SelfPublicAPI)
		c.Redirect(http.StatusFound, redirectURL)
	})

	r.GET("/session", func(c *gin.Context) {
		cookies := c.Request.Cookies()

		client := &http.Client{}
		req, err := http.NewRequest("GET", "http://127.0.0.1:9080/.ory/kratos/public/sessions/whoami", nil)
		utils.ErrorResponse(c, err)

		// Add session cookie
		for _, cookie := range cookies {
			req.AddCookie(&http.Cookie{
				Name:    cookie.Name,
				Value:   cookie.Value,
				Path:    cookie.Path,
				Expires: cookie.Expires,
			})
		}

		res, err := client.Do(req)
		utils.ErrorResponse(c, err)

		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		utils.ErrorResponse(c, err)

		c.Writer.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(c.Writer, string(body))
		return
	})

	r.GET("/session/sdk", func(c *gin.Context) {
		params := public.NewWhoamiParams()
		resp, err := publicKratosClient.Public.Whoami(params)
		utils.ErrorResponse(c, err)

		// config := resp.GetPayload().Methods["password"].Config
		s := resp.GetPayload()
		log.Println("Session:", s)
	})

	r.Run(":9080")
}
