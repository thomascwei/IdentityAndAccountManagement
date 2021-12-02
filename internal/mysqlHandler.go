package internal

import (
	"IAM/pkg/cache"
	dbb "IAM/pkg/db"
	pd "IAM/pkg/password"
	"context"
	_ "github.com/go-sql-driver/mysql"
)

//var db, err = sql.Open("mysql", "thomas:123456@/iam?charset=utf8")

// 驗證帳號密碼, 通過返回true, 失敗返回false
func VerifyPassword(username, password string) (bool, cache.TokenValue, error) {
	rows, err := db.Query("SELECT ID, Auth, password FROM accounts where username='" + username + "'")
	if err != nil {
		return false, cache.TokenValue{}, err
	}
	var tokenvalue cache.TokenValue
	var DBPwd string
	for rows.Next() {
		err = rows.Scan(&tokenvalue.Id, &tokenvalue.Auth, &DBPwd)
	}
	if err != nil {
		return false, cache.TokenValue{}, err
	}
	err = pd.CheckPassword(password, DBPwd)
	if err != nil {
		return false, cache.TokenValue{}, err
	}
	return true, tokenvalue, nil
}

/*
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
*/

var ctx = context.Background()
var queries = dbb.New(db)

// 列出全部帳戶
func QueryAllAccounts() ([]dbb.ListAccountsRow, error) {
	accounts, err := queries.ListAccounts(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]dbb.ListAccountsRow, 0)
	for _, account := range accounts {
		var temp = dbb.ListAccountsRow{
			ID:       account.ID,
			Username: account.Username,
			Email:    account.Email,
			Auth:     account.Auth,
		}
		result = append(result, temp)
	}

	return result, nil
}

// 建立新帳戶
func CreateAccount(username, password, email string, auth int32) (int64, error) {
	password, err = pd.Encryption(password)
	if err != nil {
		return 0, err
	}
	result, err := queries.CreateAccount(ctx, dbb.CreateAccountParams{
		Username: username,
		Password: password,
		Email:    email,
		Auth:     auth,
	})
	if err != nil {
		return 0, err
	}
	insertedAccountID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return insertedAccountID, nil
}

//刪除帳戶
func DeleteAccount(id int32) error {
	err := queries.DeleteAccount(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

//只改密碼
func ChangePassword(id int32, password string) error {
	password, err = pd.Encryption(password)
	if err != nil {
		return err
	}
	input := dbb.UpdatePasswordParams{
		Password: password,
		ID:       id,
	}
	err := queries.UpdatePassword(ctx, input)
	if err != nil {
		return err
	}
	return nil
}

// 欄位給足才可修改
func UpdateAccountAllField(Username string, Email string, Auth int32, ID int32) error {
	input := dbb.UpdateAccountParams{
		Username: Username,
		Email:    Email,
		Auth:     Auth,
		ID:       ID,
	}
	err := queries.UpdateAccount(ctx, input)
	if err != nil {
		return err
	}
	return nil
}
