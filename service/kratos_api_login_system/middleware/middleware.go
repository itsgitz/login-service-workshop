package middleware

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/itsgitz/ory-kratos-workshop/service/kratos_api_login_system/utils"
)

// ProxyRequestMiddleware is function act as middleware, forwarding all request to ORY Kratos Host
func ProxyRequestMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		requestURL := c.Request.RequestURI
		// cookies := c.Request.Cookies()

		// if request uri contains path '/.ory/kratos/public`, then redirect itu to PublicAPI host
		if strings.Contains(requestURL, utils.SelfPublicAPI) {
			getLengthOfSelfPublicAPIURL := len(utils.SelfPublicAPI)
			newRequestURL := fmt.Sprintf("%s/%s", utils.PublicAPI, requestURL[getLengthOfSelfPublicAPIURL:])
			// newRequestURL := fmt.Sprintf("%s/%s", utils.AdminAPI, requestURL[getLengthOfSelfPublicAPIURL:])

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
