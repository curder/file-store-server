package redis

import (
    "fmt"
    "github.com/garyburd/redigo/redis"
    "time"
)

var (
    pool          *redis.Pool
    redisHost     = "127.0.0.1:63790"
    redisPassword = ""
)

// 创建Redis连接池
func newRedisPool() *redis.Pool {
    return &redis.Pool{
        Dial: func() (conn redis.Conn, err error) {
            // 打开连接
            dial, err := redis.Dial("tcp", redisHost)
            if err != nil {
                fmt.Printf(err.Error())
                return nil, err
            }
            // 访问认证
            if _, err = dial.Do("AUTH", redisPassword); err != nil {
                dial.Close()
                return nil, err
            }

            return dial, nil
        },
        TestOnBorrow: func(c redis.Conn, t time.Time) error {
            if time.Since(t) < time.Minute {
                return nil
            }

            _, err := c.Do("PING")

            return err
        },
        MaxIdle:         100,
        MaxActive:       100,
        IdleTimeout:     300 * time.Second,
        Wait:            false,
        MaxConnLifetime: 0,
    }
}

func init() {
    pool = newRedisPool()
}

// Redis连接
func RedisPool() *redis.Pool {
    return pool
}
