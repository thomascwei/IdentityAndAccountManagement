package cache

import (
	"fmt"
	"github.com/bluele/gcache"
	"strconv"
	"time"
)

type TokenValue struct {
	Id   int
	Auth int
}

var gc = gcache.New(20).Build()

// 將token保存並設置有效期
func SetWithExpire(TokenID string, Tokenvv interface{}, seconds int) {
	gc.SetWithExpire(TokenID, Tokenvv, time.Second*time.Duration(seconds))
}

// 取得指定token的數據
func CacheGet(key string) (value interface{}, err error) {
	value, err = gc.Get(key)
	return
}

// 將指定的key從cache移除
func CacheRemove(key string) bool {
	ok := gc.Remove(key)
	return ok
}

// 取得全部cache的數據
func GetAllCache() map[interface{}]interface{} {
	return gc.GetALL(true)
}

func main() {
	TokenID := "ABCDERTJFKLDD:D"
	Tokenvv := TokenValue{Id: 1, Auth: 200}

	gc := gcache.New(1).Build()
	for i := 0; i < 10000; i++ {
		gc.SetWithExpire("AAA"+strconv.Itoa(i), TokenValue{Id: i, Auth: i}, time.Second*10)
	}

	// 設置token帶有效期限
	gc.SetWithExpire(TokenID, Tokenvv, time.Second*10)
	value, _ := gc.Get(TokenID)
	aa := value.(TokenValue)
	fmt.Println(aa.Id)
	fmt.Println(aa.Auth)

	// 完全一樣的k,v 重新寫入(每當Token驗證成功就刷新該token的有效期)
	gc.SetWithExpire(TokenID, value, time.Second*10)

	// 可以正常取得
	time.Sleep(time.Second * 5)
	value1, _ := gc.Get(TokenID)
	bb := value1.(TokenValue)
	fmt.Println(bb.Id)
	fmt.Println(bb.Auth)
	value, err := gc.Get(TokenID)
	if err != nil {
		panic(err)
	}
	fmt.Println("bbGet:", value)

	valuevv, err := gc.Get("AAA9999")
	fmt.Println("valuevv", valuevv)
	gc.Remove("AAA9998")
	valuevv, err = gc.Get("AAA9998")
	if err != nil {
		panic(err)
	}
	fmt.Println("valuevv", valuevv)
	// 過期後無法取得
	time.Sleep(time.Second * 6)
	value2, err := gc.Get(TokenID)
	if err != nil {
		panic(err)
	}
	fmt.Println(TokenID, value2)

}
