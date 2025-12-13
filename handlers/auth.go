package handlers

import (
	"ginWeb/common"
	"ginWeb/db"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func SignupFormHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		MyHTMLRender(c, http.StatusOK, "signup.html", nil)
	}
}

func SignupRegisterHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 1. 입력받은 데이터
		userId := c.PostForm("userId")
		password := c.PostForm("password")
		username := c.PostForm("username")

		// 2. DB에 저장할 User 객체 생성
		newUser := db.User{
			UserId:   userId,
			Password: password,
			Nickname: username,
		}

		// 3. DB 저장 (Create)
		if err := db.GetDB(db.MainDB).Create(&newUser).Error; err != nil {
			MyHTMLRender(c, http.StatusInternalServerError, "signup.html", gin.H{})
			return
		}

		// 가입 성공 시 로그인 페이지로 이동
		c.Redirect(http.StatusFound, "/login")
	}
}

func LoginFormHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		MyHTMLRender(c, http.StatusOK, "login.html", nil)
	}
}

func LoginHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		inputIds := c.PostForm("userId")
		inputPw := c.PostForm("password")

		var user db.User

		if err := db.GetDB(db.MainDB).Where("user_id = ? AND password = ?", inputIds, inputPw).First(&user).Error; err != nil {
			// 조회 실패 (유저가 없거나 비번이 틀림)
			MyHTMLRender(c, http.StatusInternalServerError, "login.html", gin.H{
				"error": "아이디 또는 비밀번호가 잘못되었습니다.",
			})
			return
		}

		// 로그인 성공 처리
		session := sessions.Default(c)
		session.Set(common.SessionUserIdKey, user.UserId) // 세션에는 유저의 고유 ID를 넣음
		session.Set(common.SessionUserNicknameKey, user.Nickname)
		session.Save()

		c.Redirect(http.StatusFound, common.HomeEndpoint) // "/"
	}
}
