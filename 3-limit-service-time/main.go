//////////////////////////////////////////////////////////////////////
//
// Your video processing service has a freemium model. Everyone has 10
// sec of free processing time on your service. After that, the
// service will kill your process, unless you are a paid premium user.
//
// Beginner Level: 10s max per request
// Advanced Level: 10s max per user (accumulated)
//

package main

import (
	"context"
	"sync"
	"time"
)

// User defines the UserModel. Use this to check whether a User is a
// Premium user or not
type User struct {
	ID        int
	IsPremium bool
	TimeUsed  int64 // in seconds
	sync.Mutex
}

// HandleRequest runs the processes requested by users. Returns false
// if process had to be killed
func HandleRequest(process func(), u *User) bool {
	if !u.IsPremium && u.TimeUsed > 10 {
		return false
	}

	ctx := context.Background()
	if u.IsPremium {
		process()
		return true
	}

	u.Lock()
	defer u.Unlock()

	ctx, cancel := context.WithTimeout(ctx, time.Duration(10 - u.TimeUsed) * time.Second)
	defer cancel()

	start := time.Now()
	process()
	used := time.Now().Sub(start)
	u.TimeUsed += int64(used.Seconds())
	select {
	case <-ctx.Done():
		return false
	default:
		return true
	}
}

func main() {
	RunMockServer()
}
