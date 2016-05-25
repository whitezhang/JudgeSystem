/* module_state_hier.go - hierarchical StateData */

/*
modification history
--------------------
2014/11/14, by Li Bingyi, Create
*/

/*
DESCRIPTION
    This program provides converting flat StateDate to hierarchical StateData for json output

Usage:
    import "www.badu.com/golang-lib/module_state2"
    var sd module_state2
    data, err := GetSdHierJson(&sd)
*/

package module_state2

import (
    "encoding/json"
    "fmt"
)

// hierarchical structure for StateData 
type hierStateData struct {
    SCounters       hierCounters        // for count up
    States          map[string]string   // for store states
    NumStates       hierCounters        // for store num states
    NoahKeyPrefix   string
}

// convert StateData to hierStateData
// Params:
//  - sd: flat state data
// Returns:
//  - *hierStateData: hierarchical state data
//  - error: error msg

func toHierStateData(sd *StateData) (*hierStateData, error) {
    var hsd hierStateData
    var err error

    hsd.SCounters, err = toHierCounters(sd.SCounters)
    if err != nil {
        return nil, fmt.Errorf("toHierStateData(): Scounters %s", err.Error())
    }

    hsd.States = sd.States
    hsd.NumStates, err = toHierCounters(sd.NumStates)
    if err != nil {
        return nil, fmt.Errorf("toHierStateData(): NumStates %s", err.Error())
    }

    hsd.NoahKeyPrefix = sd.NoahKeyPrefix

    return &hsd, nil
}

// get hierarchical StataData of json format
//Params:
//  - sd: flat state data
//Returns:
//  - []byte: json formated byte
//  - error: error msg
func GetSdHierJson(sd *StateData) ([]byte, error) {
    hierState, err := toHierStateData(sd)
    if err != nil {
        return nil, fmt.Errorf("GetSdHierJson(): %s", err.Error())
    }

    return json.Marshal(hierState)
}
