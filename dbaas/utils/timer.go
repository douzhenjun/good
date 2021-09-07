package utils

import "time"

func TimeStamp() int64 {
	return time.Now().Unix()
}

/*
Timer 一个简单的定时器: 延迟执行action
返回stop方法, 可在外部停止定时器
*/
func Timer(action func(), delay time.Duration) func() {
	var stopped bool
	var stopC = make(chan struct{}, 1)
	go func() {
		defer func() { stopped = true; close(stopC) }()
		select {
		case <-stopC:
			return
		case <-time.After(delay):
			action()
		}
	}()
	return func() {
		if stopped {
			return
		}
		stopC <- struct{}{}
	}
}

/*
Polling 轮询定时器: 固定间隔执行action, 如果action返回true停止轮询, count为轮询次数
返回stop方法, 可在外部停止轮询定时器
*/
func Polling(action func() bool, duration time.Duration, count int) func() {
	var stopped bool
	stopC := make(chan struct{}, 1)
	go func() {
		defer func() { stopped = true; close(stopC) }()
		ticker := time.NewTicker(duration)
		for i := 0; i < count; i++ {
			select {
			case <-stopC:
				return
			case <-ticker.C:
				if action() {
					return
				}
			}
		}
	}()
	return func() {
		if stopped {
			return
		}
		stopC <- struct{}{}
	}
}

/*
LoopTask 定时任务:  固定间隔执行action
返回stop方法, 可在外部停止定时任务
*/
func LoopTask(action func(), duration time.Duration) func() {
	var stopped bool
	stopC := make(chan struct{}, 1)
	go func() {
		defer func() { stopped = true; close(stopC) }()
		ticker := time.NewTicker(duration)
		for {
			select {
			case <-stopC:
				return
			case <-ticker.C:
				action()
			}
		}
	}()
	return func() {
		if stopped {
			return
		}
		stopC <- struct{}{}
	}
}
