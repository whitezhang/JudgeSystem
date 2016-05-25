/* interval_test.go - test for interval.go  */
/*
modification history
--------------------
2014/9/11, by Zhang Miao, create
*/
/*
DESCRIPTION
*/
package module_state2

import (
    "fmt"
    "testing"
    "time"
)

func TestNextInterval(t *testing.T) {
    // test case 1 
    now := time.Date(2009, time.November, 10, 23, 10, 10, 0, time.UTC)
    interval := NextInterval(now, 60)
    if interval != 50 {
        t.Error(fmt.Sprintf("return of NextInterval() should be 50, it's %d", interval))
    }

    // test case 2
    now = time.Date(2009, time.November, 10, 23, 10, 0, 0, time.UTC)
    interval = NextInterval(now, 60)
    if interval != 60 {
        t.Error(fmt.Sprintf("return of NextInterval() should be 50, it's %d", interval))
    }

    // test case 3
    now = time.Date(2009, time.November, 10, 23, 10, 40, 0, time.UTC)
    interval = NextInterval(now, 60)    
    if interval != 20 {
        t.Error(fmt.Sprintf("return of NextInterval() should be 20, it's %d", interval))
    }
}

