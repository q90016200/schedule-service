package common

import "time"

// MillisecondTimestamp create 13 digit timestamp
func MillisecondTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}