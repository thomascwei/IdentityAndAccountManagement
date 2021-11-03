package api

import (
	"IAM/internal"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Auth     int    `json:"auth"`
}

type login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type YamlConfig struct {
	API struct {
		SignUp int `yaml:"signUp"`
		Logout int `yaml:"logout"`
	} `yaml:"api"`
}

var CC YamlConfig
var CConfig = CC.getConf()

func (y *YamlConfig) getConf() *YamlConfig {
	yamlFile, err := ioutil.ReadFile("api/config.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, y)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return y
}

// 建立帳戶
func SignUp(c *gin.Context) {

	// 找出token
	token := ""
	for k, v := range c.Request.Header {
		if k == "Authorization" {
			token = v[0]
		}
	}
	// 找出此token的Auth
	auth, err := internal.Tokenverify(token)
	if err != nil {
		c.JSON(200, gin.H{
			"result": "fail",
			"error":  "token error, " + err.Error(),
		})
		return
	}
	// 確認此token有使用此API的權限
	if auth < CConfig.API.SignUp {
		c.JSON(200, gin.H{
			"result": "fail",
			"error":  "not authorized",
		})
		return
	}

	// 讀request body
	input := User{}
	err = c.BindJSON(&input)
	if err != nil {
		c.JSON(200, gin.H{
			"result": "fail",
			"error":  err.Error(),
		})
		return
	}
	// 確認新帳號權限小於自身帳號
	if auth < input.Auth {
		c.JSON(200, gin.H{
			"result": "fail",
			"error":  "you cannot create authority higher than your self",
		})
		return
	}
	// 建立帳號
	id, err := internal.SignUp(input.Username, input.Password, input.Email, input.Auth)
	if err != nil {
		c.JSON(200, gin.H{
			"result": "fail",
			"error":  err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"result": "ok",
		"id":     id,
	})
	return
}

func Login(c *gin.Context) {

	m := login{}
	err := c.Bind(&m)
	if err != nil {
		c.JSON(200, gin.H{
			"result": "fail",
			"error":  err.Error(),
		})
		return
	}
	if m.Username == "" || m.Password == "" {
		c.JSON(200, gin.H{
			"result": "fail",
			"error":  "request body format error",
		})
		return
	}
	token, err := internal.Login(m.Username, m.Password)
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

func Logout(c *gin.Context) {
	// 找出token
	token := ""
	for k, v := range c.Request.Header {
		if k == "Authorization" {
			token = v[0]
		}
	}
	// 刪除token, 不管是否存在
	internal.Logout(token)
	c.JSON(200, gin.H{
		"result": "ok",
	})

	return
}
