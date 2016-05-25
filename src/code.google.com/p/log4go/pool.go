/* pool.go - pool for binary data */
/*
modification history
--------------------
2015/01/27, by taochunhua, create
*/
/*
DESCRIPTION
*/

package log4go

import (
    "errors"
    "sync"
)

// buffer pool for binary data
var buf4kPool   sync.Pool
var buf16kPool  sync.Pool

var (
    ErrTooLarge = errors.New("required slice size exceed 16k")
)

// get proper []byte from pool
// if size > 16K, return ErrTooLarge
func NewBuffer(size int) ([]byte, error) {
    var pool sync.Pool

    // return buffer size
    originSize := size

    if size <= 4096 {
       size = 4096
        pool = buf4kPool
    } else if size <= 16 * 1024 {
       size = 16 * 1024
        pool = buf16kPool
    } else {
        // if message is larger than 16K, return err
        return nil, ErrTooLarge
    }

    if v := pool.Get(); v != nil {
        return v.([]byte)[:originSize], nil
    }

    return make([]byte, size)[:originSize], nil
}

func putBuffer(b []byte) {
    b = b[:cap(b)]
    if cap(b) == 4096 {
        buf4kPool.Put(b)
    }
    if cap(b) == 16 * 1024 {
        buf16kPool.Put(b)
    }
}
