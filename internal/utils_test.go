package internal

import (
	"IAM/pkg/cache"
	"IAM/pkg/model"
	"IAM/pkg/password"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"testing"
)

var db, err = sql.Open("mysql", "thomas:123456@/iam?charset=utf8&parseTime=true")
var _, _ = db.Exec("CREATE TABLE IF NOT EXISTS `accounts`\n(\n    `id`       int          NOT NULL AUTO_INCREMENT,\n    `username` VARCHAR(30)  NOT NULL,\n    `password` VARCHAR(100) NOT NULL,\n    `email`    VARCHAR(50)  not NULL,\n    `auth`     INT          not NULL,\n    UNIQUE (`username`),\n    PRIMARY KEY (`id`)\n);")
var _, _ = db.Exec("INSERT accounts\nSET username='Admin',\n    password='$2a$10$dN7Da733DxGG4CLLfRQ.5OV8UakM8H1yo5o1aWj9uOGPSBU7ZmmY6',\n    Email='admin@admin.com',\n    Auth=255;")
var _, _ = db.Exec("INSERT accounts\nSET username='Manager',\n    password='$2a$10$dN7Da733DxGG4CLLfRQ.5OV8UakM8H1yo5o1aWj9uOGPSBU7ZmmY6',\n    Email='manager@admin.com',\n    Auth=200;")

func TestSignUp(t *testing.T) {
	username := "test!!!"
	password := "test1231qaz!QAZ"
	email := "test@example.com"
	auth := 123

	got, err := SignUp(username, password, email, auth)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	rows, err := db.Query("SELECT id FROM accounts where username='test!!!'")

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

	model.DeleteAccount(int32(got))
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

func TestTokenverify(t *testing.T) {
	username := "Admin"
	password := "123456"
	token, err := Login(username, password)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	_, got, err := Tokenverify(token)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	want := 255
	if got != want {
		t.Errorf("token verify error, auth got %v, want %v", got, want)
	}
}

func TestGetAllAccount(t *testing.T) {
	allAccounts, err := GetAllAccount()
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	want := "admin@admin.com"
	got := ""
	for _, account := range allAccounts {
		if account.Username == "Admin" {
			got = account.Email
		}
	}
	if got != want {
		t.Errorf("email not match , got %v, want %v", got, want)
	}
}

func TestUpdateSingelAccount(t *testing.T) {
	fields := make(map[string]interface{})
	fields["Id"] = 1
	fields["Auth"] = 100
	fields["Username"] = "Admin"
	fields["Email"] = "admin@admin.com"
	want := 100
	var got int
	err := UpdateSingelAccount(fields)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	rows, err := db.Query("SELECT auth FROM accounts where id=1")
	for rows.Next() {
		err = rows.Scan(&got)
	}
	if err != nil {
		t.Errorf(err.Error())
	}
	if got != want {
		t.Errorf("update fail, auth not match , got %v, want %v", got, want)
	}
	// 通過測試後改回原數據
	fields["Auth"] = 255
	err = UpdateSingelAccount(fields)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
}

func TestRenewPassword(t *testing.T) {
	want := "1qaz@WSX3edcZ012"
	err := ChangeSelfPassword(2, want)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	var got string
	rows, err := db.Query("SELECT password FROM accounts where id=2")
	for rows.Next() {
		err = rows.Scan(&got)
	}
	if err != nil {
		t.Errorf(err.Error())
	}
	err = password.CheckPassword(want, got)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestInitPassword(t *testing.T) {
	want := "123456"
	err := InitPassword(2, want)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	var got string
	rows, err := db.Query("SELECT password FROM accounts where id=2")
	for rows.Next() {
		err = rows.Scan(&got)
	}
	if err != nil {
		t.Errorf(err.Error())
	}
	err = password.CheckPassword(want, got)
	if err != nil {
		t.Errorf(err.Error())
	}

}
