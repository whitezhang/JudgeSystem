/* module_state2.go - for collecting state info of a module  */
/*
modification history
--------------------
2014/4/24
*/
/*
DESCRIPTION
This is a update version of module_state

Usage:
    import "www.baidu.com/golang-lib/module_state2"

    var state module_state2.State

    state.Init()

    state.Inc("counter", 1)
    state.Set("state", "OK")
    state.SetNum("cap", 100)

    stateData := state.Get()
*/
package module_state2

import (
	"bytes"
	"fmt"
	"sync"
)

/* state, one-level for SCounters */
type StateData struct {
	SCounters     Counters          // for count up
	States        map[string]string // for store states
	NumStates     Counters          // for store num states
	NoahKeyPrefix string            // for noah key
}

// state with mutex protect
type State struct {
	lock          sync.Mutex
	data          StateData
	noahKeyPrefix string
}

//
func NewStateData() *StateData {
	sd := new(StateData)
	sd.SCounters = NewCounters()
	sd.States = make(map[string]string)
	sd.NumStates = NewCounters()

	return sd
}

// make a copy for StateData
func (sd *StateData) copy() *StateData {
	copy := new(StateData)

	copy.SCounters = sd.SCounters.copy()

	copy.States = make(map[string]string)
	for key, value := range sd.States {
		copy.States[key] = value
	}

	copy.NumStates = NewCounters()
	for numKey, numValue := range sd.NumStates {
		copy.NumStates[numKey] = numValue
	}

	return copy
}

func (sd *StateData) noahKeyGen(str string) string {
	if sd.NoahKeyPrefix == "" {
		return str
	}

	return fmt.Sprintf("%s_%s", sd.NoahKeyPrefix, str)
}

// output noah string (lines of key:value) for StateData
func (sd *StateData) NoahString() []byte {
	var buf bytes.Buffer

	// print SCounters
	for key, value := range sd.SCounters {
		key = sd.noahKeyGen(key)
		str := fmt.Sprintf("%s:%d\n", key, value)
		buf.WriteString(str)
	}

	// print States
	for key, value := range sd.States {
		key = sd.noahKeyGen(key)
		str := fmt.Sprintf("%s:%s\n", key, value)
		buf.WriteString(str)
	}

	// print NumStates
	for key, value := range sd.NumStates {
		key = sd.noahKeyGen(key)
		str := fmt.Sprintf("%s:%d\n", key, value)
		buf.WriteString(str)
	}

	return buf.Bytes()
}

/* Initialize the state */
func (s *State) Init() {
	s.data.SCounters = NewCounters()
	s.data.States = make(map[string]string)
	s.data.NumStates = NewCounters()
}

/* set noah key prefix */
func (s *State) SetNoahKeyPrefix(prefix string) {
	s.noahKeyPrefix = prefix
}

/* Increase value to key */
func (s *State) Inc(key string, value int) {
	// support s is nil
	if s == nil {
		return
	}

	s.lock.Lock()
	s.data.SCounters.inc(key, value)
	s.lock.Unlock()
}

/* Decrease value to key */
func (s *State) Dec(key string, value int) {
	// support s is nil
	if s == nil {
		return
	}

	s.lock.Lock()
	s.data.SCounters.dec(key, value)
	s.lock.Unlock()
}

/* Init counters for given keys to zero */
func (s *State) CountersInit(keys []string) {
	s.lock.Lock()
	s.data.SCounters.init(keys)
	s.lock.Unlock()
}

/* set state to key */
func (s *State) Set(key string, value string) {
	// support s is nil
	if s == nil {
		return
	}

	s.lock.Lock()
	s.data.States[key] = value
	s.lock.Unlock()
}

/* set num state to key */
func (s *State) SetNum(key string, value int64) {
	// support s is nil
	if s == nil {
		return
	}

	s.lock.Lock()
	s.data.NumStates[key] = value
	s.lock.Unlock()
}

/* Get counter value of given key    */
func (s *State) GetCounter(key string) int64 {
	s.lock.Lock()
	value, ok := s.data.SCounters[key]
	s.lock.Unlock()

	if !ok {
		value = 0
	}

	return value
}

/* Get all counters */
func (s *State) GetCounters() Counters {
	s.lock.Lock()
	counters := s.data.SCounters.copy()
	s.lock.Unlock()

	return counters
}

/* Get state value of given key    */
func (s *State) GetState(key string) string {
	s.lock.Lock()
	value, ok := s.data.States[key]
	s.lock.Unlock()

	if !ok {
		value = ""
	}

	return value
}

/* Get num state value of given key    */
func (s *State) GetNumState(key string) int64 {
	s.lock.Lock()
	value, ok := s.data.NumStates[key]
	s.lock.Unlock()

	if !ok {
		value = 0
	}

	return value
}

/* Get all states    */
func (s *State) GetAll() *StateData {
	s.lock.Lock()
	copy := s.data.copy()
	s.lock.Unlock()

	copy.NoahKeyPrefix = s.noahKeyPrefix
	return copy
}

/* Get noah prefix key */
func (s *State) GetNoahKeyPrefix() string {
	return s.noahKeyPrefix
}
