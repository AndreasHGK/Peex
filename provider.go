package peex

import (
	"github.com/google/uuid"
	"reflect"
)

// ComponentProvider represent a struct that can load & save data associated to a player for a certain component.
type ComponentProvider[c Component] interface {
	// Load loads & writes data stored under the provider UUID to the pointer to the component c. Comp can be nil if it
	// is being loaded from an offline player.
	Load(id uuid.UUID, comp *c) error
	// Save writes the component to storage, using the UUID to identify the owner of the data.
	Save(id uuid.UUID, comp *c) error
}

// ProviderWrapper is a wrapper around a ComponentProvider to ensure strict typing of components and to easily allow for Peex
// to resolve the component type.
type ProviderWrapper[c Component] struct {
	p ComponentProvider[c]
}

// WrapProvider creates a new wrapper around a provider of the desired type.
func WrapProvider[c Component](p ComponentProvider[c]) ProviderWrapper[c] {
	if p == nil {
		panic("cannot provide nil as a provider")
	}
	return ProviderWrapper[c]{
		p: p,
	}
}

/// Internal provider logic
/// -----------------------

// GenericProvider is the interface representation of any type of ProviderWrapper, allowing them to be passed in the Config.
type GenericProvider interface {
	load(id uuid.UUID, x any) error
	loadNew(id uuid.UUID) (any, error)
	save(id uuid.UUID, x any) error
	// componentId returns the type of the component that the provider provides.
	componentId(m *Manager) componentId
	componentName() string
}

func (p ProviderWrapper[c]) load(id uuid.UUID, x any) error {
	return p.p.Load(id, x.(*c))
}

func (p ProviderWrapper[c]) loadNew(id uuid.UUID) (any, error) {
	v := new(c)
	err := p.p.Load(id, v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (p ProviderWrapper[c]) save(id uuid.UUID, x any) error {
	return p.p.Save(id, x.(*c))
}

func (p ProviderWrapper[c]) componentId(m *Manager) componentId {
	return m.getComponentId(new(c))
}

func (p ProviderWrapper[c]) componentName() string {
	t := reflect.TypeOf(new(c))
	return t.Name()
}
