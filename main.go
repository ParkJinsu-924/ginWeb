package main

import (
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

	r.GET(common.LoginEndpoint, handlers.LoginFormHandler())
	r.POST(common.LoginEndpoint, handlers.LoginHandler())
	r.GET(common.SignupEndpoint, handlers.SignupFormHandler())
	r.POST(common.SignupEndpoint, handlers.SignupRegisterHandler())

	authGroup := r.Group("/")
	authGroup.Use(middlewares.LoginCheckMiddleware())
	{
		authGroup.GET(common.HomeEndpoint, handlers.HomeHandler())
		authGroup.GET(common.LogoutEndpoint, handlers.LogoutHandler())
		authGroup.GET(common.PostFormEndpoint, handlers.PostFormHandler())
		authGroup.POST(common.PostCreateEndpoint, handlers.PostCreateHandler())
		authGroup.GET(common.PostDetailEndpoint, handlers.PostDetailHandler())
		authGroup.POST(common.PostDeleteEndpoint, handlers.PostDeleteHandler())
	}

	_ = r.Run(":8080")
}
