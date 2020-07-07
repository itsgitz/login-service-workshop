package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/itsgitz/login-service-workshop/service/kratos_self-service_browser/middleware"
	"github.com/itsgitz/login-service-workshop/service/kratos_self-service_browser/utils"
	hydraModels "github.com/ory/hydra-client-go/models"
	kratos "github.com/ory/kratos-client-go/client"
	"github.com/ory/kratos-client-go/client/common"
	"golang.org/x/oauth2"
)

var (
	adminKratosClient  *kratos.OryKratos
	publicKratosClient *kratos.OryKratos
	hydraOauth2Config  *oauth2.Config
)

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
	r.Use(middleware.CustomMiddleware())

	// Main page / home page, say hello!
	r.GET("/", func(c *gin.Context) {
		log.Println("ADMIN_API:", utils.AdminAPI)
		tpl := template.Must(template.ParseFiles("views/main.html"))
		tpl.Execute(c.Writer, nil)
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

		tpl := template.Must(template.ParseFiles("views/login.html"))
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
		tpl := template.Must(template.ParseFiles("views/registration.html"))
		tpl.Execute(c.Writer, config)
	})

	// credentials page
	r.GET("/credentials", func(c *gin.Context) {
		session, err := utils.GetCurrentSession(c)
		if err != nil {
			utils.ErrorResponse(c, err)
		}

		// get query parameter
		createNewClient := c.Query("create")

		if createNewClient == "true" {
			hydraOauth2Client := hydraModels.OAuth2Client{
				RedirectUris:            []string{"http://127.0.0.1:8080/callback"},
				GrantTypes:              []string{"authorization_code", "refresh_token"},
				ResponseTypes:           []string{"code", "id_token"},
				Scope:                   "openid offline",
				ClientName:              session.Identity.Traits.Username,
				Contacts:                []string{session.Identity.Traits.Email},
				TokenEndpointAuthMethod: "client_secret_post",
			}

			jsonData, err := json.Marshal(hydraOauth2Client)
			if err != nil {
				log.Println("json encoding error:", err.Error())
				return
			}

			log.Println("json:", string(jsonData))

			httpClient := http.Client{}
			req, err := http.NewRequest("POST", "http://hydra:4445/clients", bytes.NewBuffer(jsonData))
			if err != nil {
				log.Println("http request error:", err.Error())
				return
			}

			req.Header.Set("Content-Type", "application/json")

			res, err := httpClient.Do(req)
			if err != nil {
				log.Println("http response error:", err.Error())
				return
			}

			defer res.Body.Close()

			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				log.Println("http body read error:", err.Error())
				return
			}

			log.Println(string(body))

			hydraResponse := &hydraModels.OAuth2Client{}
			err = json.Unmarshal(body, hydraResponse)
			if err != nil {
				log.Println("json unmarshal error ->", err.Error())
			}

			log.Println("client_id", hydraResponse.ClientID)
			log.Println("client_secret", hydraResponse.ClientSecret)
		}

		tpl := template.Must(template.ParseFiles("views/credentials.html"))
		err = tpl.Execute(c.Writer, nil)
		if err != nil {
			log.Println("template error:", err.Error())
			return
		}
	})

	// get current session
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

	r.GET("/logout", func(c *gin.Context) {
		director := func(req *http.Request) {
			req.URL.Scheme = "http"
			req.URL.Host = "kratos:4433"
			req.URL.Path = "/self-service/browser/flows/logout"
			req.Header.Add("X-Forwarded-Host", req.Host)
			req.Header.Add("X-Origin-Host", "kratos:4433")
		}

		proxy := &httputil.ReverseProxy{Director: director}
		proxy.ServeHTTP(c.Writer, c.Request)
		return
	})

	// Error page
	r.GET("/error", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"error": "This is error message",
		})
	})

	// Settings page
	r.GET("/settings", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "This is setting page",
		})
	})

	r.Run(":9080")
}
