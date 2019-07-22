package util

import (
	"sync"
	"time"

	"github.com/ouqiang/timewheel"
)

var lk sync.RWMutex
var Tw *timewheel.TimeWheel
var mapDelayFunc map[string]func(interface{})

type DelayParameter struct {
	Key  string      //key是map的key
	Data interface{} //Data是延时处理函数传入的参数
}

func init() {
	mapDelayFunc = make(map[string]func(interface{}))

	Tw = timewheel.New(1*time.Second, 1800, delayHandle)
	Tw.Start()
}

func delayHandle(data interface{}) {
	parameter, ok := data.(*DelayParameter)
	if !ok {
		return
	}

	lk.RLock()
	fun, ok := mapDelayFunc[parameter.Key]
	if !ok {
		lk.RUnlock()
		return
	}
	lk.RUnlock()

	fun(parameter.Data)
}

// f中interface{}参数为DelayParameter.Data
func AddDelayFunc(funcKey string, f func(interface{})) {
	lk.Lock()
	defer lk.Unlock()

	fun, ok := mapDelayFunc[funcKey]
	if ok && fun != nil {
		return
	}

	mapDelayFunc[funcKey] = f
}

// key是map的key，Data是延时处理函数传入的参数
func GenerateDelayParameter(key string, data interface{}) *DelayParameter {
	return &DelayParameter{
		Key:  key,
		Data: data,
	}
}
