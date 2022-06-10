package peex

import (
	"github.com/andreashgk/peex/eventid"
	"github.com/df-mc/atomic"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/google/uuid"
	"reflect"
)

// Manager stores all current sessions. It also contains all the registered handlers and component types.
type Manager struct {
	sessions map[uuid.UUID]*Session

	nextHandlerId handlerId
	handlerIds    map[reflect.Type]handlerId
	handlers      map[handlerId]handlerInfo
	eventHandlers map[eventid.EventId][]handlerId

	componentNextId  componentId
	componentIdTable map[reflect.Type]componentId

	started atomic.Value[bool]
}

// New creates a new Session Manager. It also inserts all the provided handlers into the manager. Events will be called
// in the order the handlers are added to the manager. These handlers can query for specific components to be present in
// a player Session in order to actually run.
func New(handlers ...Handler) *Manager {
	m := &Manager{
		sessions:         make(map[uuid.UUID]*Session),
		handlerIds:       make(map[reflect.Type]handlerId),
		handlers:         make(map[handlerId]handlerInfo),
		eventHandlers:    make(map[eventid.EventId][]handlerId),
		componentIdTable: make(map[reflect.Type]componentId),
	}
	for _, id := range eventid.AllEvents() {
		m.eventHandlers[id] = []handlerId{}
	}

	for _, h := range handlers {
		t := reflect.TypeOf(h)
		if _, ok := m.handlerIds[t]; ok {
			panic("re-registering an existing handler type")
		}

		// Assign the handler ID to the type, and generate the handlerInfo
		m.handlerIds[t] = m.nextHandlerId
		info := m.createHandlerInfo(h)
		m.handlers[m.nextHandlerId] = info
		info.events.Each(func(id eventid.EventId) {
			m.eventHandlers[id] = append(m.eventHandlers[id], m.nextHandlerId)
		})

		// Make sure to increment the handlerId for the next handler!
		m.nextHandlerId++
	}
	return m
}

// Accept assigns a Session to a player.
func (m *Manager) Accept(p *player.Player) *Session {
	m.started.Store(true)

	s := &Session{
		p:          p,
		m:          m,
		components: make(map[componentId]Component),
	}
	p.Handle(s)
	m.sessions[p.UUID()] = s
	return s
}

// QueryAll runs a query on all currently online players. This works the same as if a query is run on every Session
// separately (albeit slightly faster). A number of players on which the query executed successfully is returned.
func (m *Manager) QueryAll(queryFunc any) int {
	info := m.makeQueryFuncInfo(queryFunc)

	count := 0
	for _, s := range m.sessions {
		if s.query(queryFunc, info) {
			count++
		}
	}
	return count
}
