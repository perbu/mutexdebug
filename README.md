# Mutexdebug

This package provides a simple mutex implementation for debugging purposes.
The constructor takes a timeout parameter, which is the maximum time a mutex
can be locked before it is considered to be locked forever. This is useful
for debugging deadlocks.

It has a significant performance impact, so it should not be used in production.


## Usage

First you need to create an interface that will allow you to easily switch
between the real mutex and the debug mutex. This is done by creating a
`MutexInterface`:

```go
package mypackage
type MutexInterface interface {
    Lock()
    Unlock()
}
```

Now you can initialize the mutex in your code:

```go
package main

func useMutex() {
	var m MutexInterface
	switch debug {
	case true:
		m = mutexdebug.NewDebugMutex(time.Millisecond*40, true)
	default:
		m = &sync.Mutex{}
	}
	// now we use the mutex as usual.
	m.Lock()
	// do something
	m.Unlock()
	
	// you should be able to cast the mutex back to the debug mutex
	// then you can access the debug information, like the number of warnings issued.
	// see the tests for more information.
	
}

```
## Complete example
Now you use the interface instead of the real mutex in your code and you can switch between the two.

```go
package main

import (
	"github.com/perbu/mutexdebug"
	"time"
)

type MutexInterface interface {
	Lock()
	Unlock()
}

const debug = true

func main() {
	var m MutexInterface
	switch debug {
	case true:
		m = mutexdebug.NewDebugMutex(time.Millisecond * 40, true)
    default:
		m = &sync.Mutex{}
	}
	// now we use the mutex as usual.
	m.Lock()
	// do something
	m.Unlock()
    if debug {
        dm := m.(*mutexdebug.DebugMutex)
        fmt.Println("Number of mutexes held too long: ",dm.Warnings.Load())
    }
	
}

```