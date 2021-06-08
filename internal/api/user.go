package api

import (
	"cloudpan/internal/model"
	"cloudpan/internal/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type userRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Register 处理用户注册请求
func Register(c *gin.Context) {
	var req userRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		util.Log().Warning("bindjson error", err)
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Wrong request format!"})
		return
	}

	// 检查用户名是否重复
	user := model.GetUserByName(req.Username)
	if user.ID > 0 {
		util.Log().Warning(fmt.Sprintf("Username [%s] duplicated!", req.Username))
		c.JSON(http.StatusAccepted, gin.H{"msg": "Username has been used!"})
		return
	}

	// 向数据库中创建新用户
	user = model.NewUser()
	user.Username = req.Username
	user.Password = util.MD5([]byte(req.Password))
	user.Key = util.GenerateRandomBytes(16)

	if err = model.CreateUser(&user); err != nil {
		util.Log().Warning("create user error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Create user error!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "Register success!"})
}

// Login 处理用户登陆请求
func Login(c *gin.Context) {
	var req userRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		util.Log().Warning("bindjson error", err)
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Wrong request format!"})
		return
	}

	// 检查用户名密码是否正确
	user := model.GetUserByName(req.Username)
	if user.Username == "" {
		c.JSON(http.StatusAccepted, gin.H{"msg": "Wrong username or password!"})
		return
	}
	if user.Password != util.MD5([]byte(req.Password)) {
		util.Log().Info(fmt.Sprintf("[%s] [%s] login failed", req.Username, req.Password))
		c.JSON(http.StatusAccepted, gin.H{"msg": "Wrong username or password!"})
		return
	}

	// 设置session
	session, _ := model.GetSession(c.Request)
	session.Values["user"] = user
	if err = session.Save(c.Request, c.Writer); err != nil {
		util.Log().Warning("save session error", err)
		c.JSON(http.StatusAccepted, gin.H{"msg": "Save session error!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "Login success!"})
}

// Logout 处理用户注销请求
func Logout(c *gin.Context) {
	session, _ := model.GetSession(c.Request)

	user := session.Values["user"].(*model.User)

	session.Options.MaxAge = -1
	err := session.Save(c.Request, c.Writer)
	if err != nil {
		util.Log().Warning("save session error", err)
		c.JSON(http.StatusAccepted, gin.H{"msg": "Save session error!"})
		return
	}

	util.Log().Info(fmt.Sprintf("[%s] logout success!", user.Username))
	c.JSON(http.StatusOK, gin.H{"msg": "Logout success!"})
}
