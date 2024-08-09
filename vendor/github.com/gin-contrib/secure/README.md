# Secure

[![Build Status](https://travis-ci.org/gin-contrib/secure.svg)](https://travis-ci.org/gin-contrib/secure)
[![codecov](https://codecov.io/gh/gin-contrib/secure/branch/master/graph/badge.svg)](https://codecov.io/gh/gin-contrib/secure)
[![Go Report Card](https://goreportcard.com/badge/github.com/gin-contrib/secure)](https://goreportcard.com/report/github.com/gin-contrib/secure)
[![GoDoc](https://godoc.org/github.com/gin-contrib/secure?status.svg)](https://godoc.org/github.com/gin-contrib/secure)

Secure middleware for [Gin](https://github.com/gin-gonic/gin/) framework.

## Example

See the [example1](example/code1/example.go), [example2](example/code2/example.go).

DefaultConfig returns a Configuration with strict security settings

[embedmd]:# (secure.go go /func DefaultConfig/ /^}$/)
```go
func DefaultConfig() Config {
	return Config{
		SSLRedirect:           true,
		IsDevelopment:         false,
		STSSeconds:            315360000,
		STSIncludeSubdomains:  true,
		FrameDeny:             true,
		ContentTypeNosniff:    true,
		BrowserXssFilter:      true,
		ContentSecurityPolicy: "default-src 'self'",
		IENoOpen:              true,
		SSLProxyHeaders:       map[string]string{"X-Forwarded-Proto": "https"},
	}
}
```

[embedmd]:# (example/code1/example.go go)
```go
package main

import (
	"github.com/gin-contrib/secure"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.Use(secure.New(secure.Config{
		AllowedHosts:          []string{"example.com", "ssl.example.com"},
		SSLRedirect:           true,
		SSLHost:               "ssl.example.com",
		STSSeconds:            315360000,
		STSIncludeSubdomains:  true,
		FrameDeny:             true,
		ContentTypeNosniff:    true,
		BrowserXssFilter:      true,
		ContentSecurityPolicy: "default-src 'self'",
		IENoOpen:              true,
		ReferrerPolicy:        "strict-origin-when-cross-origin",
		SSLProxyHeaders:       map[string]string{"X-Forwarded-Proto": "https"},
	}))

	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	// Listen and Server in 0.0.0.0:8080
	router.Run()
}
```
