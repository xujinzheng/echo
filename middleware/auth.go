package middleware

import (
	"encoding/base64"
	"github.com/labstack/echo"
	"net/http"
)

type (
	AuthFunc func(string, string) bool
)

const (
	Basic = "Basic"
)

// BasicAuth returns an HTTP basic authentication middleware.
func BasicAuth(fn AuthFunc) echo.HandlerFunc {
	return func(c *echo.Context) (he *echo.HTTPError) {
		auth := c.Request.Header.Get(echo.Authorization)
		i := 0
		he = &echo.HTTPError{Code: http.StatusUnauthorized}

		for ; i < len(auth); i++ {
			c := auth[i]
			// Ignore empty spaces
			if c == ' ' {
				continue
			}

			// Check scheme
			if i < len(Basic) {
				// Ignore case
				if i == 0 {
					if c != Basic[i] && c != 'b' {
						return
					}
				} else {
					if c != Basic[i] {
						return
					}
				}
			} else {
				// Extract credentials
				b, err := base64.StdEncoding.DecodeString(auth[i:])
				if err != nil {
					return
				}
				cred := string(b)
				for i := 0; i < len(cred); i++ {
					if cred[i] == ':' {
						// Verify credentials
						if !fn(cred[:i], cred[i+1:]) {
							return
						}
						return nil
					}
				}
			}
		}
		return
	}
}
