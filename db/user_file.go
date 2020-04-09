package db

import (
    "fmt"
    mysqlDB "github.com/curder/file-store-server/db/mysql"
    "time"
)

type UserFile struct {
    UserName  string `json:"user_name"`
    FileName  string  `json:"file_name"`
    FileSha1  string  `json:"file_sha1"`
    FileSize  int64   `json:"file_size"`
    CreatedAt string  `json:"created_at"`
    UpdatedAt string  `json:"updated_at"`
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

// 获取用户文件
func GetUserFileMetas(username string, limit int) ([]UserFile, error) {
    sqlStr := "SELECT `user_name`, `file_name`, `file_sha1`, `file_size`, `created_at`, `updated_at` FROM user_files WHERE user_name = ? LIMIT ?"
    prepareStatement, err := mysqlDB.Connection().Prepare(sqlStr)
    if err != nil {
        fmt.Printf("Failed to prepare sql err: %s", err.Error())
        return nil, err
    }
    defer prepareStatement.Close()
    rows, err := prepareStatement.Query(username, limit)

    var userFiles []UserFile

    for rows.Next() {
        userFile := UserFile{}
        err := rows.Scan(&userFile.UserName, &userFile.FileName, &userFile.FileSha1, &userFile.FileSize, &userFile.CreatedAt, &userFile.UpdatedAt)
        if err != nil {
            fmt.Println(err.Error())
            break
        }
        userFiles = append(userFiles, userFile)
    }
    return userFiles, nil
}
