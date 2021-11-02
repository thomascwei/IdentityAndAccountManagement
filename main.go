package main

import (
	"IAM/pkg/cache"
	"IAM/pkg/model"
	"IAM/pkg/token"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	//id, err := model.CreateAccounts("Admin", "123456", "", 255)
	//log.Println(id, err)
	//id, err = model.CreateAccounts("Manager", "55688", "", 200)
	//log.Println(id, err)
	//id, err = model.CreateAccounts("Member", "00757", "abc@123.com", 100)
	//log.Println(id, err)

	rows, _ := model.QueryAllAccounts()
	log.Println("取得全部帳號", rows)

	// 更新account
	newnew := make(map[string]interface{})
	newnew["id"] = 3
	newnew["password"] = "0000000000"
	newnew["email"] = "444@333.com"
	err := model.UpdateAccount(newnew)
	if err != nil {
		log.Fatal(err)
	}
	tk, _ := token.GenerateToken("")
	log.Println("token:", tk)
	//gc := cache.BuildCacheObject()
	//fmt.Printf("%T", gc)
	TokenID := tk
	Tokenvv := cache.TokenValue{Id: 1, Auth: 200}
	// 設置token帶有效期限
	cache.SetWithExpire(TokenID, Tokenvv, 300)

	tk2, _ := token.GenerateToken("")
	log.Println("token2:", tk2)
	cache.SetWithExpire(tk2, cache.TokenValue{Id: 2, Auth: 400}, 300)

	value := cache.GetAllCache()
	log.Println(value)
	for k, v := range value {
		userid := v.(cache.TokenValue).Id
		if userid == 2 {
			cache.CacheRemove(k.(string))
		}
	}
	value2 := cache.GetAllCache()
	log.Println(value2)
}
