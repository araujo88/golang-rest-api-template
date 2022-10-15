package middleware

import (
	"github.com/gin-contrib/secure"
	"github.com/gin-gonic/gin"
)

func SecurityMiddleware() gin.HandlerFunc {
	return secure.New(secure.Config{
		//AllowedHosts:          []string{"example.com", "ssl.example.com"},
		//SSLRedirect:           true,
		//SSLHost:               "ssl.example.com",
		STSSeconds:            315360000,
		STSIncludeSubdomains:  true,
		FrameDeny:             true,
		ContentTypeNosniff:    true,
		BrowserXssFilter:      true,
		ContentSecurityPolicy: "default-src 'self'",
		IENoOpen:              true,
		ReferrerPolicy:        "strict-origin-when-cross-origin",
		SSLProxyHeaders:       map[string]string{"X-Forwarded-Proto": "https"},
	})
}
