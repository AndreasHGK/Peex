package peex

import (
	"github.com/df-mc/dragonfly/server/player"
	"reflect"
)

// Component represents some data type that can be stored in a player session. There can only be one component of each
// type per player, but each player does not necessarily need to have every component.
type Component interface{}

// Adder represents a Component that has a function that is called when the component is added to a Session in any way.
type Adder interface {
	Component
	// Add gets called right after the component is added to a Session. This is always called regardless of how the
	// component is added. It is also called when another component of the same type was present but got replaced. The
	// owner of the session will be passed along as an argument. Gets called after the component has been loaded.
	Add(p *player.Player)
}

// Remover represents a Component that has extra logic that runs when the component is removed from a session in any
// way.
type Remover interface {
	Component
	// Remove gets called right before the current component instance gets removed from the Session. This means the
	// method is also called when the component gets replaced with another of the same type. The owner of the session is
	// passed along as argument. Gets called before the component is saved.
	Remove(p *player.Player)
}

// ComponentFromSession returns and automatically type casts a user's component to the correct type if it is present.
func ComponentFromSession[T Component](s *Session) (T, bool) {
	comp, ok := s.Component(new(T))
	// A nil interface cannot be converted(?)
	if !ok {
		return *new(T), false
	}
	return comp.(T), true
}

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
