package eventid

func GetEventId(name string) (EventId, bool) {
	i, ok := eventNames[name]
	return i, ok
}

func AllEvents() map[string]EventId {
	return eventNames
}
