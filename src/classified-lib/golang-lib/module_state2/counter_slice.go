/* counter_slice.go - get diff of two counters  */
/*
modification history
--------------------
2014/11/12
    - move codes from waf-server for periodically get counter slice
*/
/*
DESCRIPTION

Usage:
    import "www.baidu.com/golang-lib/module_state2"

    var counter module_state2.Counter
    var counterSlice *module_state2.CounterSlice
    var state *module_state2.State

    // usage 1: get diff once
    counterSlice.Set(counter)
    // make some update to counter here
    counterSlice.Set(counter)
    // get diff between update
    diff := counterSlice.Get()

    // usage 2: update diff periodically and get when needed
    var examCnt examCounter
    // update diff periodically
    counterSlice.Init(state, interval)
    // get diff
    diff := counterSlice.Get()
*/

package module_state2

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

/* diff of two counters    */
type CounterSlice struct {
	lock sync.Mutex

	lastTime time.Time
	duration time.Duration

	countersLast Counters //  last absolute counter
	countersDiff Counters //  diff in last duration

	noahKeyPrefix string //  for noah key
}

type CounterDiff struct {
	LastTime string // time till
	Duration int    // in second

	Diff Counters

	NoahKeyPrefix string // for noah key
}

/* set for noah key prefix */
func (cs *CounterSlice) SetNoahKeyPrefix(prefix string) {
	cs.noahKeyPrefix = prefix
}

/* get noah key prefix */
func (cs *CounterSlice) GetNoahKeyPrefix() string {
	return cs.noahKeyPrefix
}

/* set to counter slice */
func (cs *CounterSlice) Set(counters Counters) {
	cs.lock.Lock()
	defer cs.lock.Unlock()

	if cs.countersLast == nil {
		// not initialized
		cs.lastTime = time.Now()
		cs.countersLast = counters.copy()
		cs.countersDiff = NewCounters()
	} else {
		now := time.Now()
		cs.duration = now.Sub(cs.lastTime)
		cs.lastTime = now

		cs.countersDiff = counters.diff(cs.countersLast)
		cs.countersLast = counters.copy()
	}
}

/* get diff from counter slice   */
func (cs *CounterSlice) Get() CounterDiff {
	var retVal CounterDiff

	cs.lock.Lock()
	defer cs.lock.Unlock()

	if cs.countersLast == nil {
		retVal.Diff = NewCounters()
	} else {
		retVal.LastTime = cs.lastTime.Format("2006-01-02 15:04:05")
		retVal.Duration = int(cs.duration.Seconds())
		retVal.Diff = cs.countersDiff.copy()
	}

	retVal.NoahKeyPrefix = cs.noahKeyPrefix

	return retVal
}

// get json format of counter diff
func (cs *CounterSlice) GetJson() ([]byte, error) {
	return json.Marshal(cs.Get())
}

func (cd CounterDiff) noahKeyGen(str string) string {
	if cd.NoahKeyPrefix == "" {
		return str
	}

	return fmt.Sprintf("%s_%s", cd.NoahKeyPrefix, str)
}

// output noah string (lines of key:value) for CounterDiff
func (cd CounterDiff) NoahString() []byte {
	var buf bytes.Buffer

	// LastTime
	str := cd.noahKeyGen("LastTime")
	str = fmt.Sprintf("%s:%s\n", str, cd.LastTime)
	buf.WriteString(str)

	// Duration
	str = cd.noahKeyGen("Duration")
	str = fmt.Sprintf("%s:%d\n", str, cd.Duration)
	buf.WriteString(str)

	// print Diff
	for key, value := range cd.Diff {
		key = cd.noahKeyGen(key)
		str := fmt.Sprintf("%s:%d\n", key, value)
		buf.WriteString(str)
	}

	return buf.Bytes()
}

// go-routine for periodically get counter slice
func (cs *CounterSlice) handleCounterSlice(s *State, interval int) {
	for {
		counter := s.GetCounters()
		cs.Set(counter)

		leftSeconds := NextInterval(time.Now(), interval)
		time.Sleep(time.Duration(leftSeconds) * time.Second)
	}
}

// init the counter diff
// Params:
//    - s: module State
//    - interval: interval to compute between two counters
// Notice: use this method only when you need to get diff between two counters periodically
func (cs *CounterSlice) Init(s *State, interval int) {
	go cs.handleCounterSlice(s, interval)
}
