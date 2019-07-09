package handler

import (
	"net/http"
	"io/ioutil"
	"filestore-server/util"
	dblayer "filestore-server/db"
)

const (
	pwd_salt = "!@#$5"
)

// 处理用户注册请求
func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		data, err := ioutil.ReadFile("./static/view/signup.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(data)
		return
	} else if r.Method == "POST" {
		r.ParseForm()

		username := r.Form.Get("username")
		password := r.Form.Get("password")

		if len(username) < 3 || len(password) < 5 {
			w.Write([]byte("Invalid paramater"))
			return
		}

		enc_password := util.Sha1([]byte(password + pwd_salt))
		suc := dblayer.UserSignUp(username, enc_password)
		if suc {
			w.Write([]byte("SUCCESS"))
		} else {
			w.Write([]byte("FAILED"))
		}
	}
}