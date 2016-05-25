/* counter_hier_test.go - unit test for counter_hier.go */
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

// init case
func TestInit_case0(t *testing.T) {
    c := NewCounters()
    c.inc("baidu.op.bfe", 1)
    c.inc("baidu.op.ps", 2)
    c.inc("baidu.inf", 3)
    mt, _ := newMultiTree(c)

    hc := newHierCounters()
    hc.init(mt)

    // for baidu.info
    if v1, ok1 := hc["baidu"].(hierCounters); ok1 {
        if v2, ok2 := v1["inf"].(int64); ok2 {
            if v2 != int64(3) {
                t.Errorf("TestToHierCounters(): hc[\"baidu\"][\"info\"] (%d) != 3",
                ((hc["baidu"].(hierCounters))["inf"]).(int64))
            }
        } else {
            t.Errorf("TestToHierCounters(): hc[\"baidu\"][\"info\"] is not an int64")
        }
    } else {
        t.Errorf("TestToHierCounters(): hc[\"baidu\"] is not an Counters")
    }

    // for baidu.op.bfe
    if v3, ok3 := hc["baidu"].(hierCounters); ok3 {
        if v4, ok4 := v3["op"].(hierCounters); ok4 {
            if v5, ok5 := v4["bfe"].(int64); ok5 {
                if v5 != int64(1) {
                    t.Errorf("TestToHierCounters(): hc[\"baidu\"][\"op\"][\"bfe\"] (%d) != 1",
                    (((hc["baidu"].(hierCounters))["op"]).(hierCounters))["bfe"].(int64))
                }
            } else {
                t.Errorf("TestToHierCounters(): hc[\"baidu\"][\"op\"][\"bfe\"] is not an int64")
            }
        } else {
            t.Errorf("TestToHierCounters(): hc[\"baidu\"][\"op\"] is not an Counters")
        }
    } else {
        t.Errorf("TestToHierCounters(): hc[\"baidu\"] is not an Counters")
    }

    // for baidu.op.ps
    if v3, ok3 := hc["baidu"].(hierCounters); ok3 {
        if v4, ok4 := v3["op"].(hierCounters); ok4 {
            if v5, ok5 := v4["ps"].(int64); ok5 {
                if v5 != int64(2) {
                    t.Errorf("TestToHierCounters(): hc[\"baidu\"][\"op\"][\"ps\"] (%d) != 1",
                    (((hc["baidu"].(hierCounters))["op"]).(hierCounters))["bfe"].(int64))
                }
            } else {
                t.Errorf("TestToHierCounters(): hc[\"baidu\"][\"op\"][\"ps\"] is not an int64")
            }
        } else {
            t.Errorf("TestToHierCounters(): hc[\"baidu\"][\"op\"] is not an Counters")
        }
    } else {
        t.Errorf("TestToHierCounters(): hc[\"baidu\"] is not an Counters")
    }
}

// normal cases
func TestToHierCounters_case0(t *testing.T) {
    c := NewCounters()
    c.inc("baidu", 1)
    _, err := toHierCounters(c)
    if err != nil {
        t.Errorf("TestToHierCounters(): %s", err.Error())
    }
}

// normal cases
func TestToHierCounters_case1(t *testing.T) {
    c := NewCounters()
    c.inc("baidu", 1)
    c.inc("baidu.op", 1)
    _, err := toHierCounters(c)
    if err == nil {
        t.Error("TestToHierCounters(): err must not be nil")
    }
}

// normal cases
func TestToHierCounters_case2(t *testing.T) {
    c := NewCounters()
    hc, err := toHierCounters(c)
    if err != nil {
        t.Errorf("TestToHierCounters(): %s", err.Error())
    }

    if len(hc) != 0 {
        t.Errorf("TestToHierCounters(): len(hc)[%d] != 0", len(hc)) 
    }
}
