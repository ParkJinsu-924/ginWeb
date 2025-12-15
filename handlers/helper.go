package handlers

import (
	"fmt"
	"ginWeb/common"
	"ginWeb/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getPostFromDB(postId string) (*db.Post, bool) {
	var post db.Post
	tx := db.GetDB(db.MainDB).First(&post, postId)
	if tx.Error != nil {
		fmt.Println("getPostFromDB Error: ", tx.Error)
		return &post, false
	}

	return &post, true
}

func getUserDataFromDB(userIndex uint) (*db.User, bool) {
	var user db.User
	tx := db.GetDB(db.MainDB).First(&user, userIndex)
	if tx.Error != nil {
		fmt.Println("getUserDataFromDB Error: ", tx.Error)
		return &user, false
	}

	return &user, true
}

func getUserDataListFromDB(userIndexList []uint) (map[uint]db.User, bool) {
	userMap := make(map[uint]db.User)
	if len(userIndexList) == 0 {
		return userMap, true
	}

	var users []db.User

	tx := db.GetDB(db.MainDB).Find(&users, userIndexList)
	if tx.Error != nil {
		fmt.Println("getUserDataFromDB Error: ", tx.Error)
		return userMap, false
	}

	for _, user := range users {
		userMap[user.ID] = user
	}

	return userMap, true
}

func getCommentsFromDB(postId string) (*[]db.Comment, bool) {
	var comment []db.Comment
	tx := db.GetDB(db.MainDB).Where("post_index = ?", postId).Order("created_at desc").Find(&comment)
	if tx.Error != nil {
		fmt.Println("getCommentFromDB Error: ", tx.Error)
		return &comment, false
	}

	return &comment, true
}

func GoHome(c *gin.Context) {
	c.Redirect(http.StatusFound, common.HomeEndpoint)
}
