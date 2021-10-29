package main

import (
	"IAM/pkg/cache"
	"IAM/pkg/token"
	"fmt"
)

func main() {
	tk := token.GenerateToken("")
	fmt.Println("token:", tk)
	//gc := cache.BuildCacheObject()
	//fmt.Printf("%T", gc)
	TokenID := tk
	Tokenvv := cache.TokenValue{Id: 1, Auth: 200}
	// 設置token帶有效期限
	cache.SetWithExpire(TokenID, Tokenvv, 300)

	tk2 := token.GenerateToken("")
	fmt.Println("token2:", tk2)
	cache.SetWithExpire(tk2, cache.TokenValue{Id: 2, Auth: 400}, 300)

	value := cache.GetAllCache()
	fmt.Println(value)
	for k, v := range value {
		userid := v.(cache.TokenValue).Id
		if userid == 2 {
			cache.CacheRemove(k.(string))
		}
	}
	value2 := cache.GetAllCache()
	fmt.Println(value2)
}
