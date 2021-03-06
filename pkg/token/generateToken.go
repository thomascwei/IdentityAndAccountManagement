package token

import (
	"crypto/md5"
	"encoding/hex"
	"golang.org/x/crypto/bcrypt"
)

// GenerateToken returns a unique token based on the provided string
func GenerateToken(ss string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(ss), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	//fmt.Println("Hash to store:", string(hash))

	hasher := md5.New()
	hasher.Write(hash)
	return hex.EncodeToString(hasher.Sum(nil)), nil
}
