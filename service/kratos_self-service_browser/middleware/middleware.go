package middleware

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/itsgitz/ory-kratos-workshop/service/kratos_self-service_browser/utils"
)

func CustomMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestURL := c.Request.RequestURI

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
