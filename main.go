package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/justincampbell/oauth-redirector/redir"
)

var (
	// RedirectURL is the base URL to redirect to that was set.
	RedirectURL *url.URL

	// Token is the auth token required to set the RedirectURL.
	Token string
)

func main() {
	Token = os.Getenv("OAUTH_REDIRECTOR_TOKEN")
	if Token == "" {
		log.Fatalf("OAUTH_REDIRECTOR_TOKEN not set")
	}

	gin.SetMode(gin.ReleaseMode)
	r := setupRouter()
	r.Run()
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/set", func(c *gin.Context) {
		authorization := c.Request.Header.Get("Authorization")
		if authorization != fmt.Sprintf("Bearer %s", Token) {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		var json redir.SetRequest
		if err := c.ShouldBindJSON(&json); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		parsed, err := url.Parse(json.URL)
		if err != nil {
			log.Printf("[ERROR] Error parsing URL: %s", err)
			c.AbortWithError(http.StatusBadRequest, err)
		}

		RedirectURL = parsed
		log.Printf("[INFO] Redirect URL set to %q", RedirectURL)

		c.String(http.StatusCreated, "OK")
	})

	r.GET("/redirect", func(c *gin.Context) {
		if RedirectURL == nil {
			log.Printf("[WARN] Redirect attempted when URL was not set")
			c.AbortWithStatus(http.StatusNotFound)
		}

		location, err := url.Parse(RedirectURL.String())
		if err != nil {
			log.Printf("[ERROR] Error parsing redirect URL: %s", err)
			c.AbortWithError(http.StatusBadRequest, err)
		}

		location.RawQuery = c.Request.URL.RawQuery

		log.Printf("[INFO] Redirecting to %q", location)
		c.Redirect(http.StatusFound, location.String())
	})

	return r
}
