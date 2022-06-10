package peex

import (
	"fmt"
	"reflect"
)

// Query is used to query for a certain component type, passing its value to the handler. Query.Load() can be used to
// get the value.
type Query[c Component] struct {
	With[c]
	val c
}

// Option is a query that allows the query to specify whether a certain Component is present in the session, passing
// along its value if it exists. This value can be accessed through Option.Load().
type Option[c Component] struct {
	With[c]
	val c
	has bool
}

// With is used in handler queries to denote that the presence of a certain component type is required, but the
// value of the component itself is not important.
type With[c Component] struct{}

// Load returns the underlying value of the Query.
func (q Query[c]) Load() c {
	return q.val
}

// Load returns the queried value, along with whether it actually exists.
func (o Option[c]) Load() (c, bool) {
	return o.val, o.has
}

/// Internal query logic
/// --------------------

func (q Query[c]) set(x any) queryType {
	q.val = x.(c)
	return q
}

func (w With[c]) getType() reflect.Type {
	v := new(c)
	return reflect.TypeOf(v).Elem()
}

func (w With[c]) optional() bool {
	return false
}

func (w With[c]) set(x any) queryType { return w }

func (o Option[c]) optional() bool {
	return true
}

func (o Option[c]) set(x any) queryType {
	o.val = x.(c)
	o.has = true
	return o
}

type queryType interface {
	getType() reflect.Type
	optional() bool
	set(x any) queryType
}

// query function stuff

type queryFuncInfo struct {
	params []queryFuncParam
}

type queryFuncParam struct {
	cId      componentId
	optional bool

	query queryType
}

func (m *Manager) makeQueryFuncInfo(f any) queryFuncInfo {
	t := reflect.TypeOf(f)
	if t.Kind() != reflect.Func {
		panic(fmt.Errorf("expected a function, got %s", t.String()))
	}

	info := queryFuncInfo{}
	for i := 0; i < t.NumIn(); i++ {
		in := t.In(i)
		query, ok := reflect.New(in).Interface().(queryType)
		if !ok {
			panic("query func must only have query types (Query, With, Option)")
		}

		cId, ok := m.componentIdTable[query.getType()]
		// If the component is not registered, the player does not have this component
		if !ok {
			continue
		}

		info.params = append(info.params, queryFuncParam{
			cId:      cId,
			optional: query.optional(),
			query:    query,
		})
	}
	return info
}
