package peex

import (
	"errors"
	"fmt"
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

// SaveAll saves every component that can be saved. If for any component an error is returned, the last error will be
// returned by this function.
func (s *Session) SaveAll() error {
	var e error
	uuid := s.Player().UUID()

	s.componentsMu.RLock()
	for id, p := range s.m.componentProvs {
		c, ok := s.components[id]
		if !ok {
			continue
		}
		err := p.save(uuid, c)
		// If there was an error saving the component, save it, so it can be returned. Will overwrite previous errors.
		// Do not automatically return on error, as we want to minimize any data loss.
		if err != nil {
			e = err
		}
	}
	s.componentsMu.RUnlock()

	if e != nil {
		return fmt.Errorf("error while saving component: %w", e)
	}
	return nil
}

// Save saves a single component type for the session. Saves the component of the same type as the argument that is
// currently present as opposed to the one provided as argument.
func (s *Session) Save(c Component) error {
	s.componentsMu.RLock()
	defer s.componentsMu.RUnlock()

	cId := s.m.getComponentId(c)
	c, ok := s.components[cId]
	if !ok {
		return errors.New("trying to save a component without a provider") // todo: should this return nil?
	}
	p, ok := s.m.componentProvs[cId]
	if !ok {
		return errors.New("trying to save a component without a provider")
	}

	err := p.save(s.Player().UUID(), c)
	if err != nil {
		return fmt.Errorf("error while saving component: %w", err)
	}
	return nil
}

// InsertComponent adds the Component to the player, keeping all it's values. An error is returned if the Component was
// already present. Also loads the component if a provider for it has been set in the config.
func (s *Session) InsertComponent(c Component) error {
	cId := s.m.getComponentId(c)

	s.componentsMu.Lock()
	defer s.componentsMu.Unlock()
	return s.insertComponent(cId, c)
	// todo: recalculate handlers here?
}

// SetComponent updates the value of a Component currently present in the Session. This is done regardless of whether
// this Component was present before. If the component was already present, and it implements Remover, the Remove method
// will first be called on the previous instance of the component. The Add method will be called on the new Component if
// it implements Adder.
//
// NOTE: does NOT load the component!
func (s *Session) SetComponent(c Component) {
	cId := s.m.getComponentId(c)
	s.componentsMu.Lock()

	p := s.Player()
	// If the component is already present, first call Remove() on the previous component if it implements it.
	if prev, ok := s.components[cId]; ok {
		if r, ok := prev.(Remover); ok {
			r.Remove(p)
		}
	}

	s.components[cId] = c
	if a, ok := c.(Adder); ok {
		a.Add(p)
	}
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
// returned. Also saves the component if a provider for it has been set in the config.
func (s *Session) RemoveComponent(c Component) (Component, error) {
	cId, ok := s.m.componentIdTable[reflect.TypeOf(c)]
	if !ok {
		return nil, errors.New("trying to remove unknown component")
	}

	s.componentsMu.Lock()
	defer s.componentsMu.Unlock()
	return s.removeComponent(cId, c)
}

/// Internal session logic
/// ----------------------

// query executes a query function on the session (if it has all the required components).
func (s *Session) query(queryFunc any, info queryFuncInfo) bool {
	val := reflect.ValueOf(queryFunc)
	args := make([]reflect.Value, 0, len(info.params))

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

// insertComponent adds a component to the session. This method is not safe for use in multiple goroutines.
func (s *Session) insertComponent(cId componentId, c Component) error {
	if _, ok := s.components[cId]; ok {
		return errors.New("session already has a component of this type")
	}

	// Try to load the component if it has a provider.
	if p, ok := s.m.componentProvs[cId]; ok {
		err := p.load(s.Player().UUID(), c)
		if err != nil {
			return fmt.Errorf("error while loading component: %w", err)
		}
	}
	s.components[cId] = c
	if a, ok := c.(Adder); ok {
		a.Add(s.Player())
	}
	return nil
}

// removeComponent removes a component from the session. This method is not safe for use in multiple goroutines.
func (s *Session) removeComponent(cId componentId, c Component) (Component, error) {
	if _, ok := s.components[cId]; !ok {
		return nil, errors.New("trying to remove a component not present in the session")
	}

	c = s.components[cId]
	if r, ok := c.(Remover); ok {
		r.Remove(s.Player())
	}
	// Try to save the component
	if p, ok := s.m.componentProvs[cId]; ok {
		err := p.save(s.Player().UUID(), c)
		if err != nil {
			return nil, fmt.Errorf("error while saving component: %w", err)
		}
	}
	delete(s.components, cId)
	// todo: recalculate handlers here?
	return c, nil
}

func (s *Session) doQuit() {
	s.componentsMu.Lock()
	defer s.componentsMu.Unlock()

	for _, comp := range s.components {
		_, err := s.removeComponent(s.m.getComponentId(comp), comp)
		if err != nil && s.m.logger != nil {
			s.m.logger.Errorf("%w", err)
		}
	}

	var p *player.Player
	// A nil player means the session is offline
	if p = s.p.Swap(nil); p == nil {
		panic("session owner has already disconnected")
	}

	s.m.sessionMu.Lock()
	delete(s.m.sessions, p.UUID())
	s.m.sessionMu.Unlock()
	s.components = nil
}
