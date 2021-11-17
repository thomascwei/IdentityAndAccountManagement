package main

import (
	"IAM/api"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"log"
	"os"
)

var (
	Trace *log.Logger
	Info  *log.Logger
	Error *log.Logger
)

func init() {
	file, err := os.OpenFile("./log/main.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("can not open log file")
	}
	Trace = log.New(os.Stdout, "TRACE: ", log.Ldate|log.Ltime|log.Lshortfile)
	Info = log.New(io.MultiWriter(file, os.Stdout), "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(io.MultiWriter(file, os.Stdout), "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}
func main() {
	Trace.Println("Starting...")
	Info.Println("Starting... INFO")
	Error.Println("Starting... ERROR")
	//以Gin框架起一個post接收數據, 收到後塞進該點位的專屬channel
	r := gin.Default()

	r.POST("/IAM/V1/Login", api.Login)
	r.GET("/IAM/V1/Logout", api.Logout)
	r.POST("/IAM/V1/create_account", api.SignUp)
	r.GET("/IAM/V1/all_accounts", api.GetAllAccount)
	r.POST("/IAM/V1/account_update", api.AccountUpdate)
	r.POST("/IAM/V1/change_password", api.ChangeSelfPassword)
	r.POST("/IAM/V1/init_password", api.InitPassword)

	r.Run(":9567")
}
