package auth

import (
	global "gin-blog/internal/global"
	"gin-blog/internal/model/dto/request"
	"gin-blog/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	svc service.AuthService
}

func NewAuthController(svc service.AuthService) *AuthController {
	return &AuthController{svc: svc}
}

func (ctrl *AuthController) Login(c *gin.Context) {
	var req request.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		global.ReturnError(c, global.ErrRequest, err)
		return
	}

	vo, err := ctrl.svc.Login(c, req)
	if err != nil {
		global.ReturnError(c, err, nil)
		return
	}

	session := sessions.Default(c)
	session.Set(global.CTX_USER_AUTH, vo.ID)
	session.Save()

	global.ReturnSuccess(c, vo)
}

func (ctrl *AuthController) Logout(c *gin.Context) {
	c.Set(global.CTX_USER_AUTH, nil)

	authIdStr := sessions.Default(c).Get(global.CTX_USER_AUTH)
	if authIdStr == nil {
		global.ReturnSuccess(c, nil)
		return
	}

	authId, _ := strconv.Atoi(authIdStr.(string)) // Or handle interface{} properly

	session := sessions.Default(c)
	session.Delete(global.CTX_USER_AUTH)
	session.Save()

	ctrl.svc.Logout(c, authId)

	global.ReturnSuccess(c, nil)
}

func (ctrl *AuthController) Register(c *gin.Context) {
	var req request.RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		global.ReturnError(c, global.ErrRequest, err)
		return
	}

	if err := ctrl.svc.Register(c, req); err != nil {
		global.ReturnError(c, err, nil)
		return
	}

	global.ReturnSuccess(c, nil)
}

func (ctrl *AuthController) SendCode(c *gin.Context) {
	var req request.SendCodeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		global.ReturnError(c, global.ErrRequest, err)
		return
	}

	if err := ctrl.svc.SendCode(c, req.Email); err != nil {
		global.ReturnError(c, err, nil)
		return
	}

	global.ReturnSuccess(c, nil)
}

func (ctrl *AuthController) VerifyCode(c *gin.Context) {
	code := c.Query("info")
	if code == "" {
		ctrl.returnErrorPage(c)
		return
	}

	if err := ctrl.svc.VerifyCode(c, code); err != nil {
		ctrl.returnErrorPage(c)
		return
	}

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(`
        <!DOCTYPE html>
        <html lang="zh-CN">
        <head>
            <meta charset="UTF-8">
            <meta name="viewport" content="width=device-width, initial-scale=1.0">
            <title>注册成功</title>
            <style>
                body {
                    font-family: Arial, sans-serif;
                    background-color: #f4f4f4;
                    display: flex;
                    justify-content: center;
                    align-items: center;
                    height: 100vh;
                    margin: 0;
                }
                .container {
                    background-color: #fff;
                    padding: 20px;
                    border-radius: 8px;
                    box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
                    text-align: center;
                }
                h1 {
                    color: #5cb85c;
                }
                p {
                    color: #333;
                }
            </style>
        </head>
        <body>
            <div class="container">
                <h1>注册成功</h1>
                <p>恭喜您，注册成功！</p>
            </div>
        </body>
        </html>
    `))
}

func (ctrl *AuthController) returnErrorPage(c *gin.Context) {
	c.Data(http.StatusInternalServerError, "text/html; charset=utf-8", []byte(`
        <!DOCTYPE html>
        <html lang="zh-CN">
        <head>
            <meta charset="UTF-8">
            <meta name="viewport" content="width=device-width, initial-scale=1.0">
            <title>注册失败</title>
            <style>
                body {
                    font-family: Arial, sans-serif;
                    background-color: #f4f4f4;
                    display: flex;
                    justify-content: center;
                    align-items: center;
                    height: 100vh;
                    margin: 0;
                }
                .container {
                    background-color: #fff;
                    padding: 20px;
                    border-radius: 8px;
                    box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
                    text-align: center;
                }
                h1 {
                    color: #d9534f;
                }
                p {
                    color: #333;
                }
            </style>
        </head>
        <body>
            <div class="container">
                <h1>注册失败</h1>
                <p>请重试。</p>
            </div>
        </body>
        </html>
    `))
}
