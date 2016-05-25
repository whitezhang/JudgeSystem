/* counter_hier.go - hier counters  */
/*
modification history
--------------------
2014/11/12, by Li Bingyi, create
*/
/*
DESCRIPTION
    This program provides converting flat Counters to hierarchical Counters
*/

package module_state2

import (
    "fmt"
)

/* hierarchical counters, just for json dump */
type hierCounters map[string]interface{}

// create new hierCounters
func newHierCounters() hierCounters {
    hCounters := make(hierCounters)
    return hCounters
}

// init hierarchical counters with multi tree
// Params:
//  - t: root node of multi tree
func (hc hierCounters) init(t *treeNode) {
    if t.children == nil {
        // t is leaf node
        hc[t.elem.key] = t.elem.value
    } else {
        for i := 0; i < len(t.children); i++ {
            child := t.children[i]
            if child.children == nil {
                // child is leaf node, value of the key is node value
                hc[child.elem.key] = child.elem.value
            } else {
                nhc := newHierCounters()
                // child is not leaf node, value of the key is a hierCounters
                hc[child.elem.key] = nhc
                // init nhc with child node
                nhc.init(child)
            }
        }
    }
}

// convert Counters(flat counters) to HierCounters(hierarchical counters)
// Params:
//  - c: flat Counters
// Returns: (hierCounters, error)
//  - hierCounters: hier counters if convert ok, else nil
//  - error: nil if convert ok, else err info
func toHierCounters(c Counters) (hierCounters, error) {

    // new multiTree with flat counters
    root, err := newMultiTree(c)
    if err != nil {
        return nil, fmt.Errorf("Counters.toHierCounters(): %s", err.Error())
    }

    // new hierCounters
    hCounters := newHierCounters()
    // init hCounters with root only when tree has child
    if root.children != nil {
        hCounters.init(root)
    }

    return hCounters, nil
}
