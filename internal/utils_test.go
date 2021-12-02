package internal

import (
	"IAM/pkg/cache"
	"IAM/pkg/password"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
	"testing"
)

var _, _ = db.Exec("CREATE TABLE IF NOT EXISTS `accounts`\n(\n    `id`       int          NOT NULL AUTO_INCREMENT,\n    `username` VARCHAR(30)  NOT NULL,\n    `password` VARCHAR(100) NOT NULL,\n    `email`    VARCHAR(50)  not NULL,\n    `auth`     INT          not NULL,\n    UNIQUE (`username`),\n    PRIMARY KEY (`id`)\n);")
var _, _ = db.Exec("INSERT accounts\nSET username='Admin',\n    password='$2a$10$dN7Da733DxGG4CLLfRQ.5OV8UakM8H1yo5o1aWj9uOGPSBU7ZmmY6',\n    Email='admin@admin.com',\n    Auth=255;")
var _, _ = db.Exec("INSERT accounts\nSET username='Manager',\n    password='$2a$10$dN7Da733DxGG4CLLfRQ.5OV8UakM8H1yo5o1aWj9uOGPSBU7ZmmY6',\n    Email='manager@admin.com',\n    Auth=200;")

func TestSignUp(t *testing.T) {
	username := "test!!!"
	password := "test1231qaz!QAZ"
	email := "test@example.com"
	auth := 123

	got, err := SignUp(username, password, email, auth)
	require.NoError(t, err)

	rows, err := db.Query("SELECT id FROM accounts where username='test!!!'")
	var want int64
	for rows.Next() {
		err = rows.Scan(&want)
	}
	require.NoError(t, err)
	require.Equal(t, want, got)

	DeleteAccount(int32(got))
}

func TestLogin(t *testing.T) {
	username := "Admin"
	password := "123456"
	token, err := Login(username, password)
	require.NoError(t, err)

	value, err := cache.CacheGet(token)
	value2 := value.(cache.TokenValue)
	require.Equalf(t, 1, value2.Id, "not correct token from cache")
}

func TestLogout(t *testing.T) {
	username := "Admin"
	passWord := "123456"
	token, err := Login(username, passWord)
	require.NoError(t, err)

	err = Logout(token)
	require.NoError(t, err)

	_, err = cache.CacheGet(token)
	require.Errorf(t, err, "remove token fail")

}

func TestTokenverify(t *testing.T) {
	username := "Admin"
	passWord := "123456"
	token, err := Login(username, passWord)
	require.NoError(t, err)

	_, got, err := Tokenverify(token)
	require.NoError(t, err)

	want := 255
	require.Equalf(t, want, got, "token verify error, auth got %v, want %v", got, want)
}

func TestGetAllAccount(t *testing.T) {
	allAccounts, err := GetAllAccount()
	require.NoError(t, err)

	want := "admin@admin.com"
	got := ""
	for _, account := range allAccounts {
		if account.Username == "Admin" {
			got = account.Email
		}
	}
	require.Equalf(t, want, got, "email not match , got %v, want %v", got, want)
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
	require.NoError(t, err)

	rows, err := db.Query("SELECT auth FROM accounts where id=1")
	for rows.Next() {
		err = rows.Scan(&got)
	}
	require.NoError(t, err)
	require.Equalf(t, want, got, "update fail, auth not match , got %v, want %v", got, want)

	// 通過測試後改回原數據
	fields["Auth"] = 255
	err = UpdateSingelAccount(fields)
	require.NoError(t, err)
}

func TestRenewPassword(t *testing.T) {
	want := "1qaz@WSX3edcZ012"
	err := ChangeSelfPassword(2, want)
	require.NoError(t, err)

	var got string
	rows, err := db.Query("SELECT password FROM accounts where id=2")
	for rows.Next() {
		err = rows.Scan(&got)
	}
	require.NoError(t, err)

	err = password.CheckPassword(want, got)
	require.NoError(t, err)

}

func TestInitPassword(t *testing.T) {
	want := "123456"
	err := InitPassword(2, want)
	require.NoError(t, err)

	var got string
	rows, err := db.Query("SELECT password FROM accounts where id=2")
	for rows.Next() {
		err = rows.Scan(&got)
	}
	require.NoError(t, err)

	err = password.CheckPassword(want, got)
	require.NoError(t, err)

}
