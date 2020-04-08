package main

import (
    "fmt"
    "github.com/curder/file-store-server/handler"
    "net/http"
)

func main() {
    //静态资源处理
    http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("./resources"))))

    http.HandleFunc("/files/uploads", handler.UploadHandler)                    // 文件上传处理
    http.HandleFunc("/files/uploads/succeeded", handler.UploadSucceededHandler) // 文件上传成功
    http.HandleFunc("/files/show", handler.GetFileMetaHandler)                  // 查询文件详情
    http.HandleFunc("/files/download", handler.DownloadHandler)                 // 文件下载
    http.HandleFunc("/files/update", handler.FileMateUpdateHandler)             // 文件更新 - 重命名
    http.HandleFunc("/files/delete", handler.FileDeleteHandler)                 // 文件删除

    http.HandleFunc("/users/sign-up", handler.SignUpHandler) // 用户注册
    http.HandleFunc("/users/sign-in", handler.SignInHandler) // 用户登录

    fmt.Println("http://127.0.0.1:8888")

    err := http.ListenAndServe(":8888", nil)
    if err != nil {
        fmt.Printf("Faild to start server, err: %s", err.Error())
    }
}
