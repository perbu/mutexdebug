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

Then you can create a `Mutex` type that implements this interface:

```go