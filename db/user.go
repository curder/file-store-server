package db

import (
    "fmt"
    mysqlDB "github.com/curder/file-store-server/db/mysql"
)

// 用户注册 - 通过用户名和密码
func UserSignUp(name, password string) bool {
    sqlStr := "INSERT INTO `users` (`name`, `password`) VAlUES (?, ?)"

    prepareStatement, err := mysqlDB.Connection().Prepare(sqlStr)
    if err != nil {
        fmt.Printf("Failed to insert,err: %s", err.Error())
        return false
    }

    defer prepareStatement.Close()

    result, err := prepareStatement.Exec(name, password)
    if err != nil {
        fmt.Printf("Failed to insert err: %s", err.Error())
        return false
    }

    if rowsAffected, err := result.RowsAffected(); err == nil && rowsAffected > 0 {
        return true
    }
    return false
}

// 用户登录 - 通过用户名和密码
func UserSignIn(name, password string) bool {
    sqlStr := "SELECT * FROM `users` WHERE name = ? LIMIT 1"

    prepareStatement, err := mysqlDB.Connection().Prepare(sqlStr)
    if err != nil {
        fmt.Printf("Failed prepare sql, err: %s", err.Error())
        return false
    }

    defer prepareStatement.Close()

    rows, err := prepareStatement.Query(name)
    if err != nil {
        fmt.Printf("query row err: %s", err.Error())
        return false
    } else if rows == nil {
        fmt.Printf("user not found: %s", name)
        return false
    }

    pRows := mysqlDB.ParseRows(rows)

    if len(pRows) > 0 && string(pRows[0]["password"].([]byte)) == password {
        return true
    }

    return false
}

// 新增或更新用户登录token
func UpdateToken(name, token string) bool {
    sqlStr := "REPLACE INTO `user_tokens` (`name`, `token`) VALUES (?, ?)"

    prepareStatement, err := mysqlDB.Connection().Prepare(sqlStr)
    if err != nil {
        fmt.Printf("Failed prepare sql error: %s", err.Error())
        return false
    }

    defer prepareStatement.Close()

    _, err = prepareStatement.Exec(name, token)
    if err != nil {
        fmt.Printf("Exec sql err: %s", err.Error())
        return false
    }

    return true
}
