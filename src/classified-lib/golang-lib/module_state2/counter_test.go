/* counter_test.go - test for counter.go  */
/*
modification history
--------------------
2014/4/24, by Zhang Miao, create
*/
/*
DESCRIPTION
*/
package module_state2

import (
    "testing"
)

func TestCounterIncDec(t *testing.T) {
    counters := NewCounters()
    counters.inc("test", 2)
    counters.dec("test", 1)

    copy := counters.copy()
    
    value, ok := copy["test"]
    if !ok || value != 1 {
        t.Error("Counters.inc() or Counters.dec() fail")
    }
}

func TestCounterInit(t *testing.T) {
    counters := NewCounters()
    
    keys := []string{"test1", "test2", "test3"}
    counters.init(keys)

    copy := counters.copy()

    for _, key := range keys {
        value, ok := copy[key]
        if !ok || value != 0 {
            t.Error("Counters.init() fail")
        }
    }
}

func TestCounterCopy(t *testing.T) {
    counters := NewCounters()
    counters["test"] = 123
    
    copy := counters.copy()
    
    value, ok := copy["test"]
    if !ok || value != 123 {
        t.Error("Counters.copy() fail")
    }    
}

func TestCounterDiff_case1(t *testing.T) {
    counters := NewCounters()
    counters["test"] = 223
    
    last := NewCounters()
    last["test"] = 123
    
    diff := counters.diff(last)

    value, ok := diff["test"]
    if !ok || value != 100 {
        t.Error("Counters.diff() fail")
    }    
}

func TestCounterDiff_case2(t *testing.T) {
    counters := NewCounters()
    counters["test"] = 123
    
    last := NewCounters()    
    
    diff := counters.diff(last)

    value, ok := diff["test"]
    if !ok || value != 123 {
        t.Error("Counters.diff() fail")
    }    
}

func TestCounterSum(t *testing.T) {
    counters1 := NewCounters()
    counters1["test1"] = 10
    counters1["test2"] = 20

    counters2 := NewCounters()
    counters2["test2"] = 20
    counters2["test3"] = 30

    counters1.Sum(counters2)
    if counters1["test1"] != 10 {
        t.Error("Counters.Sum() fail")
    }
    if counters1["test2"] != 40 {
        t.Error("Counters.Sum() fail")
    }
    if counters1["test3"] != 30 {
        t.Error("Counters.Sum() fail")
    }
}
