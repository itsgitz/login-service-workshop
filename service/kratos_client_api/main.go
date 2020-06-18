package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/itsgitz/ory-kratos-workshop/service/kratos_client_api/utils"
	kratos "github.com/ory/kratos-client-go/client"
	"github.com/ory/kratos-client-go/client/common"
)

// kratos "github.com/ory/kratos-client-go/client"

var adminKratosClient *kratos.OryKratos
var publicKratosClient *kratos.OryKratos

func customMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestURL := c.Request.RequestURI
		cookies := c.Request.Cookies()

		log.Println("cookies from middleware:", cookies)

		// if request uri contains path '/.ory/kratos/public`, then redirect itu to PublicAPI host
		if strings.Contains(requestURL, utils.SelfPublicAPI) {
			getLengthOfSelfPublicAPIURL := len(utils.SelfPublicAPI)
			newRequestURL := fmt.Sprintf("%s/%s", utils.PublicAPI, requestURL[getLengthOfSelfPublicAPIURL:])

			proxyURL, err := url.Parse(newRequestURL)
			utils.ErrorResponse(c, err)

			director := func(req *http.Request) {
				req.URL.Scheme = proxyURL.Scheme
				req.URL.Host = proxyURL.Host
				req.URL.Path = proxyURL.Path
				req.Header.Add("X-Forwarded-Host", req.Host)
				req.Header.Add("X-Origin-Host", proxyURL.Host)
			}

			proxy := &httputil.ReverseProxy{Director: director}
			proxy.ServeHTTP(c.Writer, c.Request)
		}

		c.Next()
	}
}

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

func main() {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(customMiddleware())

	// Main page / home page, say hello!
	r.GET("/", func(c *gin.Context) {
		cookies := c.Request.Cookies()

		log.Println("Show cookies", cookies)

		c.JSON(http.StatusOK, gin.H{
			"message": "hello!",
		})
	})

	r.GET("/error", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"error": "This is error message",
		})
	})

	r.GET("/settings", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "This is setting page",
		})
	})

	// registration endpoint, will redirect to /auth/registration
	r.GET("/registration", func(c *gin.Context) {
		redirectURL := fmt.Sprintf("%s%sself-service/browser/flows/registration", utils.ClientAPI, utils.SelfPublicAPI)
		c.Redirect(http.StatusSeeOther, redirectURL)
	})

	// login endpoint
	r.GET("/login", func(c *gin.Context) {
		redirectURL := fmt.Sprintf("%s%sself-service/browser/flows/login", utils.ClientAPI, utils.SelfPublicAPI)
		c.Redirect(http.StatusSeeOther, redirectURL)
	})

	// login process for ory kratos
	r.GET("/auth/login", func(c *gin.Context) {
		getRequest := c.Request.URL.Query().Get("request")

		params := common.NewGetSelfServiceBrowserLoginRequestParams()
		params.WithRequest(getRequest)

		resp, err := adminKratosClient.Common.GetSelfServiceBrowserLoginRequest(params)
		utils.ErrorResponse(c, err)

		config := resp.GetPayload().Methods["password"].Config

		tpl := template.Must(template.ParseFiles("service/kratos_client_api/views/login.html"))
		tpl.Execute(c.Writer, config)
	})

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

	r.GET("/session/manual", func(c *gin.Context) {
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

	r.GET("/session", func(c *gin.Context) {
		director := func(req *http.Request) {
			req.URL.Scheme = "http"
			req.URL.Host = "127.0.0.1:4433"
			req.URL.Path = "/sessions/whoami"
			req.Header.Add("X-Forwarded-Host", req.Host)
			req.Header.Add("X-Origin-Host", "127.0.0.1:4433")
		}

		proxy := &httputil.ReverseProxy{Director: director}
		proxy.ServeHTTP(c.Writer, c.Request)
		return
	})

	r.GET("/logout", func(c *gin.Context) {
		director := func(req *http.Request) {
			req.URL.Scheme = "http"
			req.URL.Host = "127.0.0.1:4433"
			req.URL.Path = "/self-service/browser/flows/logout"
			req.Header.Add("X-Forwarded-Host", req.Host)
			req.Header.Add("X-Origin-Host", "127.0.0.1:4433")
		}

		proxy := &httputil.ReverseProxy{Director: director}
		proxy.ServeHTTP(c.Writer, c.Request)
		return
	})

	r.Run(":9080")
}
