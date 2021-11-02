package password

import (
	"crypto/sha256"
	"fmt"
	"github.com/trustelem/zxcvbn"
)

// 檢查密碼強度
func CheckPasswordStrength(pwd string) bool {
	res := zxcvbn.PasswordStrength(pwd, nil)
	// password is safe if the zxcvbn score is >= 3
	if res.Score >= 3 {
		return true
	}
	return false
}

// 密碼加密
func Encryption(pwd string) string {
	sum := sha256.Sum256([]byte(pwd))
	return fmt.Sprintf("%x", sum)
}
