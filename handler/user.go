package handler

import (
    "fmt"
    "github.com/curder/file-store-server/db"
    "github.com/curder/file-store-server/utils"
    "io/ioutil"
    "net/http"
)

const (
    passwordSalt = "1a!c2@b3#d4$e5f%g6^"
)

// 处理用户注册
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
            w.Write([]byte("SUCCESS"))
        } else {
            w.Write([]byte("FAILED"))
        }
    }
}
