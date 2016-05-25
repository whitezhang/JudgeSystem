/* log_test.go - test for log.go */
/*
modification history
--------------------
2014/3/7, by Zhang Miao, create
2014/3/11, by Zhang Miao, modify
*/
package log

import (
	"testing"
	//"time"
)

func TestLog(t *testing.T) {
	if err := Init("test", "INFO", "./log/log", true, "M", 2); err != nil {
		t.Error("log.Init() fail")
	}

	if err := Init("test", "INFO", "./log/log", true, "M", 5); err == nil {
		t.Error("fail in process reentering log.Init()")
	}

	for i := 0; i < 100; i = i + 1 {
		Logger.Warn("warning msg: %d", i)
		Logger.Info("info msg: %d", i)

		// time.Sleep(10 * time.Second)
	}

	//time.Sleep(100 * time.Millisecond)
}
