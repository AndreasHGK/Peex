package peex

import (
	"reflect"
)

// Component represents some data type that can be stored in a player session. There can only be one component of each
// type per player, but each player does not necessarily need to have every component.
type Component interface{}

/// Internal component logic
/// ------------------------

// componentId uniquely identifies a component type in a Session Manager.
type componentId uint

// makeComponentId returns a unique integer that identifies a component type. If this identifier cannot be found, a new
// one is created for the type
func (m *Manager) getComponentId(c Component) componentId {
	return m.getComponentIdRefl(reflect.TypeOf(c))
}

func (m *Manager) getComponentIdRefl(t reflect.Type) componentId {
	id, ok := m.componentIdTable[t]
	if !ok {
		m.componentIdTable[t], id = m.componentNextId, m.componentNextId
		m.componentNextId++
	}
	return id
}
