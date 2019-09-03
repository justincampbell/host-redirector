package main

import (
	"log"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

// RedirectURL is the base URL to redirect to that was set.
var RedirectURL *url.URL

func main() {
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	r.POST("/set", func(c *gin.Context) {
		parsed, err := url.Parse(c.Query("url"))
		if err != nil {
			log.Printf("[ERROR] Error parsing URL: %s", err)
			c.AbortWithError(http.StatusBadRequest, err)
		}

		RedirectURL = parsed
		log.Printf("[INFO] Redirect URL set to %q", RedirectURL)

		c.String(http.StatusCreated, "OK")
	})

	r.GET("/redirect", func(c *gin.Context) {
		location, err := url.Parse(RedirectURL.String())
		if err != nil {
			log.Printf("[ERROR] Error parsing redirect URL: %s", err)
			c.AbortWithError(http.StatusBadRequest, err)
		}

		location.RawQuery = c.Request.URL.RawQuery

		log.Printf("[INFO] Redirecting to %q", location)
		c.Redirect(http.StatusFound, location.String())
	})

	r.Run()
}
