package main

import (
	"fmt"
	"ginWeb/chat"
	"ginWeb/common"
	"ginWeb/db"
	"ginWeb/handlers"
	"ginWeb/middlewares"
	"text/template"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	db.InitializeDB()

	r := gin.Default()
	r.SetFuncMap(template.FuncMap{
		"formatDate": func(t time.Time) string {
			return t.Format(time.DateTime)
		},
		"prettyTime": func(t time.Time) string {
			duration := time.Since(t)

			if duration.Minutes() < 1 {
				return "Just now"
			}
			if duration.Hours() < 1 {
				return fmt.Sprintf("%d mins ago", int(duration.Minutes()))
			}
			if duration.Hours() < 24 {
				return fmt.Sprintf("%d hours ago", int(duration.Hours()))
			}

			return t.Format(time.DateTime)
		},
	})
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")

	store := cookie.NewStore([]byte("!@%%#!@#$!@" /*대칭키*/))
	r.Use(sessions.Sessions("sessionCookie" /*cookieName*/, store))
	r.NoRoute(func(c *gin.Context) { handlers.GoHome(c) })

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
		authGroup.GET(common.PostDetailEndpoint+"/:id", handlers.PostDetailHandler())
		authGroup.POST(common.PostDeleteEndpoint, handlers.PostDeleteHandler())
		authGroup.POST(common.PostCommentsCreateEndpoint, handlers.PostCommentsCreateHandler())
		authGroup.POST(common.PostCommentsDeleteEndpoint, handlers.PostCommentsDeleteHandler())

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
