package main

import (
    "encoding/json"
    "fmt"
    "time"
)

import (
    "code.google.com/p/log4go"
    "www.baidu.com/golang-lib/log"
    "www.baidu.com/golang-lib/module_state"    
)

type Output struct {
    Counters    module_state.StateTable // some counters
}


func main() {
    str := "this is test log for test test test test test"
    
    log4go.SetLogBufferLength(10000)
    log4go.SetLogWithBlocking(false)
    log4go.SetWithModuleState(true)
    log.Init("test", "DEBUG", "./log", false, "D", 5)

    var total, count int64
    var output Output
    total = 0
    count = 0
    for i := 1; i < 1000000; i = i + 1 {
        start := time.Now()
        
        log.Logger.Info(str)

        now := time.Now()
        /* get duration from start to now, in Microsecond   */
        duration := now.Sub(start).Nanoseconds() / 1000
        
        total += duration
        count += 1
        
        if (count % 100) == 0 {
            fmt.Printf("ave=%d\n", total/count)
            total = 0
            count = 0
            output.Counters = log4go.GetModuleState()
            buff, err := json.Marshal(output)
            if err == nil {
                fmt.Printf("%s\n", buff)
            }
        }
    }
}