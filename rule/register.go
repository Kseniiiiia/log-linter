package rule

import (
	"fmt"
	"sync"

	"log-linter/config"
)

type Factory func(cfg config.RuleConfig) (Rule, error)

type Registry struct {
	mu        sync.RWMutex
	factories map[string]Factory
}

func NewRegistry() *Registry {
	return &Registry{
		factories: make(map[string]Factory),
	}
}

func (r *Registry) Register(name string, factory Factory) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.factories[name]; ok {
		return fmt.Errorf("rule %q already registered", name)
	}

	r.factories[name] = factory
	return nil
}

func (r *Registry) Create(name string, cfg config.RuleConfig) (Rule, error) {
	r.mu.RLock()
	factory, ok := r.factories[name]
	r.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("rule %q not found", name)
	}

	return factory(cfg)
}

func (r *Registry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.factories))
	for name := range r.factories {
		names = append(names, name)
	}

	return names
}

var Global = NewRegistry()
