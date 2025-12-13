package handlers

import (
	"fmt"
	"ginWeb/common"
	"ginWeb/db"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func HomeHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		var posts []db.Post

		tx := db.GetDB(db.MainDB).Order("id desc").Find(&posts)

		if tx.Error != nil {
			fmt.Println("HomeHandler Error: ", tx.Error)
			GoHome(c)
			return
		}

		MyHTMLRender(c, http.StatusOK, "index.html", gin.H{
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
		MyHTMLRender(c, http.StatusOK, "post_form.html", gin.H{
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
			fmt.Println("PostCreateHandler Error: ", err)
			GoHome(c)
			return
		}

		GoHome(c)
	}
}

func PostDetailHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		postId := c.Param("id")

		post, res := getPostFromDB(postId)
		if !res {
			GoHome(c)
			return
		}

		MyHTMLRender(c, http.StatusOK, "post_detail.html", gin.H{
			"Post":           post,
			"DeletePostPath": common.PostDeleteEndpoint,
		})
	}
}

func PostDeleteHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		postId := c.PostForm("id")
		if postId == "" {
			GoHome(c)
			return
		}

		post, res := getPostFromDB(postId)
		if !res {
			GoHome(c)
			return
		}

		currentUserId := c.MustGet(common.SessionUserIdKey).(string)

		if post.UserId != currentUserId {
			GoHome(c)
			return
		}

		if err := db.GetDB(db.MainDB).Delete(&post).Error; err != nil {
			fmt.Println("Delete Error:", err)
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{"Message": "삭제 중 오류 발생"})
			return
		}

		GoHome(c)
	}
}
