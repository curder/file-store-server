package mysql

import (
    "database/sql"
    "fmt"
    _ "github.com/go-sql-driver/mysql"
    "log"
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

// 分析数据库行
func ParseRows(rows *sql.Rows) []map[string]interface{} {
    columns, _ := rows.Columns()
    scanArgs := make([]interface{}, len(columns))
    values := make([]interface{}, len(columns))
    for j := range values {
        scanArgs[j] = &values[j]
    }

    record := make(map[string]interface{})
    records := make([]map[string]interface{}, 0)
    for rows.Next() {
        // 将行数据保存到record字典
        if err := rows.Scan(scanArgs...); err != nil {
            log.Fatal(err)
            panic(err)
        }

        for i, col := range values {
            if col != nil {
                record[columns[i]] = col
            }
        }
        records = append(records, record)
    }
    return records
}
