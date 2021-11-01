package model

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type Account struct {
	id       int    `json:"id"`
	username string `json:"username"`
	email    string `json:"email"`
	auth     int    `json:"auth"`
}

func checkErr(err error) error {
	return err
}

var db, err = sql.Open("mysql", "root:123456@/iam?charset=utf8")

func CreateAccounts(username, password, email string, auth int) (int64, error) {
	stmt, err := db.Prepare("INSERT accounts SET username=?,password=?,email=?,auth=?")
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

func QueryAllAccounts() ([]Account, error) {
	rows, err := db.Query("SELECT * FROM accounts")
	if err != nil {
		return nil, err
	}
	var result = make([]Account, 0)
	for rows.Next() {
		var account Account
		var pwd string
		err := rows.Scan(&account.id, &account.username, &pwd, &account.email, &account.auth)
		if err != nil {
			log.Fatal(err)
		}
		result = append(result, account)
	}
	return result, err
}

// 驗證帳號密碼, 通過返回true, 失敗返回false
func VerifyPassword(username, password string) (bool, error) {
	rows, err := db.Query("SELECT password FROM accounts where username='" + username + "'")
	if err != nil {
		return false, err
	}
	var pwd string
	//var account Account
	for rows.Next() {
		err = rows.Scan(&pwd)
	}
	if err != nil {
		return false, err
	}
	if password == pwd {
		return true, nil
	}
	return false, nil
}

// 可接受部分欄位變更
func UpdateAccount(params map[string]interface{}) error {
	//編成prepare string
	SetSection := " SET "
	WhereClause := " where id =?"
	args := make([]interface{}, 0)
	iidd, ok := params["id"]
	if !ok {
		return fmt.Errorf("id column not found")
	}

	for k, v := range params {
		if k != "id" {
			SetSection += k + "= ?,"
			args = append(args, v)
		}
	}
	args = append(args, iidd)
	SetSection = SetSection[:len(SetSection)-1]
	RawString := "update accounts "
	RawString = RawString + SetSection + WhereClause
	//fmt.Println("update syntax", RawString)
	stmt, err := db.Prepare(RawString)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(args...)
	if err != nil {
		return err
	}
	return nil
}
