package db

import (
    "database/sql"
    "fmt"
    mysqlDB "github.com/curder/file-store-server/db/mysql"
)

// 文件上传完毕写入数据库行操作
func OnFileUploadFinished(fileName, fileSha1 string, fileSize int64, location string) bool {
    sqlStr := "INSERT INTO files (`name`, `sha1`, `size`, `path`, `status`) VALUES (?, ?, ?, ?, 1);"
    prepareStatement, err := mysqlDB.Connection().Prepare(sqlStr)
    if err != nil {
        fmt.Printf("Failed to prepare statement. err: %s", err.Error())
        return false
    }

    defer prepareStatement.Close() // 关闭连接

    result, err := prepareStatement.Exec(fileName, fileSha1, fileSize, location)
    if err != nil {
        fmt.Printf("Failed to exec sql err: %s", err.Error())
        return false
    }

    if res, err := result.RowsAffected(); err == nil {
        if res <= 0 {
            fmt.Printf("File with hash: %s has been uploaded before", fileSha1)
        }
        return true
    }
    return false
}

type TableFile struct {
    FileSha1 string
    FileName sql.NullString
    FileSize sql.NullInt64
    Location sql.NullString
}

// 通过Sha1获取文件原信息
func GetFileMeta(fileSha1 string) (*TableFile, error) {
    sqlStr := "SELECT `sha1`, `name`, `size`, `path` FROM files WHERE sha1= ? AND `status` = 1 limit 1"
    prepareStatement, err := mysqlDB.Connection().Prepare(sqlStr)
    if err != nil {
        fmt.Printf("Failed to prepare sql err: %s", err.Error())
        return nil, err
    }

    defer prepareStatement.Close() // 关闭连接

    file := TableFile{}

    err = prepareStatement.QueryRow(fileSha1).Scan(&file.FileSha1, &file.FileName, &file.FileSize, &file.Location)
    if err != nil {
        fmt.Printf(err.Error())
        return nil, err
    }

    return &file, nil
}
