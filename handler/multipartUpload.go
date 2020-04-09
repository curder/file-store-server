package handler

import (
    "github.com/curder/file-store-server/utils"
    "math"
    "net/http"
    "strconv"
    "time"

    redisPool "github.com/curder/file-store-server/cache/redis"
)

// 初始化分块信息结构体
type MultipartUploadInfo struct {
    FileSha1   string // 文件hash
    FileSize   int    // 文件大小
    UploadID   string // 上传ID
    ChunkSize  int    // 分块大小
    ChunkCount int    // 分块数
}

// 初始化分块信息
func InitiateMultipartUploadHandler(w http.ResponseWriter, r *http.Request) {
    // 解析用户请求参数
    r.ParseForm()

    name := r.Form.Get("name")          // 用户名
    fileSha1 := r.Form.Get("file_sha1") // 文件hash值
    fileSize, err := strconv.Atoi(r.Form.Get("file_size"))

    if err != nil {
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
        UploadID:   name + string(time.Now().UnixNano()),
        ChunkSize:  5 * 1024 * 1024,                                       // 分块大小 = 5M
        ChunkCount: int(math.Ceil(float64(fileSize) / (5 * 1024 * 1024))), // fileSize/ChunkSize再向上取整
    }

    // 将初始化信息写入到Redis缓存
    _, _ = redisConnection.Do("HSET", "multipart_upload:"+multipartInfo.UploadID, "file_sha1", multipartInfo.FileSha1)
    _, _ = redisConnection.Do("HSET", "multipart_upload:"+multipartInfo.UploadID, "chunk_count", multipartInfo.ChunkCount)
    _, _ = redisConnection.Do("HSET", "multipart_upload:"+multipartInfo.UploadID, "file_size", multipartInfo.FileSize)

    // 将响应初始化数据返回给客户端
    w.Write(utils.NewResponse(0, "Successful", multipartInfo).JSONBytes())
}
