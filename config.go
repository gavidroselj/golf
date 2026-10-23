package golf

import (
	"strings"
	"sync"
)

type config struct {
	BasePath                string `json:"base_path"`
	DynamicValuePlaceholder string `json:"dynamic_value_placeholder"`
	IgnoreSubtestNames      bool   `json:"ignore_subtest_names"`

	inhibitUpdate bool
}

var (
	defaults = config{
		BasePath:                "testdata",
		DynamicValuePlaceholder: "-- DYNAMIC VALUE --",
		IgnoreSubtestNames:      false,

		inhibitUpdate: false,
	}

	mu      sync.RWMutex
	configs = make(map[string]config)
)

type Option func(*config)

func SetDefaults(opts ...Option) {
	mu.Lock()
	defer mu.Unlock()

	for _, opt := range opts {
		opt(&defaults)
	}
}

func Configure(t TestingT, opts ...Option) {
	mu.RLock()
	config := defaults
	mu.RUnlock()

	for _, opt := range opts {
		opt(&config)
	}

	mu.Lock()
	configs[TopLevelTestName(t)] = config
	mu.Unlock()

	t.Cleanup(func() {
		mu.Lock()
		defer mu.Unlock()

		delete(configs, TopLevelTestName(t))
	})
}

func TopLevelTestName(t TestingT) string {
	segments := strings.Split(t.Name(), "/")
	if len(segments) <= 0 {
		return t.Name()
	}

	return segments[0]
}

func getConfig(t TestingT) config {
	mu.RLock()
	defer mu.RUnlock()

	if c, ok := configs[TopLevelTestName(t)]; ok {
		return c
	}

	return defaults
}

func WithBasePath(basePath string) Option {
	return func(c *config) {
		c.BasePath = basePath
	}
}

func WithDynamicValuePlaceholder(dynamicValuePlaceholder string) Option {
	return func(c *config) {
		c.DynamicValuePlaceholder = dynamicValuePlaceholder
	}
}

func WithIgnoreSubtestNames(ignoreSubtestNames bool) Option {
	return func(c *config) {
		c.IgnoreSubtestNames = ignoreSubtestNames
	}
}

func withInhibitUpdate(inhibitUpdate bool) Option {
	return func(c *config) {
		c.inhibitUpdate = inhibitUpdate
	}
}
