package service

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"landlord/model"
	"net/http"
)

// Register 注册账号
func Register(c *gin.Context) {
	var request model.RegisterRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		// 请求数据格式错误
		c.JSON(http.StatusBadRequest, BadRequestError{err.Error()})
		return
	}
	response, err := request.Register()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response)
	} else {
		c.JSON(http.StatusOK, response)
	}
}

// SendValidateCode 发送验证码
func SendValidateCode(c *gin.Context) {
	var request model.SendValidateCodeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, BadRequestError{err.Error()})
		return
	}
	response, err := request.SendValidateCode()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response)
	} else {
		c.JSON(http.StatusOK, response)
	}
}

// Login 登录
func Login(c *gin.Context) {
	var request model.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, BadRequestError{err.Error()})
		return
	}
	response, err := request.Login()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	// 设置session
	session := sessions.Default(c)
	session.Set("user_id", response.UserID)
	session.Set("nickname", response.Nickname)
	c.JSON(http.StatusOK, response)
}