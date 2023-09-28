package middleware

import (
	t "github.com/chdlvy/go-chatFrame/jwt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 验证是否登录过、是否有token
func VerifyTokenMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHandler := c.Request.Header.Get("authorization")
		if authHandler == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "请求头部auth为空，登录过期，重新登录"})
		} else {
			_, err := t.ParseToken(authHandler)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"msg": "token失效，请重新登录",
				})
				return
			}
			// 刷新token
			newToken := refushToken(authHandler)
			c.AbortWithStatusJSON(http.StatusOK, gin.H{
				"token": newToken,
			})
			// c.Next()
		}

	}
}

// 根据旧token返回一个新token
func refushToken(token string) string {
	msg, _ := t.ParseToken(token)
	phone := msg["phone"]
	password := msg["password"]
	m := map[string]interface{}{
		"phone":    phone,
		"password": password,
	}
	newToken := t.GetToken(m)
	return newToken
}

// 跨域
func Cors() gin.HandlerFunc {
	return func(context *gin.Context) {
		method := context.Request.Method

		context.Header("Access-Control-Allow-Origin", "*")
		context.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token, x-token")
		context.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PATCH, PUT")
		context.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		context.Header("Access-Control-Allow-Credentials", "true")

		if method == "OPTIONS" {
			context.AbortWithStatus(http.StatusNoContent)
		}
	}
}
