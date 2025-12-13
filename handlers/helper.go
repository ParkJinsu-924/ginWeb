package handlers

import (
	"fmt"
	"ginWeb/common"
	"ginWeb/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getPostFromDB(postId string) (db.Post, bool) {
	var post db.Post
	tx := db.GetDB(db.MainDB).First(&post, postId)
	if tx.Error != nil {
		fmt.Println("getPostFromDB Error: ", tx.Error)
		return post, false
	}

	return post, true
}

func GoHome(c *gin.Context) {
	c.Redirect(http.StatusFound, common.HomeEndpoint)
}
