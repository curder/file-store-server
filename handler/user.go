package handler

import (
    "encoding/json"
    "fmt"
    "github.com/curder/file-store-server/db"
    "github.com/curder/file-store-server/utils"
    "io/ioutil"
    "net/http"
)

const (
    passwordSalt = "1a!c2@b3#d4$e5f%g6^"
)

// 用户注册
func SignUpHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        tpl, err := ioutil.ReadFile("./resources/views/sign-up.html")
        if err != nil {
            fmt.Printf("Failed to read template file err: %s", err.Error())
            return
        }
        w.Write(tpl)
    } else if r.Method == "POST" { // 注册逻辑
        name := r.PostFormValue("name")         // 操作类型 0 为重命名
        password := r.PostFormValue("password") // 修改的文件sha1

        if len(name) < 3 || len(password) < 5 {
            w.Write([]byte("Invalid parameter"))
            return
        }

        encodePassword := utils.Sha1([]byte(password + passwordSalt))

        result := db.UserSignUp(name, encodePassword)

        if result {
            res := map[string]interface{}{
                "code":         10000,
                "redirect_url": "/users/sign-in",
            }
            data, _ := json.Marshal(res)
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusOK)
            w.Write(data)
        } else {
            w.Write([]byte("FAILED"))
        }
    }
}

// 用户登录
func SignInHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        tpl, err := ioutil.ReadFile("./resources/views/sign-in.html")
        if err != nil {
            fmt.Printf("Failed error to read file %s", err.Error())
            return
        }
        w.Write(tpl)
    } else if r.Method == "POST" {
        // 获取用户提交表单数据
        name := r.PostFormValue("name")
        password := r.PostFormValue("password")

        if len(name) < 3 || len(password) < 5 {
            w.Write([]byte("Invalid params"))
            return
        }

        // 查询用户
        hasUser := db.UserSignIn(name, utils.Sha1([]byte(password+passwordSalt)))
        if !hasUser { // 存在用户
            w.Write([]byte("FAILED"))
            return
        }

        // 生成访问凭证
        token := utils.GenerateToken(name)
        if updateToken := db.UpdateToken(name, token); !updateToken {
            w.Write([]byte("FAILED"))
            return
        }

        // 构造返回的结构体，返回json
        response := utils.Response{
            Code:    0,
            Message: "登录成功",
            Data: struct {
                Location string `json:"location"`
                Name     string `json:"name"`
                Token    string `json:"token"`
            }{
                Location: "http://" + r.Host + "/resources/views/home.html",
                Name:     name,
                Token:    token,
            },
        }
        w.Header().Set("Content-Type", "application/json")
        w.Write(response.JSONBytes())
    }
}

// 查询用户信息
func UserInfoHandler(w http.ResponseWriter, r *http.Request) {
    // 解析请求参数
    r.ParseForm()
    name := r.Form.Get("name")
    // token := r.Form.Get("token")

    // 验证 token 是否生效
    //if isValid := utils.IsTokenValid(name, token); !isValid {
    //    fmt.Printf("token is invalid")
    //    return
    //}

    // 查询用户信息
    user, err := db.GetUserInfo(name)
    if err != nil {
        w.WriteHeader(http.StatusForbidden)
        return
    }

    // 组装并响应用户数据
    response := utils.Response{
        Code:    0,
        Message: "获取用户详细信息",
        Data:    user,
    }

    w.Write(response.JSONBytes())
}
