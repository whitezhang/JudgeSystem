/* interval.go - calc interval for get counter slice    */
/*
modification history
--------------------
2014/9/11
*/
/*
DESCRIPTION

*/
package module_state2

import (
	"time"
)

// Calculate seconds left for next interval
//
// Params:
//      - interval
// Returns:
//      seconds left for next interval
//
func NextInterval(now time.Time, interval int) int {
	seconds := now.Second()

	return interval - seconds%interval
}
