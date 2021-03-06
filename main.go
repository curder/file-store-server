package main

import (
    "fmt"
    "github.com/curder/file-store-server/handler"
    "log"
    "net/http"
)

func main() {
    //静态资源处理
    http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("./resources"))))

    http.HandleFunc("/files/uploads", handler.UploadHandler)                    // 文件上传处理
    http.HandleFunc("/files/uploads/succeeded", handler.UploadSucceededHandler) // 文件上传成功
    http.HandleFunc("/files/index", handler.GetFilesMetaHandler)                // 文件列表
    http.HandleFunc("/files/show", handler.GetFileMetaHandler)                  // 查询文件详情
    http.HandleFunc("/files/download", handler.DownloadHandler)                 // 文件下载
    http.HandleFunc("/files/update", handler.FileMateUpdateHandler)             // 文件更新 - 重命名
    http.HandleFunc("/files/delete", handler.FileDeleteHandler)                 // 文件删除
    http.HandleFunc("/files/fast-uploads", handler.TypeFastUploadHandler)       // 秒传

    // 分块上传
    http.HandleFunc("/files/multipart-uploads/init", handler.InitiateMultipartUploadHandler) // 初始化分块信息
    http.HandleFunc("/files/multipart-uploads/upload-part", handler.UploadPartHandler)       // 分块上传
    http.HandleFunc("/files/multipart-uploads/complete", handler.CompleteUploadPartHandler)  // 分块上传完毕

    http.HandleFunc("/users/sign-up", handler.SignUpHandler)                         // 用户注册
    http.HandleFunc("/users/sign-in", handler.SignInHandler)                         // 用户登录
    http.HandleFunc("/users/info", handler.HttpInterceptor(handler.UserInfoHandler)) // 用户详细信息

    fmt.Println("http://127.0.0.1:8888")

    log.Fatal(http.ListenAndServe(":8888", nil))
}
