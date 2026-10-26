package golf

import (
	"strings"
	"sync"
)

type config struct {
	BasePath                string `json:"base_path"`
	DynamicValuePlaceholder string `json:"dynamic_value_placeholder"`
	IgnoreSubtestNames      bool   `json:"ignore_subtest_names"`
	InhibitUpdate           bool   `json:"inhibit_update"`
}

var (
	defaults = config{
		BasePath:                "testdata",
		DynamicValuePlaceholder: "-- DYNAMIC VALUE --",
		IgnoreSubtestNames:      false,
		InhibitUpdate:           false,
	}

	mu      sync.RWMutex
	configs = make(map[string]config)
)

type Option func(*config)

// SetDefaults configures the default options applied to all golf tests.
// These defaults can then be modified using Configure.
func SetDefaults(opts ...Option) {
	mu.Lock()
	defer mu.Unlock()

	for _, opt := range opts {
		opt(&defaults)
	}
}

// Configure changes test options for specific tests. Configure should be called at the start of a top level test.
// Configure will apply to the top level test and all subtests. The test hirearchy is determined via the top level test name.
// Having multiple top level tests with the same name is untested and may cause unintended behavior.
func Configure(t TestingT, opts ...Option) {
	mu.RLock()
	config := defaults
	mu.RUnlock()

	for _, opt := range opts {
		opt(&config)
	}

	mu.Lock()
	configs[topLevelTestName(t)] = config
	mu.Unlock()

	t.Cleanup(func() {
		mu.Lock()
		defer mu.Unlock()

		delete(configs, topLevelTestName(t))
	})
}

// topLevelTestName returns the name of the top level test. This is used for determining the test hiearchy.
func topLevelTestName(t TestingT) string {
	segments := strings.Split(t.Name(), "/")
	if len(segments) <= 0 {
		return t.Name()
	}

	return segments[0]
}

// getConfig gets the config for the current test.
func getConfig(t TestingT) config {
	mu.RLock()
	defer mu.RUnlock()

	if c, ok := configs[topLevelTestName(t)]; ok {
		return c
	}

	return defaults
}

// WithBasePath configures the base path for golf test artifacts to be stored.
// The default value is "testdata", the path is relative from the folder the test file is in.
func WithBasePath(basePath string) Option {
	return func(c *config) {
		c.BasePath = basePath
	}
}

// WithDynamicValuePlaceholder configures the value with which dynamic values will be replaced in the stored golden files.
// The default value is "-- DYNAMIC VALUE --"
func WithDynamicValuePlaceholder(dynamicValuePlaceholder string) Option {
	return func(c *config) {
		c.DynamicValuePlaceholder = dynamicValuePlaceholder
	}
}

// WithIgnoreSubtestNames configures whether golden files should belong to the top level test, or be individual for every subtest.
// This is useful when you have multiple subtests which are expecting the same data via different means.
// The default value is "false".
func WithIgnoreSubtestNames(ignoreSubtestNames bool) Option {
	return func(c *config) {
		c.IgnoreSubtestNames = ignoreSubtestNames
	}
}

// WithInhibitUpdate configures whether the golden files should ignore the `-update` flag.
// The default value is "false".
func WithInhibitUpdate(inhibitUpdate bool) Option {
	return func(c *config) {
		c.InhibitUpdate = inhibitUpdate
	}
}
