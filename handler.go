package peex

import (
	"errors"
	"github.com/df-mc/dragonfly/server/player"
	"reflect"
)

// Handler is a struct that handles player-related events. It can query for certain components contained in the player
// session, and will only run if all those components are present.
type Handler interface {
}

/// Internal handler logic
/// ----------------------

type handlerId uint

// handlerInfo contains data about a specific handler type.
type handlerInfo struct {
	h   Handler
	typ reflect.Type

	components []componentQuery
	events     map[eventId]struct{}

	playerField  int
	sessionField int
	managerField int
}

type componentQuery struct {
	id       componentId
	fieldNum int
	optional bool
}

// createHandlerInfo creates a new handler info struct for a type of handler.
func (m *Manager) createHandlerInfo(h Handler) handlerInfo {
	v := reflect.ValueOf(h)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		panic(errors.New("A handler must be of type struct, or a pointer to a one"))
	}

	info := handlerInfo{
		h:            h,
		typ:          reflect.TypeOf(h),
		events:       getHandlerEvents(h),
		playerField:  -1,
		sessionField: -1,
		managerField: -1,
	}
	for i := 0; i < v.NumField(); i++ {
		// Fields marked with a `ignore:""` tag will be ignored by the library.
		if _, ok := v.Type().Field(i).Tag.Lookup("ignore"); ok {
			continue
		}

		// Ignore unexported fields
		if !v.Field(i).CanInterface() {
			m.logger.Debugf("warning: unexported handler fields cannot be copied by Peex")
			continue
		}
		switch x := v.Field(i).Interface().(type) {
		case queryType:
			fieldType := x.getType()

			cId := m.getComponentIdRefl(fieldType)
			info.components = append(info.components, componentQuery{
				id:       cId,
				fieldNum: i,
				optional: x.optional(),
			})
		case *player.Player:
			// We don't need to pass the same player multiple times
			if info.playerField != -1 {
				continue
			}
			info.playerField = i
		case *Session:
			// We don't need to pass the same session multiple times
			if info.sessionField != -1 {
				continue
			}
			info.sessionField = i
		case *Manager:
			// We don't need to pass the same manager multiple times
			if info.managerField != -1 {
				continue
			}
			info.managerField = i
		default:
			continue
		}
	}

	return info
}

// handleEvent handles all shared logic for events, such as assigning query values.
func (s *Session) handleEvent(eventId eventId, f func(h Handler)) {
	s.componentsMu.RLock()
	defer s.componentsMu.RUnlock()
handlerLoop:
	for _, id := range s.m.eventHandlers[eventId] {
		info := s.m.handlers[id]

		// Figure out which components to set in the queries
		comps := make([]componentQuery, 0, len(info.components))
		for _, compQuery := range info.components {
			_, isPresent := s.components[compQuery.id]
			if !isPresent && !compQuery.optional {
				continue handlerLoop
			} else if !isPresent && compQuery.optional {
				continue
			}

			comps = append(comps, compQuery)
		}

		actualType := reflect.New(info.typ).Elem()
		structType := actualType
		if actualType.Kind() == reflect.Pointer {
			val := reflect.New(info.typ.Elem())
			actualType.Set(val.Elem().Addr())

			structType = actualType.Elem()
		}

		for _, compQuery := range comps {
			field := structType.Field(compQuery.fieldNum)
			query := field.Interface().(queryType)

			field.Set(reflect.ValueOf(query.set(s.components[compQuery.id].(Component))))
		}

		if info.playerField != -1 {
			structType.Field(info.playerField).Set(reflect.ValueOf(s.Player()))
		}
		if info.sessionField != -1 {
			structType.Field(info.sessionField).Set(reflect.ValueOf(s))
		}
		if info.managerField != -1 {
			structType.Field(info.managerField).Set(reflect.ValueOf(s.m))
		}

		f(actualType.Interface().(Handler))
	}
}
