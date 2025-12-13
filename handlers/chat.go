package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ChatPageHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		MyHTMLRender(c, http.StatusOK, "chat.html", gin.H{})
	}
}
