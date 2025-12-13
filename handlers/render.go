package handlers

import (
	"ginWeb/common"
	"strings"

	"github.com/gin-gonic/gin"
)

func MyHTMLRender(c *gin.Context, code int, fileName string, data gin.H) {
	if data == nil {
		data = gin.H{}
	}

	var keyVal any
	var exists bool

	if keyVal, exists = c.Get(common.SessionUserIdKey); exists {
		data["UserId"] = keyVal.(string)
	}

	if keyVal, exists = c.Get(common.SessionUserNicknameKey); exists {
		data["UserNickname"] = keyVal.(string)
	}

	currentPath := c.Request.URL.Path
	if strings.Contains(currentPath, "/chat") {
		data["Menu"] = "chat"
	} else {
		data["Menu"] = "board"
	}

	c.HTML(code, fileName, data)
}
