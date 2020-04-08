package mysql

import (
    "database/sql"
    "fmt"
    _ "github.com/go-sql-driver/mysql"
    "os"
)

var db *sql.DB

// 初始化数据库连接
func init() {
    var err error
    dns := "root:@tcp(127.0.0.1:3306)/file_store_server?charset=utf8"
    db, err = sql.Open("mysql", dns)
    db.SetConnMaxLifetime(1000) // 设置活跃连接数
    if err != nil {
        fmt.Printf("Failed to connect to mysql, err:%s", err.Error())
        os.Exit(1)
    }
}

// 返回数据库连接对象
func Connection() *sql.DB {
    return db
}
