package main

import (
    "bufio"
    "bytes"
    "fmt"
    "io"
    "io/ioutil"
    "net/http"
    "net/url"
    "os"
    "strconv"

    jsonit "github.com/json-iterator/go"
)

func multipartUpload(filename string, targetURL string, chunkSize int) error {
    f, err := os.Open(filename)
    if err != nil {
        fmt.Println(err)
        return err
    }
    defer f.Close()

    bfRd := bufio.NewReader(f)
    index := 0

    ch := make(chan int)
    buf := make([]byte, chunkSize) // 每次读取chunkSize大小的内容
    for {
        n, err := bfRd.Read(buf)
        if n <= 0 {
            break
        }
        index++

        bufCopied := make([]byte, 5*1024*1024)
        copy(bufCopied, buf)

        go func(b []byte, curIdx int) {
            fmt.Printf("upload_size: %d\n", len(b))

            resp, err := http.Post(
                targetURL+"&chunk_index="+strconv.Itoa(curIdx),
                "multipart/form-data",
                bytes.NewReader(b))
            if err != nil {
                fmt.Println(err)
            }

            body, er := ioutil.ReadAll(resp.Body)
            fmt.Printf("%+v %+v\n", string(body), er)
            resp.Body.Close()

            ch <- curIdx
        }(bufCopied[:n], index)

        //遇到任何错误立即返回，并忽略 EOF 错误信息
        if err != nil {
            if err == io.EOF {
                break
            } else {
                fmt.Println(err.Error())
            }
        }
    }

    for idx := 0; idx < index; idx++ {
        select {
        case res := <-ch:
            fmt.Println(res)
        }
    }

    return nil
}

func main() {
    name := "curder"
    token := "96f8c8b042ad6e99945fde8a5edd380e5e8ed82e"
    fileHash := "fd7c5327c68fcf94b62dc9f58fc1cdb3c8c01258"

    // 1. 请求初始化分块上传接口
    resp, err := http.PostForm(
        "http://127.0.0.1:8888/files/multipart-uploads/init",
        url.Values{
            "name":     {name},
            "token":    {token},
            "file_hash": {fileHash},
            "file_size": {"209715200"},
        })

    if err != nil {
        fmt.Println(err.Error())
        os.Exit(-1)
    }

    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Println(err.Error())
        os.Exit(-1)
    }

    // 2. 得到uploadID以及服务端指定的分块大小chunkSize
    uploadID := jsonit.Get(body, "data").Get("upload_id").ToString()
    chunkSize := jsonit.Get(body, "data").Get("chunk_size").ToInt()
    fmt.Printf("upload_id: %s  chunk_size: %d\n", uploadID, chunkSize)

    // 3. 请求分块上传接口
    fileName := "/tmp/go-test-files/go1.10.3.linux-amd64.tar.gz"
    tURL := "http://127.0.0.1:8888/files/multipart-uploads/upload-part?" +
        "name=" + name + "&token=" + token + "&upload_id=" + uploadID
    err = multipartUpload(fileName, tURL, chunkSize)
    if err != nil {
        fmt.Println(err.Error())
    }


    // 4. 请求分块完成接口
    resp, err = http.PostForm(
        "http://127.0.0.1:8888/files/multipart-uploads/complete",
        url.Values{
            "name":      {name},
            "upload_id": {uploadID},
            "token":     {token},
            "file_hash": {fileHash},
            "file_size": {"209715200"},
            "file_name":  {"go1.10.3.linux-amd64.tar.gz"},
        })

    if err != nil {
        fmt.Println(err.Error())
        os.Exit(-1)
    }

    defer resp.Body.Close()
    body, err = ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Println(err.Error())
        os.Exit(-1)
    }
    fmt.Printf("complete result: %s\n", string(body))


}
