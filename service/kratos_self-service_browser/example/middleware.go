func customMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		requestURL := c.Request.RequestURI
		cookies := c.Request.Cookies()

		log.Println("request url:", requestURL)
		log.Println("cookies from middleware", cookies)

		// if request uri contains path '/.ory/kratos/public`, then redirect itu to PublicAPI host
		if strings.Contains(requestURL, utils.SelfPublicAPI) {
			getLengthOfSelfPublicAPIURL := len(utils.SelfPublicAPI)
			log.Println("cut", requestURL[getLengthOfSelfPublicAPIURL:])

			newRequestURL := fmt.Sprintf("%s/%s", utils.PublicAPI, requestURL[getLengthOfSelfPublicAPIURL:])
			log.Println("New Request URL", newRequestURL)

			// redirect using http
			// c.Redirect(http.StatusFound, newRequestURL)

			proxyURL, err := url.Parse(newRequestURL)
			if err != nil {
				log.Println(err.Error())
			}

			proxy := httputil.NewSingleHostReverseProxy(proxyURL)
			proxy.ServeHTTP(c.Writer, c.Request)

			// director := func(req *http.Request) {
			// 	req.URL.Scheme = "http"
			// 	req.URL.Host = "127.0.0.1:4433"
			// 	req.URL.Path = fmt.Sprintf("/%s", requestURL[getLengthOfSelfPublicAPIURL:])

			// 	for _, cookie := range cookies {
			// 		req.AddCookie(&http.Cookie{
			// 			Name:  cookie.Name,
			// 			Value: cookie.Value,
			// 			Path:  "/",
			// 		})
			// 	}
			// }

			// proxy := &httputil.ReverseProxy{Director: director}
			// proxy.ServeHTTP(c.Writer, c.Request)
		}

		c.Next()
	}
}