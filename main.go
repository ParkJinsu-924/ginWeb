package main

import (
	"ginWeb/chat"
	"ginWeb/common"
	"ginWeb/db"
	"ginWeb/handlers"
	"ginWeb/middlewares"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	db.InitializeDB()

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")

	store := cookie.NewStore([]byte("!@%%#!@#$!@" /*대칭키*/))
	r.Use(sessions.Sessions("sessionCookie" /*cookieName*/, store))
	r.NoRoute(func(c *gin.Context) { c.Redirect(http.StatusFound, common.HomeEndpoint) })

	wsHub := chat.NewHub()
	go wsHub.Run()

	frontGroup := r.Group("/")
	frontGroup.Use(middlewares.NotLoginCheckMiddleware())
	{
		frontGroup.GET(common.LoginEndpoint, handlers.LoginFormHandler())
		frontGroup.POST(common.LoginEndpoint, handlers.LoginHandler())
		frontGroup.GET(common.SignupEndpoint, handlers.SignupFormHandler())
		frontGroup.POST(common.SignupEndpoint, handlers.SignupRegisterHandler())
	}

	authGroup := r.Group("/")
	authGroup.Use(middlewares.LoginCheckMiddleware())
	{
		authGroup.GET(common.HomeEndpoint, handlers.HomeHandler())
		authGroup.GET(common.LogoutEndpoint, handlers.LogoutHandler())

		// posting
		authGroup.GET(common.PostFormEndpoint, handlers.PostFormHandler())
		authGroup.POST(common.PostCreateEndpoint, handlers.PostCreateHandler())
		authGroup.GET(common.PostDetailEndpoint, handlers.PostDetailHandler())
		authGroup.POST(common.PostDeleteEndpoint, handlers.PostDeleteHandler())

		// chat
		authGroup.GET(common.ChatEndpoint, handlers.ChatPageHandler())
		authGroup.GET("/ws", func(c *gin.Context) {
			session := sessions.Default(c)
			nickname := session.Get(common.SessionUserNicknameKey).(string)
			chat.ServeWs(wsHub, c, nickname)
		})
	}

	_ = r.Run(":8080")
}
