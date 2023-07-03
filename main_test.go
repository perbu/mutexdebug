package mutexdebug_test

import (
	"fmt"
	"github.com/perbu/mutexdebug"
	"sync"
	"testing"
	"time"
)

func TestDebugMutex_UnlockBeforeTimeout(t *testing.T) {
	mutex := mutexdebug.NewDebugMutex(20*time.Millisecond, true)
	mutex.Lock()

	time.Sleep(10 * time.Millisecond)

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Mutex panicked when it should not have")
		}
	}()
	mutex.Unlock()
	if mutex.Warnings.Load() != 0 {
		t.Errorf("Expected 0 warnings, got %d", mutex.Warnings.Load())
	}
}

// TestDebugMutex_Relock locks and relocks the mutex 10 times in separate goroutines.
// it should not deadlock and issue no warnings.
func TestDebugMutex_Relock(t *testing.T) {
	mutex := mutexdebug.NewDebugMutex(10*time.Millisecond, true)
	wg := sync.WaitGroup{}
	mutex.Lock()
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mutex.Lock()
			time.Sleep(time.Millisecond)
			mutex.Unlock()
		}()
	}
	start := time.Now()
	mutex.Unlock()
	wg.Wait()
	fmt.Println("unlock took:", time.Since(start))
	if mutex.Warnings.Load() != 0 {
		t.Errorf("Expected 0 warnings, got %d", mutex.Warnings.Load())
	}
}

func TestDebugMutex_AfterTimeout(t *testing.T) {
	mutex := mutexdebug.NewDebugMutex(1*time.Millisecond, true)
	mutex.Lock()
	time.Sleep(3 * time.Millisecond)
	mutex.Unlock()
	if mutex.Warnings.Load() != 1 {
		t.Errorf("Expected 1 warning, got %d", mutex.Warnings.Load())
	}
	fmt.Println("done")

}
