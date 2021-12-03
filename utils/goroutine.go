package utils

import (
	log "github.com/sirupsen/logrus"
)

func SafeGroutine(fn func(options ...interface{})){
	defer func() {
		if err:= recover(); err != nil{
			log.Errorf("未知错误, error:%+v", err)
		}
	}()
	fn()
}

// WrapRecover goroutine里的panic都最好被处理掉,否则会造成整个程序的崩溃.
func WrapRecover(){
	if err:= recover(); err != nil {
		log.Errorf("未知错误, error:%+v", err)
	}
}