package db

import (
    "fmt"
    mysqlDB "github.com/curder/file-store-server/db/mysql"
    "time"
)

type UserFile struct {
    UserName  string
    FileName  string
    FileSha1  string
    FileSize  int64
    CreatedAt string
    UpdatedAt string
}

// 更新用户文件表
func OnUserFileUploadFinished(userName, fileName, fileSha1 string, fileSize int64) bool {
    // 预处理语句
    sqlStr := "REPLACE INTO `user_files` (`user_name`, `file_name`, `file_sha1`, `file_size`, `created_at`) VALUES (?, ?, ?, ?, ?);"
    prepareStatement, err := mysqlDB.Connection().Prepare(sqlStr)
    if err != nil {
        fmt.Printf("Failed to prepare sql err: %s", err.Error())
        return false
    }
    defer prepareStatement.Close()

    // 执行语句
    _, err = prepareStatement.Exec(userName, fileName, fileSha1, fileSize, time.Now().Format("2006-01-02 15:04:05"))
    if err != nil {
        fmt.Printf("Failed to exec sql err: %s", err.Error())
        return false
    }

    return true
}
