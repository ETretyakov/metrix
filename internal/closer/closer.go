// Module "closer" implements a gracefull closer shutdown pattern for the main service.
package closer

import (
	"context"

	"metrix/pkg/logger"
)

var globalCloser = New()

// Add - a function that adds a closer function to the global closer.
func Add(f ...func() error) {
	globalCloser.Add(f...)
}

// CloseAll - a function that initiates close functions in the closer.
func CloseAll() {
	globalCloser.CloseAll()
}

// Closer - the closer struct that stores all closer functions.
type Closer struct {
	funcs []func() error
}

// New - the function that initaites Closer structure.
func New() *Closer {
	return &Closer{
		funcs: make([]func() error, 0),
	}
}

// Add - the function that adds function to the specific closer.
func (c *Closer) Add(f ...func() error) {
	c.funcs = append(c.funcs, f...)
}

// CloseAll - the fucntion that intiates stored cose functions sequentialy.
func (c *Closer) CloseAll() {
	for _, f := range c.funcs {
		if err := f(); err != nil {
			logger.Error(context.TODO(), "error close", err)
		}
	}
}
