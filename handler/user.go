package handler

import (
	"net/http"
	"io/ioutil"
	"filestore-server/util"
	dblayer "filestore-server/db"
	"time"
	"fmt"
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

// 登录接口
func SignInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		data, err := ioutil.ReadFile("./static/view/signin.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(data)
		return
	}
	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	enc_password := util.Sha1([]byte(password + pwd_salt))
	// 1. 校验用户名及密码
	pwdChecked := dblayer.UserSignIn(username, enc_password)
	if !pwdChecked {
		w.Write([]byte("FAILED"))
		return
	}
	// 2. 生成访问凭证(Token)
	token := GenToken(username)
	tkChecked := dblayer.UpdateToken(username, token)
	if !tkChecked {
		w.Write([]byte("FAILED"))
		return
	}

	// 3. 登录成功后重定向到首页
	// w.Write([]byte("http://"+ r.Host +"/static/view/home.html"))
	resp := util.RespMsg{
		Code: 0,
		Msg: "ok",
		Data: struct {
			Location string
			Username string
			Token string
		}{
			Location: "/static/view/home.html",
			Username: username,
			Token: token,
		},
	}
	w.Write(resp.JSONBytes())
}

// 查询用户信息
func UserInfoHandler(w http.ResponseWriter, r *http.Request) {
	// 1. 解析请求参数
	r.ParseForm()
	username := r.Form.Get("username")
	token := r.Form.Get("token")
	// 2. 验证 Token 是否有效
	isValidTk := IsTokenVaild(token)
	if !isValidTk {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	// 3. 查询用户信息
	user, err := dblayer.GetUserInfo(username)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	// 4. 组装并且响应用户数据
	resp := util.RespMsg{
		Code: 0,
		Msg: "ok",
		Data: user,
	}
	w.Write(resp.JSONBytes())
}

// 生成访问凭证
func GenToken(username string) string {
	// 40位字符：md5(username + timestamp+token_salt) + timestamp[:8]
	ts := fmt.Sprintf("%x", time.Now().Unix())
	token_prefix := util.MD5([]byte(username + ts + "_tokensalt"))
	return token_prefix + ts[:8]
}

// Token 是否有效
func IsTokenVaild(token string) bool {
	// TODO: 从数据表 tbl_user_token 查询 username 对应的 Token 信息
	// TODO: 判断 Token 的时效性
	// TODO: 对比两个 Token 是否一致
	return true
}