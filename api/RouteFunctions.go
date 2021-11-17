package api

import (
	"IAM/internal"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Auth     int    `json:"auth"`
}

type changePassword struct {
	Id          int    `json:"id"`
	NewPassword string `json:"new_password"`
}

type login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type YamlConfig struct {
	API struct {
		SignUp        int `yaml:"signUp"`
		Logout        int `yaml:"logout"`
		GetAllAccount int `yaml:"getAllAccount"`
		Update        int `yaml:"update"`
		InitPassword  int `yaml:"initpassword"`
	} `yaml:"api"`
}

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

var CC YamlConfig
var CConfig = CC.getConf()
var (
	file, _ = os.OpenFile("./log/RouteFunctions.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	//Trace   = log.New(os.Stdout, "TRACE: ", log.Ldate|log.Ltime|log.Lshortfile)
	Info = log.New(io.MultiWriter(file, os.Stdout), "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	//Error   = log.New(io.MultiWriter(file, os.Stdout), "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
)

func ReturnError(c *gin.Context, info string) {
	Info.Println(info)
	c.JSON(200, gin.H{
		"result": "fail",
		"error":  info,
	})
}

// 建立帳戶
func SignUp(c *gin.Context) {
	// 找出header裡的token
	token := ""
	for k, v := range c.Request.Header {
		if k == "Authorization" {
			token = v[0]
		}
	}
	// 確認token有效並得到此token的auth
	_, auth, err := internal.Tokenverify(token)
	if err != nil {
		c.JSON(200, gin.H{
			"result": "fail",
			"error":  "token error, " + err.Error(),
		})
		return
	}
	// 確認有使用此API的權限
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
	fmt.Println(104, input.Auth)
	if err != nil {
		c.JSON(200, gin.H{
			"result": "fail",
			"error":  err.Error(),
		})
		return
	}
	// 確認新帳號權限小於自身帳號
	if auth <= input.Auth {
		c.JSON(200, gin.H{
			"result": "fail",
			"error":  "you cannot create authority higher than or equal to yourself",
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
	_, _, _ = internal.Tokenverify(token)
	c.JSON(200, gin.H{
		"result": "ok",
		"id":     id,
	})
	return
}

// 登入
func Login(c *gin.Context) {
	m := login{}
	err := c.Bind(&m)
	if err != nil {
		ReturnError(c, err.Error())
		return
	}
	// 檢查request body轉type後是否少欄位
	if m.Username == "" || m.Password == "" {
		ReturnError(c, "request body format error")
		return
	}
	token, err := internal.Login(m.Username, m.Password)
	if err != nil {
		ReturnError(c, err.Error())
		return
	}
	c.JSON(200, gin.H{
		"result": "ok",
		"token":  token,
	})
	Info.Println(m.Username, "login successfully")
	return
}

// 登出
func Logout(c *gin.Context) {
	// 找出token
	token := ""
	for k, v := range c.Request.Header {
		if k == "Authorization" {
			token = v[0]
		}
	}
	// 刪除token, 不管是否存在
	_ = internal.Logout(token)
	c.JSON(200, gin.H{
		"result": "ok",
	})
}

// 取得帳號戶清單及內容
func GetAllAccount(c *gin.Context) {
	// 找出header裡的token
	token := ""
	for k, v := range c.Request.Header {
		if k == "Authorization" {
			token = v[0]
		}
	}
	// 確認token有效並得到此token的auth
	_, auth, err := internal.Tokenverify(token)
	if err != nil {
		ReturnError(c, err.Error())
		return
	}
	// 確認有使用此API的權限
	if auth < CConfig.API.GetAllAccount {
		ReturnError(c, "not authorized")
		return
	}
	result, err := internal.GetAllAccount()
	if err != nil {
		ReturnError(c, err.Error())
		return
	}
	_, _, _ = internal.Tokenverify(token)
	c.JSON(200, gin.H{
		"result":   "ok",
		"accounts": result,
	})
}

// 更新帳戶內容
func AccountUpdate(c *gin.Context) {
	// 找出header裡的token
	token := ""
	for k, v := range c.Request.Header {
		if k == "Authorization" {
			token = v[0]
		}
	}
	// 確認token有效並得到此token的auth
	id, auth, err := internal.Tokenverify(token)
	if err != nil {
		ReturnError(c, err.Error())
		return
	}

	// 確認有使用此API的權限
	if auth < CConfig.API.Update {
		ReturnError(c, "ID:"+strconv.Itoa(id)+" not authorized")
		return
	}
	// 讀request body
	input := User{}
	err = c.BindJSON(&input)
	fmt.Println("user:", input)
	if err != nil {
		ReturnError(c, err.Error())
		return
	}
	// 新權限必須小於此token
	if input.Auth >= auth {
		ReturnError(c, "auth too high")
		return
	}
	// 轉成map
	UserMap := structs.Map(input)
	// 刪除password, 密碼要在其他API單獨改
	delete(UserMap, "Password")
	fmt.Println("UserMap:", UserMap)
	err = internal.UpdateSingelAccount(UserMap)
	if err != nil {
		ReturnError(c, err.Error())
		return
	}
	_, _, _ = internal.Tokenverify(token)
	c.JSON(200, gin.H{
		"result": "ok",
	})
}

// 改自己的密碼
func ChangeSelfPassword(c *gin.Context) {
	// 找出header裡的token
	token := ""
	for k, v := range c.Request.Header {
		if k == "Authorization" {
			token = v[0]
		}
	}
	// 確認token有效並得到此token的auth
	id, _, err := internal.Tokenverify(token)
	if err != nil {
		ReturnError(c, err.Error())
		return
	}
	// 讀request body
	input := changePassword{}
	err = c.BindJSON(&input)
	if err != nil {
		ReturnError(c, err.Error())
		return
	}
	err = internal.ChangeSelfPassword(id, input.NewPassword)
	if err != nil {
		ReturnError(c, err.Error())
		return
	}
	_, _, _ = internal.Tokenverify(token)
	c.JSON(200, gin.H{
		"result": "ok",
	})

}

// 初始化密碼(忘記密碼時使用)
func InitPassword(c *gin.Context) {
	// 找出header裡的token
	token := ""
	for k, v := range c.Request.Header {
		if k == "Authorization" {
			token = v[0]
		}
	}
	// 確認token有效並得到此token的auth
	_, auth, err := internal.Tokenverify(token)
	if err != nil {
		ReturnError(c, "token error, "+err.Error())
		return
	}
	// 確認有使用此API的權限
	if auth < CConfig.API.InitPassword {
		ReturnError(c, "not authorized")
		return
	}
	// 讀request body
	input := changePassword{}
	err = c.BindJSON(&input)
	if err != nil {
		ReturnError(c, err.Error())
		return
	}
	err = internal.InitPassword(input.Id, input.NewPassword)
	if err != nil {
		ReturnError(c, err.Error())
		return
	}
	// 更新token時效
	_, _, _ = internal.Tokenverify(token)
	c.JSON(200, gin.H{
		"result": "ok",
	})

}

// 確認token是否有效
func Tokenverify(c *gin.Context) {
	// 找出header裡的token
	token := ""
	for k, v := range c.Request.Header {
		if k == "Authorization" {
			token = v[0]
		}
	}
	// 確認token有效並得到此token的auth
	_, _, err := internal.Tokenverify(token)
	if err != nil {
		ReturnError(c, "token error, "+err.Error())
		return
	}
	c.JSON(200, gin.H{
		"result": "ok",
	})
}
