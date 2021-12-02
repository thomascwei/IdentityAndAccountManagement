package internal

import (
	"fmt"
	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
	yamlFile, err := ioutil.ReadFile("config/auth.yaml")
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

func ReturnError(c *gin.Context, statusCode int, info string) {
	Info.Println(info)
	c.JSON(statusCode, gin.H{
		"result": "fail",
		"error":  info,
	})
}

// 建立帳戶
func SignUpRoute(c *gin.Context) {
	authInterface, ok := c.Get("auth")
	if !ok {
		ReturnError(c, 500, "cannot get auth from middleware")
		return
	}
	auth := authInterface.(int)
	// 確認有使用此API的權限
	if auth < CConfig.API.SignUp {
		ReturnError(c, http.StatusUnauthorized, "cannot get auth from middleware")
		return
	}
	// 讀request body
	input := User{}
	err := c.BindJSON(&input)
	if err != nil {
		ReturnError(c, http.StatusBadRequest, err.Error())
		return
	}
	// 確認新帳號權限小於自身帳號
	if auth <= input.Auth {
		ReturnError(c, http.StatusUnprocessableEntity, "you cannot create authority higher than or equal to yourself")
		return
	}
	// 建立帳號
	id, err := SignUp(input.Username, input.Password, input.Email, input.Auth)
	if err != nil {
		ReturnError(c, http.StatusInternalServerError, err.Error())
		return
	}
	// 成功建立帳號後返回id
	c.JSON(200, gin.H{
		"result": "ok",
		"id":     id,
	})
	return
}

// 登入
func LoginRoute(c *gin.Context) {
	m := login{}
	err := c.Bind(&m)
	if err != nil {
		ReturnError(c, http.StatusBadRequest, err.Error())
		return
	}
	// 檢查request body轉type後是否少欄位
	if m.Username == "" || m.Password == "" {
		ReturnError(c, http.StatusBadRequest, "request body format error")
		return
	}
	token, err := Login(m.Username, m.Password)
	if err != nil {
		ReturnError(c, http.StatusInternalServerError, err.Error())
		return
	}
	// 登入成功後返回token
	c.JSON(200, gin.H{
		"result": "ok",
		"token":  token,
	})
	Info.Println(m.Username, "login successfully")
	return
}

// 登出
func LogoutRoute(c *gin.Context) {
	// 找出token
	token := ""
	for k, v := range c.Request.Header {
		if k == "Authorization" {
			token = v[0]
		}
	}
	// 刪除token, 不管是否存在
	_ = Logout(token)
	c.JSON(200, gin.H{
		"result": "ok",
	})
}

// 取得帳號戶清單及內容
func GetAllAccountRoute(c *gin.Context) {
	authInterface, ok := c.Get("auth")
	if !ok {
		ReturnError(c, http.StatusInternalServerError, "cannot get auth from middleware")
		return
	}
	auth := authInterface.(int)
	// 確認有使用此API的權限
	if auth < CConfig.API.GetAllAccount {
		ReturnError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	result, err := GetAllAccount()
	if err != nil {
		ReturnError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(200, gin.H{
		"result":   "ok",
		"accounts": result,
	})
}

// 更新帳戶內容
func AccountUpdateRoute(c *gin.Context) {
	// 從middleware context取得權限
	authInterface, ok := c.Get("auth")
	if !ok {
		ReturnError(c, http.StatusInternalServerError, "cannot get auth from middleware")
		return
	}
	auth := authInterface.(int)

	// 確認有使用此API的權限
	if auth < CConfig.API.Update {
		ReturnError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	// 讀request body
	input := User{}
	err := c.BindJSON(&input)
	if err != nil {
		ReturnError(c, http.StatusBadRequest, err.Error())
		return
	}
	// 新權限必須小於此token
	if input.Auth >= auth {
		ReturnError(c, http.StatusUnprocessableEntity, "you cannot create authority higher than or equal to yourself")
		return
	}
	// 轉成map
	UserMap := structs.Map(input)
	// 刪除password, 密碼要在其他API單獨改
	delete(UserMap, "Password")
	fmt.Println("UserMap:", UserMap)
	err = UpdateSingelAccount(UserMap)
	if err != nil {
		ReturnError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(200, gin.H{
		"result": "ok",
	})
}

// 改自己的密碼
func ChangeSelfPasswordRoute(c *gin.Context) {
	// 從middleware context取得權限
	idInterface, ok := c.Get("id")
	if !ok {
		ReturnError(c, http.StatusInternalServerError, "cannot get auth from middleware")
		return
	}
	id := idInterface.(int)

	// 讀request body
	input := changePassword{}
	err := c.BindJSON(&input)
	if err != nil {
		ReturnError(c, http.StatusBadRequest, err.Error())
		return
	}
	err = ChangeSelfPassword(id, input.NewPassword)
	if err != nil {
		ReturnError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(200, gin.H{
		"result": "ok",
	})

}

// 初始化密碼(忘記密碼時使用)
func InitPasswordRoute(c *gin.Context) {
	// 從middleware context取得權限
	authInterface, ok := c.Get("auth")
	if !ok {
		ReturnError(c, http.StatusInternalServerError, "cannot get auth from middleware")
		return
	}
	auth := authInterface.(int)
	// 確認有使用此API的權限
	if auth < CConfig.API.InitPassword {
		ReturnError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	// 讀request body
	input := changePassword{}
	err := c.BindJSON(&input)
	if err != nil {
		ReturnError(c, http.StatusBadRequest, err.Error())
		return
	}
	err = InitPassword(input.Id, input.NewPassword)
	if err != nil {
		ReturnError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(200, gin.H{
		"result": "ok",
	})

}

// 確認token是否有效
func TokenverifyRoute(c *gin.Context) {
	// middleware通過後回傳ok
	c.JSON(200, gin.H{
		"result": "ok",
	})
}
