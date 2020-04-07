package main

import (
    "fmt"
    "github.com/curder/file-store-server/handler"
    "net/http"
)

func main() {
    http.HandleFunc("/files/uploads", handler.UploadHandler)                   // 文件上传处理
    http.HandleFunc("/file/uploads/succeeded", handler.UploadSucceededHandler) // 文件上传成功

    fmt.Println("http://127.0.0.1:8888")

    err := http.ListenAndServe(":8888", nil)
    if err != nil {
        fmt.Printf("Faild to start server, err: %s", err.Error())
    }
}
