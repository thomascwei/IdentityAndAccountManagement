package api

import (
	"IAM/internal"
	"fmt"
	"github.com/gin-gonic/gin"
)

type User struct {
	username string `json:"username"`
	password string `json:"password"`
	email    string `json:"email"`
	auth     int    `json:"auth"`
}

func SignUp(c *gin.Context) (int64, error) {
	for k, v := range c.Request.Header {
		fmt.Println(k, v)
	}
	input := User{}
	c.BindJSON(&input)
	id, err := internal.SignUp(input.username, input.password, input.email, input.auth)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func Login(c *gin.Context) {
	var m map[string]interface{}
	err := c.Bind(&m)
	if err != nil {
		c.JSON(200, gin.H{
			"result": "fail",
			"error":  err.Error(),
		})
		return
	}

	token, err := internal.Login(m["username"].(string), m["password"].(string))
	if err != nil {
		c.JSON(200, gin.H{
			"result": "fail",
			"error":  err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"result": "ok",
		"token":  token,
	})
	return
}
