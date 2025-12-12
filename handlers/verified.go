package handlers

import (
	"ginWeb/common"
	"ginWeb/db"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func HomeHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		val, exist := c.Get(common.SessionUserIdKey)
		if !exist {
			c.Abort()
			return
		}

		var posts []db.Post

		tx := db.GetDB(db.MainDB).Order("id desc").Find(&posts)

		if tx.Error != nil {
			panic(tx.Error)
		}

		c.HTML(http.StatusOK, "index.html", gin.H{
			"User":         val,
			"Posts":        posts,
			"PostFormPath": common.PostFormEndpoint,
		})
	}
}

func LogoutHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		session.Clear()
		session.Options(sessions.Options{
			Path:   "/",
			MaxAge: -1,
		})
		session.Save()
		c.Redirect(http.StatusFound, common.LoginEndpoint)
	}
}

func PostFormHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "post_form.html", gin.H{
			"PostCreatePath": common.PostCreateEndpoint,
		})
	}
}

func PostCreateHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		authorUserId := session.Get(common.SessionUserIdKey).(string)
		title := c.PostForm("title")
		content := c.PostForm("content")

		//time
		newPost := db.Post{
			UserId:           authorUserId,
			Title:            title,
			Content:          content,
			CreatedTimestamp: time.Now().Format("2006-01-02 15:04:05"),
		}

		if err := db.GetDB(db.MainDB).Create(&newPost).Error; err != nil {
			panic(err)
		}

		c.Redirect(http.StatusFound, common.HomeEndpoint)
	}
}

func PostDetailHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		postId := c.Param("id")

		var post db.Post
		tx := db.GetDB(db.MainDB).First(&post, postId)
		if tx.Error != nil {
			panic(tx.Error)
		}

		c.HTML(http.StatusOK, "post_detail.html", gin.H{
			"Post":           post,
			"DeletePostPath": common.PostDeleteEndpoint,
		})
	}
}

func PostDeleteHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		postId := c.PostForm("id")
		if postId == "" {
			c.Redirect(http.StatusFound, common.HomeEndpoint)
			return
		}

		tx := db.GetDB(db.MainDB).Delete(&db.Post{}, postId)
		if tx.Error != nil {
			panic(tx.Error)
		}

		c.Redirect(http.StatusFound, common.HomeEndpoint)
	}
}
