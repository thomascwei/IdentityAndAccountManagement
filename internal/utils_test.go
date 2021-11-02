package internal

import (
	"IAM/pkg/cache"
	"IAM/pkg/model"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"testing"
)

var db, err = sql.Open("mysql", "root:123456@/iam?charset=utf8")

func TestSignUp(t *testing.T) {
	username := "test"
	password := "test1231qaz!QAZ"
	email := "test@example.com"
	auth := 123

	got, err := SignUp(username, password, email, auth)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	rows, err := db.Query("SELECT id FROM accounts where username='test'")

	var want int64
	for rows.Next() {
		err = rows.Scan(&want)
	}
	if err != nil {
		t.Errorf(err.Error())
	}
	if got != want {
		t.Errorf("ID mismatch, create account fail")
	}

	model.DeleteAccount(got)
}

func TestLogin(t *testing.T) {
	username := "Admin"
	password := "123456"
	token, err := Login(username, password)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	value, err := cache.CacheGet(token)
	value2 := value.(cache.TokenValue)
	if value2.Id != 1 {
		t.Errorf("not get correct token from cache")
		return
	}
}

func TestLogout(t *testing.T) {
	username := "Admin"
	password := "123456"
	token, err := Login(username, password)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	err = Logout(token)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	_, err = cache.CacheGet(token)
	if err == nil {
		t.Errorf("remove token fail")
	}

}
