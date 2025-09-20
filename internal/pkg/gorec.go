package pkg

import "fmt"

// GoWithRecoverAndSemaphore is a function that executes the provided function in a goroutine
// sem is a semaphore channel that limits the number of concurrent executions
func GoWithRecoverAndSemaphore(fn func(), handlePanic func(err error), sem chan struct{}) {
	select {
	case sem <- struct{}{}:
		go func() {
			defer func() {
				<-sem
				if r := recover(); r != nil {
					if err, ok := r.(error); ok {
						handlePanic(err)
					} else {
						handlePanic(fmt.Errorf("panic occurred: %v", r))
					}
				}
			}()
			fn()
		}()
	default:
		return
	}
}
