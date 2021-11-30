package password

import (
	"fmt"
	"github.com/trustelem/zxcvbn"
	"golang.org/x/crypto/bcrypt"
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
func Encryption(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}

// CheckPassword checks if the provided password is correct or not
func CheckPassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
