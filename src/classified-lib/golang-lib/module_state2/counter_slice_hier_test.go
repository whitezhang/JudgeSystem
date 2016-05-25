/* counter_slice_hier_test.go - unit test for counter_slice_hier.go */

/*
modification history
--------------------
2014/11/24, by Li Bingyi, create
*/

/*
DESCRIPTION
*/
package module_state2

import (
    "testing"
)

func TestToHierCounterDiff_case0(t *testing.T) {
    var cd CounterDiff
    cd.LastTime = "lastTime"
    cd.Duration = 20
    cd.Diff = NewCounters()
    cd.Diff.inc("baidu.op", 1)
    hcd, err := toHierCounterDiff(&cd)
    if err != nil {
        t.Errorf("TestToHierCounterDiff(): %s", err.Error())
    }

    if hcd.LastTime != cd.LastTime {
        t.Errorf("hcd.LastTime[%s] != cd.LastTime[%s]", hcd.LastTime, cd.LastTime)
    }

    if hcd.Duration != cd.Duration {
        t.Errorf("hcd.Duration[%d] != cd.Duration[%d]", hcd.Duration, cd.Duration)
    }

}
