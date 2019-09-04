package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRedirect(t *testing.T) {
	tests := []struct {
		name             string
		redirectURL      string
		callbackParams   string
		expectedStatus   int
		expectedLocation string
	}{
		{
			"not set",
			"",
			"foo=bar",
			http.StatusNotFound,
			"",
		},
		{
			"set https",
			"https://example.com/a/b/c",
			"foo=bar",
			http.StatusFound,
			"https://example.com/a/b/c?foo=bar",
		},
		{
			"set http + port",
			"http://localhost:1234/a/b/c",
			"foo=bar",
			http.StatusFound,
			"http://localhost:1234/a/b/c?foo=bar",
		},
	}

	router := setupRouter()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.redirectURL != "" {
				parsed, err := url.Parse(tt.redirectURL)
				require.NoError(t, err)
				RedirectURL = parsed
			} else {
				RedirectURL = nil
			}

			w := httptest.NewRecorder()
			reqURL := fmt.Sprintf("/redirect?%s", tt.callbackParams)
			req, _ := http.NewRequest("GET", reqURL, nil)
			router.ServeHTTP(w, req)

			require.Equal(t, tt.expectedStatus, w.Code)
			assert.Equal(t, tt.expectedLocation, w.Header().Get("Location"))
		})
	}

}
