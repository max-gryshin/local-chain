package pkg

import "fmt"

func GoWithRecover(fn func(), handlePanic func(err error)) {
	go func() {
		defer func() {
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
}
