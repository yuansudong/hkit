package hkit

/****************
Author: yuansudong
Date: 2018.06.23
*****************/

import (
	"net/http"
	"time"
)

// Request 用于描述一个请求
type Request struct {
	Method string
	URL    string
	Body   []byte
	Header []HTTPOption
}

// Response 用于描述一个响应
type Response struct {
	Status int
	Header http.Header
	Body   []byte
}

// RetryConfig 重试配置
type RetryConfig struct {
	_MaxRetryTimes int
	_RetryInterval time.Duration
}

// Engine 用于描述一个httpclient
type Engine struct {
	_Client     *http.Client
	_Opt        *Option
	_BaseHeader []HTTPOption
}

// HTTPOption 定义一个http的选项
type HTTPOption func(req *http.Request)

// SetHeader 设置头部
func SetHeader(key string, value string) HTTPOption {
	return func(r *http.Request) {
		r.Header.Set(key, value)
	}
}
