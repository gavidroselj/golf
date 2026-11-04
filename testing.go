package golf

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestingT interface {
	assert.TestingT
	Name() string
	Cleanup(func())
}

type T struct {
	assert.TestingT
	name      string
	collector *Collector
}

func NewT(t *testing.T) *T {
	return &T{TestingT: t, name: t.Name()}
}

func (t *T) Name() string {
	return t.name
}

func (t *T) Cleanup(func()) {
}

func (t *T) Collect() *T {
	collector := &Collector{}
	return &T{
		TestingT:  collector,
		name:      t.name,
		collector: collector,
	}
}

func (t *T) Errors() []error {
	if t.collector == nil {
		return nil
	}

	return t.collector.errors
}

type Collector struct {
	errors []error
}

func (c *Collector) Errorf(format string, args ...any) {
	c.errors = append(c.errors, fmt.Errorf(format, args...))
}

func (c *Collector) Errors() []error {
	return c.errors
}

func ExpectFailure(t *testing.T, f func(t TestingT)) {
	ct := NewT(t).Collect()

	f(ct)

	assert.GreaterOrEqual(t, len(ct.Errors()), 1, "Expected test to fail, but it passed")
}
