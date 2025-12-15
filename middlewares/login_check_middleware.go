package middlewares

import (
	"ginWeb/common"
	"ginWeb/handlers"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func LoginCheckMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userId := session.Get(common.SessionUserIdKey)

		if userId == nil {
			c.Redirect(http.StatusFound, common.LoginEndpoint)
			c.Abort()
			return
		}

		// 컨텍스트에 유저 정보 세팅
		c.Set(common.SessionUserUUIDKey, session.Get(common.SessionUserUUIDKey))
		c.Set(common.SessionUserIdKey, userId)
		c.Set(common.SessionUserNicknameKey, session.Get(common.SessionUserNicknameKey))
		c.Set(common.SessionUserTagKey, session.Get(common.SessionUserTagKey))
		c.Next()
	}
}

func NotLoginCheckMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userId := session.Get(common.SessionUserIdKey)

		if userId != nil {
			handlers.GoHome(c)
			c.Abort()
			return
		}

		c.Next()
	}
}
