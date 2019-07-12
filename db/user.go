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
	if len(pRows) > 0 && string(pRows[0]["user_pwd"].([]byte)) == encpwd {
		return true
	}
	return false;
}

// 刷新用户登录的 Token
func UpdateToken(username string, token string) bool {
	stmt, err := mydb.DBConn().Prepare("replace into tbl_user_token (`user_name`, `user_token`) values (?, ?)")
	if err != nil {
		fmt.Println("Failed to login, err:"+ err.Error())
		return false
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, token)
	if err != nil {
		fmt.Println("Failed to login, err:"+ err.Error())
		return false
	}
	return true
}

type User struct {
	Username string
	Email string
	Phone string
	SignupAt string
	LastActiveAt string
	Status int
}

// 用户信息查询
func GetUserInfo(username string) (User, error) {
	user := User{}

	stmt, err := mydb.DBConn().Prepare("select user_name, email, phone, signup_at, last_active, status from tbl_user where user_name = ? limit 1")
	if err != nil {
		fmt.Println("Failed to select, err:"+ err.Error())
		return user, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(username).Scan(&user.Username, &user.SignupAt)
	if err != nil {
		return user, err
	}
	return user, nil
}