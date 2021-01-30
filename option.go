package hkit

/****************
Author: yuansudong
Date: 2018.06.23
*****************/

import "time"

// Option 用于描述一个选项配置
type Option struct {
	MaxIdleConns        int           `yaml:"max_idle_conns"`
	IdleConnTimeout     int           `yaml:"idle_conn_timeout"`
	MaxIdleConnsPerHost int           `yaml:"max_idle_conn_per_host"`
	DisableCompression  bool          `yaml:"disable_compress"`
	SkipTLSVerify       bool          `yaml:"skip_tls_verify"`
	RetryInterval       int           `yaml:"retry_interval"`
	_RetryTime          time.Duration `yaml:"-"`
	RetryCount          int           `yaml:"retry_count"`
	EnableProxy         bool          `yaml:"enable_proxy"`
	ProxyAddr           string        `yaml:"proxy_addr"`
	EnableHTTP2         bool          `yaml:"enable_http2"`
}
