package hkit

/****************

Author: yuansudong
Date: 2018.06.23
*****************/

import "time"

const (
	//  _DefaultRetryTime 默认的重试次数
	_DefaultRetryCount = 3
	// _DefaultRetryTime 默认的重试时间
	_DefaultRetryTime = time.Second * 5
)

const (
	// _HeaderContentType 类型
	_HeaderKeyContentType = "Content-Type"
	// _HeaderValueApplicationJSON
	_HeaderValueApplicationJSON = "application/json"
)

const (
	// _LayoutError 错误格式
	_LayoutError = "unknown http code:%d,info:%s"
)
