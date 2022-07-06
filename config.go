package peex

import "github.com/df-mc/dragonfly/server"

// Config is a struct passed to the New function when creating a new Manager. It allows for customizing several aspects
// such as specifying handlers and component providers. Many of these cannot be modified after the manager has been
// created.
type Config struct {
	// Logger allows for an optional logger to be supplied. This will log things such as errors when saving components
	// when a player is leaving.
	Logger server.Logger
	// Handlers contains all the handlers that will run during the lifetime of the manager. These will always be active,
	// but can be controlled through adding or removing components from users.
	Handlers []Handler
	// Providers allows for passing of a list of ComponentProviders which can load & save components for players at
	// runtime. The providers must be wrapped in a ProviderWrapper using the WrapProvider function.
	Providers []ComponentProvider
}
