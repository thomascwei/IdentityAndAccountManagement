package internal

import (
	"IAM/pkg/cache"
	"IAM/pkg/model"
	pd "IAM/pkg/password"
	"IAM/pkg/token"
	"errors"
)

func SignUp(username, password, email string, auth int) (int64, error) {
	// 檢查密碼強度
	ok := pd.CheckPasswordStrength(password)
	if !ok {
		return 0, errors.New("not enough password strength")
	}
	// 新增到DB
	id, err := model.CreateAccounts(username, password, email, auth)

	return id, err
}

// 比對SQL的帳密, 成功返回token,nil ; 失敗返回""與error
func Login(username, password string) (string, error) {
	// 驗證帳密
	ok, IdAuth, err := model.VerifyPassword(username, password)
	if !ok {
		return "", err
	}
	// 砍掉原本的token
	allToken := cache.GetAllCache()
	for token, vv := range allToken {
		if vv.(cache.TokenValue).Id == IdAuth.Id {
			cache.CacheRemove(token.(string))
		}
	}
	// 建立新token
	tokenid, err := token.GenerateToken(username)
	if err != nil {
		return "", err
	}
	// 寫進cache
	cache.SetWithExpire(tokenid, IdAuth, 300)
	// return token
	return tokenid, nil
}

func Logout(token string) error {
	// 刪除原有token
	ok := cache.CacheRemove(token)
	if ok {
		return nil
	} else {
		return errors.New("token not exist")
	}
}

// 確認token是否有效, 同時更新時效
func Tokenverify(token string) (int, error) {
	// 從cachek,v
	value, err := cache.CacheGet(token)
	// cache取不到返回失敗
	if err != nil {
		return -1, err
	}
	// 更新token時效
	cache.SetWithExpire(token, value, 300)

	return value.(cache.TokenValue).Auth, nil
}

// 取得所有帳號信息, 除了密碼
func GetAllAccount() ([]model.AccountFields, error) {
	allAccounts, err := model.QueryAllAccounts()
	if err != nil {
		return nil, err
	}
	return allAccounts, nil
}

// 更新帳號內容
func UpdateSingelAccount(params map[string]interface{}) error {
	// 改密碼要用RenewPassword, map裡如果有password先移除
	_, ok := params["password"]
	if ok {
		delete(params, "password")
	}
	err := model.UpdateAccount(params)
	if err != nil {
		return err
	}
	return nil
}

// 更新密碼
func RenewPassword(id int, password string) error {
	// 檢查密碼強度
	ok := pd.CheckPasswordStrength(password)
	if !ok {
		return errors.New("password is not strong enough")
	}
	params := make(map[string]interface{})
	params["id"] = id
	params["password"] = password
	err := model.UpdateAccount(params)
	if err != nil {
		return err
	}
	return nil

}

// admin初始化密碼專用,不檢查密碼強度
func InitPassword(id int, password string) error {
	params := make(map[string]interface{})
	params["id"] = id
	params["password"] = password
	err := model.UpdateAccount(params)
	if err != nil {
		return err
	}
	return nil

}
