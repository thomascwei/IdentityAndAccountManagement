package main

import (
	"IAM/api"
	_ "github.com/go-sql-driver/mysql"
	"log"
	//"net/http"
	"github.com/gin-gonic/gin"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	//以Gin框架起一個post接收數據, 收到後塞進該點位的專屬channel
	r := gin.Default()

	r.POST("/IAM/V1/Login", api.Login)

	r.Run(":9567")
}
