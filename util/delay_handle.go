package util

import (
	"os"
	"time"

	"github.com/ouqiang/timewheel"
)

var Tw *timewheel.TimeWheel

func DelayInit() {
	Tw = timewheel.New(1*time.Second, 1800, delayHandle)

	Tw.Start()
}

func delayHandle(data interface{}) {
	os.Remove(data.(string))
}
