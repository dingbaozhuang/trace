package service

import (
	"fmt"
	"time"

	"github.com/ouqiang/timewheel"
)

var Tw *timewheel.TimeWheel

func DelayInit() {
	Tw = timewheel.New(1*time.Second, 3600, func(data interface{}) {
		fmt.Println("time wheel ~", data)
	})

	Tw.Start()
}
