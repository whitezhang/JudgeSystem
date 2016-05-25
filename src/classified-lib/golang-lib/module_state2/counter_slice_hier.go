/* counter_slice_hier.go - hierarchical CounterDiff */

/*
modification history
--------------------
2014/11/14, by Li Bingyi, create
*/

/*
DESCRIPTION
    This program provides converting flat CounterDiff to hierarchical CounterDiff for json output

Usage:
    import "www.badu.com/golang-lib/module_state2"
    var cd module_state2
    data, err := GetCdHierJson(&cd)

*/
package module_state2

import (
    "encoding/json"
    "fmt"
)

/* hierarchical structure for counter diff */
type hierCounterDiff struct {
	LastTime        string      // time till
	Duration        int         // in second

    Diff            hierCounters
    NoahKeyPrefix   string      //  for noah key
}

// convert CounterDiff to hierCounterDiff
// Params:
//  - cd: flat counter diff
// Returns:
//  - *hierCounterDiff: hierarchical counter diff
//  - error: error msg
func toHierCounterDiff(cd *CounterDiff) (*hierCounterDiff, error) {
    var hcd hierCounterDiff
    var err error

    hcd.Diff, err = toHierCounters(cd.Diff)
    if err != nil {
        return nil, fmt.Errorf("toHierCounterDiff(): %s", err.Error())
    }

    hcd.LastTime = cd.LastTime
    hcd.Duration = cd.Duration
    hcd.NoahKeyPrefix = cd.NoahKeyPrefix

    return &hcd, nil
}

// get hierarchical counter diff of json format
//Params:
//  - cd: flat counter diff
//Returns:
//  - []byte: json formated byte
//  - error: error msg
func GetCdHierJson(cd *CounterDiff) ([]byte, error) {
    hierCounterDiff, err := toHierCounterDiff(cd)
    if err != nil {
        return nil, fmt.Errorf("GetCdHierJson(): %s", err.Error())
    }

    return json.Marshal(hierCounterDiff)
}
