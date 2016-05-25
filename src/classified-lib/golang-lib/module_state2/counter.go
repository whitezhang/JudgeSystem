/* counter.go - counters  */
/*
modification history
--------------------
2014/4/24, by Zhang Miao, create
*/
/*
DESCRIPTION
*/

package module_state2

/* flat counters    */
type Counters map[string]int64

// create new Counters
func NewCounters() Counters {
    counters := make(Counters)
    return counters
}

// increase value to key
func (c *Counters) inc(key string, value int) {
    _, ok := (*c)[key]
    
    if !ok {
        (*c)[key] = int64(value)
    } else {
        (*c)[key] += int64(value)
    }    
}

// decrease value to key
func (c *Counters) dec(key string, value int) {    
    _, ok := (*c)[key]
    
    if !ok {
        (*c)[key] = int64(value) * -1
    } else {
        (*c)[key] -= int64(value)
    }
}

// init counter for given keys to zero
// invoke init() to show all counters after start
func (c *Counters) init(keys []string) {
    for _, key := range keys {
        (*c)[key] = 0
    }
}

// make a copy of Counters
func (c *Counters) copy() Counters {
    copy := make(Counters)
    for key, value := range *c {
        copy[key] = value
    }
    return copy
}

// get change between two counters
func (c *Counters) diff(last Counters) Counters {
    diff := make(Counters)
    
    for key, value := range *c {
        old_value, ok := last[key]
        if ok {
            // exist in last
            diff[key] = value - old_value
        } else {
            diff[key] = value
        }
    }
    return diff
}

// calc sum of two Counters
func (c *Counters) Sum(c2 Counters) {
    for key, value2 := range c2 {
        value, ok := (*c)[key]
        if ok {
            (*c)[key] = value + value2
        } else {
            (*c)[key] = value2
        }
    }
}
