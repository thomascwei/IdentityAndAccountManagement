package internal

import (
	"IAM/pkg/cache"
	"IAM/pkg/model"
	pd "IAM/pkg/password"
	"IAM/pkg/token"
	"errors"
)

// 建立新帳號
func SignUp(username, password, email string, auth int) (int64, error) {
	ok := pd.CheckPasswordStrength(password)
	if !ok {
		return 0, errors.New("not enough password strength")
	}
	id, err := model.CreateAccounts(username, password, email, auth)

	return id, err
}

// 比對SQL的帳密, 成功返回token, 失敗返回失敗原因
func Login(username, password string) (string, error) {

	// 驗證帳密
	ok, tokenvalue, err := model.VerifyPassword(username, password)
	if !ok {
		return "", err
	}
	// 砍掉原本的token
	allToken := cache.GetAllCache()
	for token, vv := range allToken {
		if vv.(cache.TokenValue).Id == tokenvalue.Id {
			cache.CacheRemove(token.(string))
		}
	}
	// 建立新token
	tokenid, err := token.GenerateToken(username)
	if err != nil {
		return "", err
	}
	// 寫進cache
	cache.SetWithExpire(tokenid, tokenvalue, 300)
	// return token
	return tokenid, nil
}

// login驗證密碼通過後刪除原有token
func Logout(token string) error {
	ok := cache.CacheRemove(token)
	if ok {
		return nil
	} else {
		return errors.New("remove token fail")
	}

}
