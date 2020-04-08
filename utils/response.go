package utils

import (
    "encoding/json"
    "log"
)

type Response struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data"`
}

// 生成response对象
func (r *Response) NewResponse(code int, message string, data interface{}) *Response {
    return &Response{
        Code:    code,
        Message: message,
        Data:    data,
    }
}

// JSONBytes : 对象转json格式的二进制数组
func (r *Response) JSONBytes() (result []byte) {
    result, err := json.Marshal(r)
    if err != nil {
        log.Println(err)
    }
    
    return
}
