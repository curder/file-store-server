package handler

import (
    "fmt"
    "github.com/curder/file-store-server/db"
    "github.com/curder/file-store-server/utils"
    "github.com/garyburd/redigo/redis"
    "math"
    "net/http"
    "os"
    "path"
    "strconv"
    "strings"
    "time"

    redisPool "github.com/curder/file-store-server/cache/redis"
)

// 初始化分块信息结构体
type MultipartUploadInfo struct {
    FileSha1   string `json:"file_sha1"`   // 文件hash
    FileSize   int    `json:"file_size"`   // 文件大小
    UploadID   string `json:"upload_id"`   // 上传ID
    ChunkSize  int    `json:"chunk_size"`  // 分块大小
    ChunkCount int    `json:"chunk_count"` // 分块数
}

// 初始化分块信息
func InitiateMultipartUploadHandler(w http.ResponseWriter, r *http.Request) {
    // 解析用户请求参数
    r.ParseForm()

    name := r.Form.Get("name")          // 用户名
    fileSha1 := r.Form.Get("file_hash") // 文件hash值
    fileSize, err := strconv.Atoi(r.Form.Get("file_size"))

    if err != nil {
        w.Header().Set("Content-Type", "application/json")
        _, _ = w.Write(utils.NewResponse(-1, "Invalid params", nil).JSONBytes())
        return
    }

    // 获取Redis连接
    redisConnection := redisPool.RedisPool().Get()
    defer redisConnection.Close()

    // 生成分块上传的初始化信息
    multipartInfo := MultipartUploadInfo{
        FileSha1:   fileSha1,
        FileSize:   fileSize,
        UploadID:   name + fmt.Sprintf("%x", time.Now().UnixNano()),
        ChunkSize:  5 * 1024 * 1024,                                       // 分块大小 = 5M
        ChunkCount: int(math.Ceil(float64(fileSize) / (5 * 1024 * 1024))), // fileSize/ChunkSize再向上取整
    }

    // 将初始化信息写入到Redis缓存
    _, _ = redisConnection.Do("HSET", "multipart_upload:"+multipartInfo.UploadID, "file_sha1", multipartInfo.FileSha1)
    _, _ = redisConnection.Do("HSET", "multipart_upload:"+multipartInfo.UploadID, "chunk_count", multipartInfo.ChunkCount)
    _, _ = redisConnection.Do("HSET", "multipart_upload:"+multipartInfo.UploadID, "file_size", multipartInfo.FileSize)

    // 将响应初始化数据返回给客户端
    w.Header().Set("Content-Type", "application/json")
    w.Write(utils.NewResponse(0, "Successful", multipartInfo).JSONBytes())
}

// 分块上传
func UploadPartHandler(w http.ResponseWriter, r *http.Request) {
    // 解析用户请求参数
    r.ParseForm()

    uploadID := r.Form.Get("upload_id")     // 上传的id
    chunkIndex := r.Form.Get("chunk_index") // 分块索引

    // 获得Redis连接
    redisConnection := redisPool.RedisPool().Get()
    defer redisConnection.Close()

    // 获得文件句柄，用于存储分块内容
    filePath := "/tmp/go-test-files/multipart-uploads/" + uploadID + "/" + chunkIndex
    err := os.MkdirAll(path.Dir(filePath), 0744) // 创建目录
    if err != nil {
        fmt.Println(err.Error())
        w.Header().Set("Content-Type", "application/json")
        w.Write(utils.NewResponse(-1, "make dir failed", nil).JSONBytes())
    }

    fd, err := os.Create(filePath)
    if err != nil {
        w.Header().Set("Content-Type", "application/json")
        w.Write(utils.NewResponse(-1, "Upload part failed", nil).JSONBytes())
        return
    }
    defer fd.Close()

    // 每次读取 5M
    buffer := make([]byte, 5 * 1024*1024)
    for {
        n, err := r.Body.Read(buffer)
        fd.Write(buffer[:n])
        if err != nil {
            break
        }
    }
    // 更新redis缓存状态
    _, _ = redisConnection.Do("HSET", "multipart_upload:"+uploadID, "chunk_index:"+chunkIndex, 1)

    // 返回处理结果到客户端
    w.Header().Set("Content-Type", "application/json")
    w.Write(utils.NewResponse(0, "Successful", nil).JSONBytes())
}

// 分块上传完毕
func CompleteUploadPartHandler(w http.ResponseWriter, r *http.Request) {
    // 解析请求参数
    r.ParseForm()

    name := r.Form.Get("name")                           // 用户名
    uploadId := r.Form.Get("upload_id")                  // 文件上传ID
    fileSha1 := r.Form.Get("file_hash")                  // 文件hash
    fileSize, _ := strconv.Atoi(r.Form.Get("file_size")) // 文件大小
    fileName := r.Form.Get("file_name")                  // 文件名称

    // 获取redis连接
    redisConnection := redisPool.RedisPool().Get()
    defer redisConnection.Close()

    // 通过upload_id判断是否所有分块上传完毕
    data, err := redis.Values(redisConnection.Do("HGETALL", "multipart_upload:"+uploadId))
    if err != nil {
        w.Header().Set("Content-Type", "application/json")
        w.Write(utils.NewResponse(-1, "complete upload failed, err:" + err.Error(), nil).JSONBytes())
        return
    }
    totalCount := 0 // 总分块数
    chunkCount := 0 // 分块数量

    for i := 0; i < len(data); i += 2 {
        k := string(data[i].([]byte))
        v := string(data[i+1].([]byte))

        if k == "chunk_count" {
            totalCount, _ = strconv.Atoi(v)
        } else if strings.HasPrefix(k, "chunk_index:") && v == "1" {
            chunkCount += 1
        }
    }

    if totalCount != chunkCount {
        w.Header().Set("Content-Type", "application/json")
        w.Write(utils.NewResponse(-2, "invalid request", nil).JSONBytes())
        return
    }

    // TODO 合并分块

    // 更新文件表和用户关联文件表
    db.OnFileUploadFinished(fileName, fileSha1, int64(fileSize), "")
    db.OnUserFileUploadFinished(name, fileName, fileSha1, int64(fileSize))

    // 向客户端响应结果
    w.Header().Set("Content-Type", "application/json")
    w.Write(utils.NewResponse(0, "successful", nil).JSONBytes())
}
