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

		tx := db.GetDB(db.MainDB).InnerJoins("User").Order("posts.id desc").Find(&posts)

		if tx.Error != nil {
			fmt.Println("HomeHandler Error: ", tx.Error)
			c.Abort()
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
		authorUserUid := session.Get(common.SessionUserUUIDKey).(uint)
		title := c.PostForm("title")
		content := c.PostForm("content")

		//time
		newPost := db.Post{
			UserIndex:        authorUserUid,
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

		comments, res := getCommentsFromDB(postId)
		if !res {
			GoHome(c)
			return
		}

		var userIdList []uint
		userIdList = append(userIdList, post.UserIndex)

		for _, comment := range *comments {
			userIdList = append(userIdList, comment.UserIndex)
		}

		userMap, res := getUserDataListFromDB(userIdList)
		if !res {
			GoHome(c)
			return
		}

		authorUser := userMap[post.UserIndex]

		type commentData struct {
			Comment db.Comment
			User    db.User
		}

		var commentDataList []commentData

		for _, comment := range *comments {
			commentUser := userMap[comment.UserIndex]
			commentDataList = append(commentDataList, commentData{comment, commentUser})
		}

		MyHTMLRender(c, http.StatusOK, "post_detail.html", gin.H{
			"Post":               post,
			"PostAuthorNickname": authorUser.Nickname,
			"PostAuthorTag":      authorUser.Tag,
			"CommentsCreatePath": common.PostCommentsCreateEndpoint,
			"DeletePostPath":     common.PostDeleteEndpoint,
			"Comments":           commentDataList,
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

		currentUserUid := c.MustGet(common.SessionUserUUIDKey).(uint)

		if post.UserIndex != currentUserUid {
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

func PostCommentsCreateHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		postId := c.PostForm("post_id")
		content := c.PostForm("content")
		if postId == "" || content == "" {
			GoHome(c)
			return
		}

		post, res := getPostFromDB(postId)
		if !res {
			GoHome(c)
			return
		}

		currentUserUid := c.MustGet(common.SessionUserUUIDKey).(uint)

		newComment := db.Comment{
			PostIndex: post.ID,
			UserIndex: currentUserUid,
			Content:   content,
		}

		if err := db.GetDB(db.MainDB).Create(&newComment).Error; err != nil {
			fmt.Println("Comment Create Error: ", err)
			GoHome(c)
			return
		}

		c.Redirect(http.StatusFound, fmt.Sprintf(common.PostDetailEndpoint+"/%d", post.ID))
	}
}
