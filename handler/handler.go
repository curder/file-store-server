package handler

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

// 文件上传逻辑
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" { // // 展示文件上传表单
		// 加载模版文件
		tpl, err := ioutil.ReadFile("./resources/views/uploads.html")
		if err != nil {
			_, _ = io.WriteString(w, "internal server error")
			return
		}

		_, err = io.WriteString(w, string(tpl))

		if err != nil {
			fmt.Println("template error: ", err)
		}
	} else if r.Method == "POST" { // 处理文件上传
		file, head, err := r.FormFile("file") // 获取表单file字段的文件
		if err != nil {
			fmt.Printf("Failed to get data, err: %s", err.Error())
			return
		}
		defer func() { // 延迟关闭文件
			err = file.Close()
			if err != nil {
				fmt.Printf("close file err: %s", err.Error())
			}
		}()

		newFile, err := os.Create("/tmp/go-test-files/" + head.Filename) // 创建一个文件，如果存在这个文件则清空，如果不存在则创建

		defer func() { // 延迟关闭文件
			err = newFile.Close()
			if err != nil {
				fmt.Printf("new file close err: %s", err.Error())
			}
		}()

		if err != nil {
			fmt.Printf("Failed to create file, err: %s", err.Error())
			return
		}

		_, err = io.Copy(newFile, file) // 将用户上传的文件拷贝给上面创建的文件
		if err != nil {
			fmt.Printf("Failed to save data into file, err: %s", err.Error())
			return
		}

		http.Redirect(w, r, "/file/uploads/succeeded", http.StatusFound) // 重定向到新页面
	}
}

// 上传成功
func UploadSucceededHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = io.WriteString(w, "Upload finished!")
}
