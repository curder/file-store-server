package utils

import (
    "crypto/md5"
    "crypto/sha1"
    "encoding/hex"
    "fmt"
    "github.com/curder/file-store-server/db"
    "hash"
    "io"
    "os"
    "path/filepath"
    "strconv"
    "time"
)

type Sha1Stream struct {
    _sha1 hash.Hash
}

// 更新文件hash
func (obj *Sha1Stream) Update(data []byte) {
    if obj._sha1 == nil {
        obj._sha1 = sha1.New()
    }
    obj._sha1.Write(data)
}

// 获取文件sum值
func (obj *Sha1Stream) Sum() string {
    return hex.EncodeToString(obj._sha1.Sum([]byte("")))
}

func Sha1(data []byte) string {
    _sha1 := sha1.New()
    _sha1.Write(data)
    return hex.EncodeToString(_sha1.Sum([]byte("")))
}

func FileSha1(file *os.File) string {
    _sha1 := sha1.New()
    io.Copy(_sha1, file)
    return hex.EncodeToString(_sha1.Sum(nil))
}

func MD5(data []byte) string {
    _md5 := md5.New()
    _md5.Write(data)
    return hex.EncodeToString(_md5.Sum([]byte("")))
}

func FileMD5(file *os.File) string {
    _md5 := md5.New()
    _, _ = io.Copy(_md5, file)
    return hex.EncodeToString(_md5.Sum(nil))
}

// 文件是否存在
func PathExists(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil {
        return true, nil
    }
    if os.IsNotExist(err) {
        return false, nil
    }
    return false, err
}

// 获取文件大小
func GetFileSize(filename string) int64 {
    var result int64
    _ = filepath.Walk(filename, func(path string, f os.FileInfo, err error) error {
        result = f.Size()
        return nil
    })
    return result
}

// 生成访问凭证 规则是 md5(name + timestamp + slat) + timestamp[:8]
func GenerateToken(name string) string {
    ts := fmt.Sprintf("%x", time.Now().Unix())

    tokenPrefix := MD5([]byte(name + ts + "token-salt"))

    return tokenPrefix + ts[:8]
}

// 检查token
func IsTokenValid(name, token string) bool {
    // 判断token是否过期
    if len(token) != 40 {
        return false
    }

    ts := fmt.Sprintf("%x", time.Now().Unix())[:8]
    userTime, _ := strconv.ParseUint(token[32:], 16, 32)
    now, _ := strconv.ParseUint(ts[:8], 16, 32)

    const TokenInvalidTime = 86400
    if userTime+TokenInvalidTime < now {
        fmt.Println("token超时")
        return false
    }

    // 数据库查询 name 对应的token信息
    isTokenUsed := db.CheckUserToken(name, token)
    if !isTokenUsed {
        fmt.Println("token错误")
        return false
    }

    return true
}
