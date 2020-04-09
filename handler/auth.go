package handler

import (
    "github.com/curder/file-store-server/utils"
    "net/http"
)

// 权限验证拦截器
func HttpInterceptor(h http.HandlerFunc) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 获取表单数据
        r.ParseForm()
        name := r.Form.Get("name")
        token := r.Form.Get("token")

        if len(name) < 3 || !utils.IsTokenValid(name, token) {
            w.WriteHeader(http.StatusForbidden)
            return
        }
        h(w, r)
    })
}
