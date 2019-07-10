package db

import (
	mydb "filestore-server/db/mysql"
	"fmt"
)

// 通过用户名及密码完成 user 表的操作
func UserSignUp(username string, password string) bool {
	stmt, err := mydb.DBConn().Prepare("insert ignore into tbl_user (`user_name`, `user_pwd`) values (?, ?)")
	if err != nil {
		fmt.Println("Failed to insert, err:"+ err.Error());
		return false
	}
	defer stmt.Close()

	ret, err := stmt.Exec(username, password)
	if err != nil {
		fmt.Println("Faile to insert, err:"+ err.Error())
		return false
	}

	if rowsAffected, err := ret.RowsAffected(); nil == err && rowsAffected > 0 {
		return true
	}
	return false
}

// 判断密码是否一致
func UserSignIn(username string, encpwd string) bool {
	stmt, err := mydb.DBConn().Prepare("select * from tbl_user where user_name = ? limit 1")
	if err != nil {
		fmt.Println("Failed to select, err:"+ err.Error())
		return false
	}
	defer stmt.Close()

	rows, err := stmt.Query(username)
	if err != nil {
		fmt.Println("Failed to select, err:"+ err.Error())
		return false
	} else if rows == nil {
		fmt.Println("username not found:"+ username)
		return false
	}

	pRows := mydb.ParseRows(rows)
	if len(pRows) > 0 && string(pRows[0]["user_pwd"].([]byte) == encpwd) {
		return true
	}
	return false;
}