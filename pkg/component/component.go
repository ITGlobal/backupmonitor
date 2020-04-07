package component

import (
	"fmt"
	"sync"

	"github.com/sarulabs/di"
)

// Factory creates an instance of component from a DI container
type Factory func(c di.Container) (T, error)

// Builder is a component and service registry builder
type Builder interface {
	// Add new service
	AddService(def di.Def)

	// Add new component factory
	AddComponent(factory Factory)

	// Build a component and service registry
	Build() (Registry, error)
}

// Registry is a component and service registry
type Registry interface {
	// Get registered services
	Services() di.Container

	// Get registered components
	Components() []T
}

// T is a runnable component that supports graceful shutdown
type T interface {
	// Start a component
	Start(group *sync.WaitGroup, stop chan interface{})
}

// NewBuilder creates a new instance of Builder
func NewBuilder() Builder {
	services, _ := di.NewBuilder()

	builder := &builderImpl{
		services:   services,
		keys:       make([]string, 0),
		components: make([]Factory, 0),
		errors:     make([]error, 0),
	}

	return builder
}

type builderImpl struct {
	services   *di.Builder
	keys       []string
	components []Factory
	errors     []error
}

// Add new service
func (b *builderImpl) AddService(def di.Def) {
	err := b.services.Add(def)
	if err != nil {
		b.errors = append(b.errors, err)
	}

	b.keys = append(b.keys, def.Name)
}

// Add new component factory
func (b *builderImpl) AddComponent(factory Factory) {
	b.components = append(b.components, factory)
}

// Build a component and service registry
func (b *builderImpl) Build() (Registry, error) {
	if len(b.errors) > 0 {
		return nil, b.errors[0]
	}

	services := b.services.Build()
	if services == nil {
		return nil, fmt.Errorf("Unable to build DI container")
	}

	for _, key := range b.keys {
		_, err := services.SafeGet(key)
		if err != nil {
			return nil, fmt.Errorf("Unable to build service \"%s\": %v", key, err)
		}
	}

	components := make([]T, len(b.components))
	for i, factory := range b.components {
		component, err := factory(services)
		if err != nil {
			return nil, err
		}

		components[i] = component
	}

	return &registryImpl{services, components}, nil
}

type registryImpl struct {
	services   di.Container
	components []T
}

// Get registered services
func (r *registryImpl) Services() di.Container {
	return r.services
}

// Get registered components
func (r *registryImpl) Components() []T {
	return r.components
}
