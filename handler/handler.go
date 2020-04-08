package handler

import (
	"encoding/json"
	"fmt"
	"github.com/curder/file-store-server/meta"
	"github.com/curder/file-store-server/utils"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
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

		// 创建文件原始数据对象
		fileMeta := meta.FileMeta{
			FileName:  head.Filename,
			Location:  "/tmp/go-test-files/" + head.Filename,
			UpdatedAt: time.Now().Format("2006-01-02 15:04:05"),
		}

		newFile, err := os.Create(fileMeta.FileName) // 创建一个文件，如果存在这个文件则清空，如果不存在则创建

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

		fileMeta.FileSize, err = io.Copy(newFile, file) // 将用户上传的文件拷贝给上面创建的文件
		if err != nil {
			fmt.Printf("Failed to save data into file, err: %s", err.Error())
			return
		}

		_, _ = newFile.Seek(0, 0)
		fileMeta.FileSha1 = utils.FileSha1(newFile) // 计算文件的sha1值
		meta.UpdateFileMeta(fileMeta)

		http.Redirect(w, r, "/files/uploads/succeeded", http.StatusFound) // 重定向到新页面
	}
}

// 上传成功
func UploadSucceededHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = io.WriteString(w, "Upload finished!")
}

// 文件查询
func GetFileMetaHandle(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Printf("parse form err: %s", err.Error())
		return
	}

	fileHash := r.Form["file_hash"][0]

	fileMeta := meta.GetFileMeta(fileHash)

	data, err := json.Marshal(fileMeta) // json格式化

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(data)
}

// 文件下载
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Printf("parse form err: %s", err.Error())
		return
	}

	fileHash := r.Form.Get("file_hash")
	fileMeta := meta.GetFileMeta(fileHash)

	file, err := os.Open(fileMeta.Location) // 打开本地文件
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer func() {
		err = file.Close()
		if err != nil {
			fmt.Printf("Close file err: %s", err.Error())
		}
	}()

	data, err := ioutil.ReadAll(file) // 读取文件
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	// 添加下载的头信息
	w.Header().Set("Content-Type", "application/octect-stream")
	w.Header().Set("Content-Disposition", `attachment;fileName="`+fileMeta.FileName+`"`)
	w.Write(data)
}
