package main

import (
	"IAM/api"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	//以Gin框架起一個post接收數據, 收到後塞進該點位的專屬channel
	r := gin.Default()

	r.POST("/IAM/V1/Login", api.Login)
	r.GET("/IAM/V1/Logout", api.Logout)
	r.POST("/IAM/V1/create_account", api.SignUp)
	r.GET("/IAM/V1/all_accounts", api.GetAllAccount)
	r.POST("/IAM/V1/account_update", api.AccountUpdate)
	r.POST("/IAM/V1/change_password", api.ChangePassword)
	r.POST("/IAM/V1/init_password", api.InitPassword)

	r.Run(":9567")
}
