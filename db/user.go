package db

import (
    "fmt"
    mysqlDB "github.com/curder/file-store-server/db/mysql"
)

// 用户注册 - 通过用户名和密码
func UserSignUp(username, password string) bool {
    sqlStr := "INSERT INTO `users` (`name`, `password`) VAlUES (?, ?)"

    prepareStatement, err := mysqlDB.Connection().Prepare(sqlStr)
    if err != nil {
        fmt.Printf("Failed to insert,err: %s", err.Error())
        return false
    }

    defer prepareStatement.Close()

    result, err := prepareStatement.Exec(username, password)
    if err != nil {
        fmt.Printf("Failed to insert err: %s", err.Error())
        return false
    }

    if rowsAffected, err := result.RowsAffected(); err == nil && rowsAffected > 0 {
        return true
    }
    return false
}
