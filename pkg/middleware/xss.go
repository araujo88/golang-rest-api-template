package middleware

import (
	xss "github.com/araujo88/gin-gonic-xss-middleware"
	"github.com/gin-gonic/gin"
)

func Xss() gin.HandlerFunc {
	var xssMdlwr xss.XssMw
	return xssMdlwr.RemoveXss()
}
