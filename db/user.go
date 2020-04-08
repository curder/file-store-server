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

// 数据库查询user对应的token信息
func CheckUserToken(name, token string) bool {
    sqlStr := "SELECT * FROM `user_tokens` WHERE name = ? LIMIT 1"
    stmt, err := mysqlDB.Connection().Prepare(sqlStr)
    if err != nil {
        fmt.Println(err.Error())
        return false
    }

    defer stmt.Close()
    rows, err := stmt.Query(name)
    if err != nil {
        fmt.Println(err.Error())
        return false
    } else if rows == nil {
        fmt.Println("user name not found " + name)
        return false
    }

    pRows := mysqlDB.ParseRows(rows)

    if len(pRows) > 0 && string(pRows[0]["token"].([]byte)) == token {
        return true
    }
    return false
}

// 用户结构体
type User struct {
    Name       string `json:"name"`
    Email      string `json:"email"`
    Phone      string `json:"phone"`
    SignUpAt   string `json:"sign_up_at"`
    LastActive string `json:"last_active"`
    Status     int    `json:"status"`
}

// 查询用户信息
func GetUserInfo(name string) (user User, err error) {
    u := User{}

    sqlStr := "SELECT `name`, `email`, `phone`, `sign_up_at`, `last_active`, `status` FROM `users` WHERE `name` = ? LIMIT 1"
    prepareStatement, err := mysqlDB.Connection().Prepare(sqlStr)
    if err != nil {
        fmt.Printf("prepare sql err: %s", err.Error())
        return u, err
    }
    defer prepareStatement.Close()

    err = prepareStatement.QueryRow(name).Scan(&u.Name, &u.Email, &u.Phone, &u.SignUpAt, &u.LastActive, &u.Status)
    if err != nil {
        return u, err
    }

    return u, nil
}
