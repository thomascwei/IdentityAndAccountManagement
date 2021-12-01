package internal

import (
	"github.com/gin-gonic/gin"
)

//驗證token, 通過後向後傳遞auth
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 找出header裡的token
		token := ""
		for k, v := range c.Request.Header {
			if k == "Authorization" {
				token = v[0]
			}
		}
		// 確認token有效並得到此token的auth
		id, auth, err := Tokenverify(token)
		if err != nil {
			ReturnError(c, err.Error())
			c.Abort()
		}
		c.Set("auth", auth)
		c.Set("id", id)
		c.Next()
	}
}
