package peex

import (
	"errors"
	"fmt"
	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/google/uuid"
	"reflect"
	"strings"
	"sync"
)

// Manager stores all current sessions. It also contains all the registered handlers and component types.
type Manager struct {
	logger server.Logger

	sessions  map[uuid.UUID]*Session
	sessionMu sync.RWMutex

	handlerNextId  handlerId
	handlerIdTable map[reflect.Type]handlerId
	handlers       map[handlerId]handlerInfo
	eventHandlers  map[eventId][]handlerId

	componentNextId  componentId
	componentIdTable map[reflect.Type]componentId
	componentProvs   map[componentId]GenericProvider
	// todo: component cache
}

// New creates a new Session Manager. It also inserts all the provided handlers into the manager. Events will be called
// in the order the handlers are added to the manager. These handlers can query for specific components to be present in
// a player Session in order to actually run.
func New(cfg Config) *Manager {
	m := &Manager{
		logger:           cfg.Logger,
		sessions:         map[uuid.UUID]*Session{},
		handlerIdTable:   map[reflect.Type]handlerId{},
		handlers:         map[handlerId]handlerInfo{},
		eventHandlers:    map[eventId][]handlerId{},
		componentIdTable: map[reflect.Type]componentId{},
		componentProvs:   map[componentId]GenericProvider{},
	}
	for _, id := range allEvents {
		m.eventHandlers[id] = []handlerId{}
	}

	for _, h := range cfg.Handlers {
		t := reflect.TypeOf(h)
		if _, ok := m.handlerIdTable[t]; ok {
			panic("re-registering an existing handler type")
		}

		// Assign the handler ID to the type, and generate the handlerInfo
		m.handlerIdTable[t] = m.handlerNextId
		info := m.createHandlerInfo(h)
		m.handlers[m.handlerNextId] = info
		for id, _ := range info.events {
			m.eventHandlers[id] = append(m.eventHandlers[id], m.handlerNextId)
		}

		// Check if the handler (partially) implements a possibly outdated version of player.Handler, preventing events
		// from being silently ignored.
		for eventName, id := range allEvents {
			if _, ok := info.events[id]; ok {
				continue
			}

			// If the handler does have the method but does not implement the one specified in Peex it is probably an
			// outdated handler method.
			methodName := "Handle" + strings.TrimPrefix(eventName, "event")
			if _, ok := t.MethodByName(methodName); ok {
				panic("incompatible handler method: " + methodName + " (is Peex or the handler outdated?)")
			}
		}

		// Make sure to increment the handlerId for the next handler!
		m.handlerNextId++
	}
	for _, p := range cfg.Providers {
		id := p.componentId(m)
		if _, ok := m.componentProvs[id]; ok {
			panic("cannot register multiple providers for the same component (" + p.componentName() + ")")
		}
		m.componentProvs[id] = p
	}
	return m
}

// Accept assigns a Session to a player. This also works for disconnected players or fake players. Initial components
// can be provided for the player to start with. The add function will be called on any component that implements Adder.
// Providing multiple components of the same type is not allowed and will return an error.
func (m *Manager) Accept(p *player.Player, components ...Component) (*Session, error) {
	m.sessionMu.Lock()
	defer m.sessionMu.Unlock()

	if _, ok := m.sessions[p.UUID()]; ok {
		return nil, errors.New("trying to handle a player that already has a handler")
	}
	s := &Session{
		m:          m,
		components: make(map[componentId]Component),
	}
	s.p.Store(p)
	p.Handle(s)
	// Insert all the components into the session. No mutex lock is needed, as it is not yet possible for any other
	// goroutine to have access to the session yet.
	for _, comp := range components {
		err := s.insertComponent(m.getComponentId(comp), comp)
		if err != nil {
			return nil, err
		}
	}
	m.sessions[p.UUID()] = s
	return s, nil
}

