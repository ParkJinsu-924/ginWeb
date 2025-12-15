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

	if val, exist := c.Get(common.SessionUserUUIDKey); exist {
		data["UserIndex"] = val.(uint)
	}
	if val, exist := c.Get(common.SessionUserIdKey); exist {
		data["UserId"] = val.(string)
	}
	if val, exist := c.Get(common.SessionUserNicknameKey); exist {
		data["UserNickname"] = val.(string)
	}
	if val, exist := c.Get(common.SessionUserTagKey); exist {
		data["UserTag"] = val.(uint64)
	}

	currentPath := c.Request.URL.Path
	if strings.Contains(currentPath, "/chat") {
		data["Menu"] = "chat"
	} else {
		data["Menu"] = "board"
	}

	c.HTML(code, fileName, data)
}
