package closer

import (
	"log"
	"os"
	"os/signal"
	"sync"
)

// globalCloser is the default instance used by package-level functions.
var globalCloser = New()

// Add registers one or more cleanup functions with the global closer.
// These functions will be executed when CloseAll is called, in the order they were added.
// Each function should return an error if the cleanup fails.
func Add(f ...func() error) {
	globalCloser.Add(f...)
}

// Wait blocks until all cleanup functions registered with the global closer have completed.
// Typically called after CloseAll to wait for cleanup to finish.
func Wait() {
	globalCloser.Wait()
}

// CloseAll executes all registered cleanup functions with the global closer.
// Functions are executed concurrently, and CloseAll is idempotent (will only run once).
func CloseAll() {
	globalCloser.CloseAll()
}

// Closer manages a collection of cleanup functions and provides
// thread-safe methods for graceful shutdown management.
type Closer struct {
	mu    sync.Mutex     // protects funcs slice
	once  sync.Once      // ensures CloseAll only runs once
	done  chan struct{}  // signals when closing is complete
	funcs []func() error // registered cleanup functions
}

// New creates a new Closer instance. If signals are provided,
// it starts a goroutine that will trigger CloseAll when any of
// those signals are received.
func New(sig ...os.Signal) *Closer {
	c := &Closer{done: make(chan struct{})}
	if len(sig) > 0 {
		go func() {
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, sig...)
			<-ch
			signal.Stop(ch)
			c.CloseAll()
		}()
	}
	return c
}

// Add registers one or more cleanup functions with this Closer instance.
// Functions are appended to the existing list of cleanup functions.
func (c *Closer) Add(f ...func() error) {
	c.mu.Lock()
	c.funcs = append(c.funcs, f...)
	c.mu.Unlock()
}

// Wait blocks until CloseAll has completed execution of all registered functions.
// Returns immediately if CloseAll hasn't been called yet.
func (c *Closer) Wait() {
	<-c.done
}

// CloseAll executes all registered cleanup functions concurrently.
// This method is idempotent - subsequent calls will have no effect.
// Errors from cleanup functions are logged but otherwise ignored.
func (c *Closer) CloseAll() {
	c.once.Do(func() {
		defer close(c.done)

		c.mu.Lock()
		funcs := c.funcs
		c.funcs = nil
		c.mu.Unlock()

		errs := make(chan error, len(funcs))
		for _, f := range funcs {
			go func(f func() error) {
				errs <- f()
			}(f)
		}

		for i := 0; i < cap(errs); i++ {
			if err := <-errs; err != nil {
				log.Println("error returned from Closer")
			}
		}
	})
}
