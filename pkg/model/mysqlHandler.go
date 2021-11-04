package model

import (
	"IAM/pkg/cache"
	"IAM/pkg/password"
	pd "IAM/pkg/password"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type AccountFields struct {
	ID       int    `json:"ID"`
	Username string `json:"username"`
	Email    string `json:"Email"`
	Auth     int    `json:"Auth"`
}

var db, err = sql.Open("mysql", "thomas:123456@/iam?charset=utf8")

func CreateAccounts(username, password, email string, auth int) (int64, error) {
	password = pd.Encryption(password)
	stmt, err := db.Prepare("INSERT accounts SET username=?,password=?,Email=?,Auth=?")
	if err != nil {
		return 0, err
	}
	res, err := stmt.Exec(username, password, email, auth)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, err
}

func QueryAllAccounts() ([]AccountFields, error) {
	rows, err := db.Query("SELECT * FROM accounts")
	if err != nil {
		return nil, err
	}
	var result = make([]AccountFields, 0)
	for rows.Next() {
		var account AccountFields
		var pwd string
		err := rows.Scan(&account.ID, &account.Username, &pwd, &account.Email, &account.Auth)
		if err != nil {
			log.Fatal(err)
		}
		result = append(result, account)
	}
	return result, err
}

// 驗證帳號密碼, 通過返回true, 失敗返回false
func VerifyPassword(username, password string) (bool, cache.TokenValue, error) {
	rows, err := db.Query("SELECT ID, Auth, password FROM accounts where username='" + username + "'")
	if err != nil {
		return false, cache.TokenValue{}, err
	}
	var tokenvalue cache.TokenValue
	var pwd string
	//var account AccountFields
	for rows.Next() {
		err = rows.Scan(&tokenvalue.Id, &tokenvalue.Auth, &pwd)
	}
	if err != nil {
		return false, cache.TokenValue{}, err
	}
	if pd.Encryption(password) == pwd {
		return true, tokenvalue, nil
	}
	return false, cache.TokenValue{}, errors.New("account password invalid")
}

// 可接受部分欄位變更
func UpdateAccount(params map[string]interface{}) error {
	//編成prepare string
	SetSection := " SET "
	WhereClause := " where ID =?"
	args := make([]interface{}, 0)
	iidd, ok := params["Id"]
	if !ok {
		return fmt.Errorf("Id column not found")
	}

	for k, v := range params {
		if k == "Password" {
			SetSection += k + "= ?,"
			args = append(args, password.Encryption(v.(string)))
			//fmt.Println("args0: ", args)

		} else {
			if k != "Id" {
				SetSection += k + "= ?,"
				args = append(args, v)
			}
		}
	}
	args = append(args, iidd)
	SetSection = SetSection[:len(SetSection)-1]
	RawString := "update accounts "
	RawString = RawString + SetSection + WhereClause
	//fmt.Println("update syntax: ", RawString)
	stmt, err := db.Prepare(RawString)
	if err != nil {
		return err
	}
	//fmt.Println("args1: ", args)
	_, err = stmt.Exec(args...)
	if err != nil {
		return err
	}
	return nil
}

func DeleteAccount(id int64) error {
	//刪除資料
	stmt, err := db.Prepare("delete from accounts where ID=?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	return nil
}
