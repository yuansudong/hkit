package hkit

import "log"

var (
	_EngineInst *Engine
)

// Init 初始化
func Init(iOpt *Option) {
	var err error
	_EngineInst, err = _New(iOpt)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

// Get 获取单例
func Get() *Engine {
	return _EngineInst
}

// Release 释放
func Release() {
	if _EngineInst != nil {
		_EngineInst._Client.CloseIdleConnections()
	}
	_EngineInst = nil
}