// Sessions returns every session currently stored in the manager.
func (m *Manager) Sessions() []*Session {
	m.sessionMu.RLock()
	sessions := make([]*Session, 0, len(m.sessions))
	for _, s := range m.sessions {
		sessions = append(sessions, s)
	}
	m.sessionMu.RUnlock()
	return sessions
}

// SessionFromPlayer returns a player's session.
func (m *Manager) SessionFromPlayer(p *player.Player) (*Session, bool) {
	return m.SessionFromUUID(p.UUID())
}

// SessionFromUUID returns the session of the player with the corresponding UUID. Only works for currently online
// players.
func (m *Manager) SessionFromUUID(id uuid.UUID) (*Session, bool) {
	m.sessionMu.RLock()
	s, ok := m.sessions[id]
	m.sessionMu.RUnlock()
	return s, ok
}

// QueryID executes a query on a player by their UUID, regardless of whether they are online or not. If the player is
// online with a stored session, components will be fetched from that session. If the player is not online, or one of
// the components is not present, those components will be loaded if they have a provider. In this case, if at least one
// component has no provider or there was a provider error, the query will not run. Loaded components will be saved
// again.
// Returns any error that occurred and whether the query ran. Should be handled independently.
func (m *Manager) QueryID(id uuid.UUID, queryFunc any) (bool, error) {
	info := m.makeQueryFuncInfo(queryFunc)

	m.sessionMu.RLock()
	defer m.sessionMu.RUnlock()
	s, hasSession := m.sessions[id]

	val := reflect.ValueOf(queryFunc)
	args := make([]reflect.Value, 0, len(info.params))
	var compSaveQueue []any
	var compSaveIds []componentId

	if hasSession {
		s.componentsMu.RLock()
		defer s.componentsMu.RUnlock()
	}
	// Retrieve or load all required components.
	for _, param := range info.params {
		c, ok, err := func() (any, bool, error) {
			// Case 1: the player is online and has the component.
			if hasSession {
				c, ok := s.components[param.cId]
				if ok {
					return c, true, nil
				}
			}

			// Case 2: the player is not online or does not have the component.
			p, ok := m.componentProvs[param.cId]
			if !ok {
				return nil, false, nil
			}

			v, err := p.loadNew(id)
			compSaveQueue = append(compSaveQueue, v)
			compSaveIds = append(compSaveIds, param.cId)
			if err != nil {
				return nil, false, fmt.Errorf("error loading component: %w", err)
			}
			// THe component was successfully loaded.
			return v, true, nil
		}()
		// See if we can turn the value into an argument, or if something went wrong.
		if ok && err == nil {
			if param.direct {
				args = append(args, reflect.ValueOf(c))
				continue
			}
			args = append(args, reflect.ValueOf(param.query.set(c)))
		} else if err != nil {
			return false, err
		} else if !ok {
			if !param.optional {
				return false, nil
			}
			args = append(args, reflect.ValueOf(param.query))
			continue
		}
	}

	val.Call(args)
	// Save all the components that are were loaded because of this query.
	for i, c := range compSaveQueue {
		p, ok := m.componentProvs[compSaveIds[i]]
		if !ok {
			// Should not happen: we should have already loaded this component using its provider.
			panic("component does not have a provider")
		}
		// Try actually save it
		err := p.save(id, c)
		if err != nil {
			return true, fmt.Errorf("error saving component: %w", err)
		}
	}
	return true, nil
}

// QueryAll runs a query on all currently online players. This works the same as if a query is run on every Session
// separately (albeit slightly faster). A number of players on which the query executed successfully is returned.
func (m *Manager) QueryAll(queryFunc any) int {
	info := m.makeQueryFuncInfo(queryFunc)

	count := 0
	m.sessionMu.RLock()
	for _, s := range m.sessions {
		if s.query(queryFunc, info) {
			count++
		}
	}
	m.sessionMu.RUnlock()
	return count
}
