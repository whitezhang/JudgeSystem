/* module_state2_test.go - test for module_state.go  */
/*
modification history
--------------------
2014/3/25, by Zhang Miao, modify from module_state_test.go
2014/7/9, by Li Bingyi, add unit test for SetNum() and GetNumState()
*/
/*
DESCRIPTION
*/
package module_state2

import (
    "fmt"
    "testing"
)

func TestModuleState(t *testing.T) {
    var state State
    var ok bool
    var value int64
    var vStr string
    var num  int64

    state.Init()        
    state.Inc("counter", 1)
    state.Inc("counter", 2)
    state.Dec("counter", 1)
    state.Set("state", "OK")
    state.SetNum("cap", 100)
    
    // test GetAll()
    data := state.GetAll()
    fmt.Println(*data)
    
    value, ok = data.SCounters["counter"]
    if !ok || value != 2 {
        t.Error("err in GetAll(), value should be 2")
    }
    
    vStr, ok = data.States["state"]
    if !ok || vStr != "OK" {
        t.Error("err in GetAll(), value should be OK")
    }
    
    // test GetCounter()
    value = state.GetCounter("counter")
    if value != 2 {
        t.Error("err in GetCounter(), value should be 2")
    }
    
    // test GetCounters()
    state.Inc("counter2", 3)
    counters := state.GetCounters()
    value, ok = counters["counter"]
    if !ok || value != 2 {
        t.Error("err in GetCounters(), value should be 2")
    }
    value, ok = counters["counter2"]
    if !ok || value != 3 {
        t.Error("err in GetCounters(), value should be 3")
    }

    // test GetState()
    vStr = state.GetState("state")
    if vStr != "OK" {
        t.Error("err in GetState(), value should be OK")
    }

    // test GetNumState()
    num = state.GetNumState("cap")
    if num != 100 {
        t.Error("err in GetNumSate(), num should be 100")
    }
}

func TestModuleStateCountersInit(t *testing.T) {
    var state State

    state.Init()
    
    // init counters
    keys := []string{"test1", "test2", "test3"}
    state.CountersInit(keys)

    // check counters
    counters := state.GetCounters()
    for _, key := range keys {
        value, ok := counters[key]
        
        if !ok || value != 0 {
            t.Error("err in CountersInit(), value should be 0")
        }
    }    
}

func TestModuleStateNil(t *testing.T) {
    // test support of nil for Inc(), Dec(), Set(), SetNum()
    var pState *State
    
    if pState != nil {
        t.Error("pState should be nil")
    }
    
    pState.Inc("test", 1)
    pState.Dec("test", 1)
    pState.Set("state", "ok")
    pState.SetNum("num", 1)
}

// test for StateData.NoahString()
func TestStateData_NoahString(t *testing.T) {
    sd := NewStateData()
    
    sd.SCounters.inc("counter", 1)
    sd.States["state"] = "ok"
    sd.NumStates["num_state"] = 1
    
    strOK := "counter:1\n" + "state:ok\n" + "num_state:1\n"
    
    if string(sd.NoahString()) != strOK {
        t.Error("err in StateData.NoahString()")
    }
}
