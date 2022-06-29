package peex

import (
	"errors"
	"github.com/df-mc/atomic"
	"github.com/df-mc/dragonfly/server/player"
	"reflect"
	"sync"
)

//go:generate go run ./cmd/events/main.go -o events.go -p peex -m ./go.mod
var _ player.Handler = (*Session)(nil)

// Session is a unique object that stores a player's data and handles player events. Data is stored in components, which
// can be added and removed from the Session at any time.
type Session struct {
	p atomic.Value[*player.Player]
	m *Manager

	components   map[componentId]Component
	componentsMu sync.RWMutex
}

// Player returns the Player that owns the Session. Returns nil if the Session is owned by a player that is no longer
// online.
func (s *Session) Player() *player.Player {
	return s.p.Load()
}

// Query runs a query function on the session. This function must use the Query, With and Option types as input
// parameters. These will work like they do in a Handler. True is returned if the query actually ran, else false.
func (s *Session) Query(queryFunc any) bool {
	info := s.m.makeQueryFuncInfo(queryFunc)
	return s.query(queryFunc, info)
}

// InsertComponent adds the Component to the player, keeping all it's values. An error is returned if the Component was
// already present.
func (s *Session) InsertComponent(c Component) error {
	cId := s.m.getComponentId(c)

	s.componentsMu.Lock()
	defer s.componentsMu.Unlock()
	if _, ok := s.components[cId]; ok {
		return errors.New("session already has a component of this type")
	}

	s.components[cId] = c
	// todo: recalculate handlers here?
	return nil
}

// SetComponent updates the value of a Component currently present in the Session. This is done regardless of whether
// this Component was present before.
func (s *Session) SetComponent(c Component) {
	cId := s.m.getComponentId(c)
	s.componentsMu.Lock()
	s.components[cId] = c
	// todo: recalculate handlers here?
	s.componentsMu.Unlock()
}

// Component returns the Component in the Session of the same type as the argument if it was found.
func (s *Session) Component(c Component) (Component, bool) {
	cId, ok := s.m.componentIdTable[reflect.TypeOf(c)] // we don't need to create a component id here
	if !ok {
		return nil, false
	}

	s.componentsMu.RLock()
	comp, ok := s.components[cId]
	s.componentsMu.RUnlock()
	return comp, ok
}

// RemoveComponent tries to remove the component with the same type as the provided argument from the Session. The value
// of the component will also be returned, If the Session does not have the component, nothing happens, and nil is
// returned.
func (s *Session) RemoveComponent(c Component) Component {
	cId, ok := s.m.componentIdTable[reflect.TypeOf(c)]
	if !ok {
		return nil
	}

	s.componentsMu.Lock()
	defer s.componentsMu.Unlock()
	if _, ok := s.components[cId]; !ok {
		return nil
	}

	c = s.components[cId]
	delete(s.components, cId)
	// todo: recalculate handlers here?
	return c
}

/// Internal session logic
/// ----------------------

// query executes a query function on the session (if it has all the required components).
func (s *Session) query(queryFunc any, info queryFuncInfo) bool {
	val := reflect.ValueOf(queryFunc)
	var args []reflect.Value

	s.componentsMu.RLock()
	defer s.componentsMu.RUnlock()
	for _, param := range info.params {
		c, ok := s.components[param.cId]
		if !ok && !param.optional {
			return false
		} else if !ok && param.optional {
			continue
		}

		args = append(args, reflect.ValueOf(param.query.set(c)))
	}

	val.Call(args)
	return true
}

func (s *Session) doQuit() {
	s.componentsMu.Lock()
	defer s.componentsMu.RUnlock()

	var p *player.Player
	// A nil player means the session is offline
	if !s.p.CompareAndSwap(p, nil) {
		panic("session has already disconnected")
	}

	s.m.sessionMu.Lock()
	delete(s.m.sessions, p.UUID())
	s.m.sessionMu.Unlock()
}
