package handlers

import (
	"ginWeb/common"

	"github.com/gin-gonic/gin"
)

func MyHTMLRender(c *gin.Context, code int, fileName string, data gin.H) {
	if data == nil {
		data = gin.H{}
	}

	if user, exists := c.Get(common.SessionUserNicknameKey); exists {
		data["User"] = user.(string)
	}

	c.HTML(code, fileName, data)
}
