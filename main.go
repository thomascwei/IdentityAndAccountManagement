package main

import (
	"IAM/internal"
	"io"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
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

	r.POST("/IAM/V1/Login", internal.LoginRoute)
	r.GET("/IAM/V1/Logout", internal.LogoutRoute)

	// 加一層middleware驗證token
	authRoutes := r.Group("/").Use(internal.AuthMiddleware())
	authRoutes.GET("/IAM/V1/all_accounts", internal.GetAllAccountRoute)
	authRoutes.POST("/IAM/V1/create_account", internal.SignUpRoute)
	authRoutes.POST("/IAM/V1/account_update", internal.AccountUpdateRoute)
	authRoutes.POST("/IAM/V1/change_password", internal.ChangeSelfPasswordRoute)
	authRoutes.POST("/IAM/V1/init_password", internal.InitPasswordRoute)
	authRoutes.GET("/IAM/V1/token_verify", internal.TokenverifyRoute)

	r.Run(":9567")
}
