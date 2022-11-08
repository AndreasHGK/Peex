# Peex

A multi-handler & player session system for Dragonfly, partly inspired by ECS and Dragonfly's command system.

Peex aims to keep a modular approach, without boilerplate code while keeping enough speed and simplicity.
I have personally tried multiple approaches for multiple handlers per player in the past,
from manually calling other handlers in the main handler to more sophisticated approaches.
Ultimately I think this approach is my favourite one so far.

## How it works

### Basics
This section will show the basics of how peex works.
The example used here will be a basic implementation of some sort of minigame system.

#### The manager & sessions
Firstly you will need to make a new `*peex.Manager`.
This will store all active sessions, and will allow you to assign a session to a player.

```go
manager := peex.New(peex.Config{
	// ... (some fields are omitted)
	Handlers: []peex.Handler{ /* ... handlers go here (more on that shortly). */ }
})

// Ideally, run this when the player joins to assign them a session.
session := manager.Accept(player)
```
As can be seen in this example, you can provide all handlers that will run when creating the manager.
They **cannot** be added after it has been created.
You can still control when handlers run using components.
Let's go over those first before explaining handlers in more detail.

#### Components
Components are what actually stores a player's data.
A player can have multiple components, but they are stored by type so multiple components
of the same type is not possible.
They are usually simple structs with data, or pointers to ones.
Keep in mind that if your component is not a pointer it cannot be modified in handlers.

In our example, lets create a MinigamePlayer component.
```go
type MinigamePlayer struct {
    Game  *Minigame
    Score int
    Team  Team
}
```
That's all you need to do!
You can add any number of fields (or no fields), just like a normal struct.
To give a player this component, you can do the following:
```go
err := session.InsertComponent(&MinigamePlayer{
    // values...
})
```
The function will return an error if a player already has a component of said type.
Use `session.SetComponent(component)` to set or overwrite a component regardless of whether
it was already present.
Components can also be removed using `session.RemoveComponent(component)`.
This will remove the component with the same type as the argument, if it exists, and return it.

In our example you would add the component when a player joins a minigame and remove it when they leave it.

#### Handlers
Now that our player has components, we can write handlers to handle events for the player.
A handler is just a struct that implements some methods fom `player.Handler`.
Note that your handler does not actually need to implement `player.Handler`.
In fact, it is recommended to **not implement events you dont use** for performance reasons.

Struct fields can be used to add different queries to the handler.
The handler will only run if all the queried components are present in the session
and will also allow the handler to access these values.

Let's create a handler that will handle events when the player is in a minigame.
We will make a simple one that subtracts from the score when the player dies.
```go
type MinigameHandler struct {
    // peex will set the first *player.Player field it finds to the 
    // player that is the events. Has to be exported!
    Player  *player.Player
    Session *peex.Session // same as above but for *session.Session
    Manager *peex.Manager // ^
    
    // This parameter will make it so the handler only runs when the
    // specified component type is present. Different query types
    // also exist, like With if you do not wish to access any values
    // and Optional, which will make the handler run even if the
    // component is not present. All queries need to be exported!
    MinigamePlayer peex.Query[*MinigamePlayer]
    // You can add as many queries for different types as you like!
}

func (m MinigameHandler) HandleDeath() {
    m.MinigamePlayer.Load().Score -= 1
}
```
As seen before, handlers need to be registered when creating the manager.
This means you cannot remove handlers on runtime.
This should not be a problem due to the query system:
you can specify which handlers run by adding or removing components to/from a session.
When you register a handler to the manager,
it will automatically detect which events are implemented and only handle those events.

#### Query functions
Sometimes you want to run some logic on certain components, or only if certain
components are present.
You can either use `component, ok := session.Component(type)`,
or use the `session.Query(queryFuncion)` method.

A query function is similar to a handler: you can specify queries as function parameters,
and the query will only run if all component are present.
Lets run a query to change a player's team, which would for example be useful in a /changeteam command.
```go
didRun := session.Query(func(q1 peex.Query[*MinigamePlayer]) {
    q1.Load().Team = newTeam
})
```
Here didRun is a boolean that returns whether the query was able to run or not.
In query functions the `peex.Query[]` around the component type can be omitted.
When using another query type like Option, you will still need to include it.

You can also run queries on multiple players at once, using the manager.QueryAll() method.
This works the same as session.Query(), just for every player.
The method will return the amount of players the function actually ran for.

### Data Persistence

You may want to automatically load and save data for some components.
This is also supported in the library using component providers,
and working with persistent data for both online and offline players is similar to before.

#### Providers

A provider is just a struct that implements a Save and Load method for a component.
Say we want to have a provider for `*SampleComponent`, the provider would look something like this:
```go
type SampleProvider struct { /* ... */ }

func (SampleProvider) Load(id uuid.UUID, comp *SampleComponent) error {
	/* implementation ... */
}

func (SampleProvider) Save(id uuid.UUID, comp *SampleComponent) error {
    /* implementation ... */
}
```
Note that the component type must be a pointer.

After creating the provider for the component, you may register it to the manager by adding it to the `peex.Config`.
It needs to be wrapped in a `peex.ProviderWrapper` to allow peex to use the provider regardless of the component type
while keeping strict typing.
```go
manager := peex.New(peex.Config{
	// ... (some fields are omitted)
	Providers: []peex.ComponentProvider{
		peex.WrapProvider(SampleProvider{}),
		/* ... more providers go here. */
	}
})
```
Now, when the `SampleComponent` is inserted into a session,
the provider will first have its Load function called to load any data into the component.
When the component is removed, it will also be saved again.

Notice that we did not have to modify the actual component at all.
This allows for providers to be seamlessly swapped out.

#### UUID Queries

You may want to query a component regardless of whether the player's session currently has this component,
or even when the player is offline altogether.
For this you can perform a query function by a player UUID.

This is almost the same as a normal query, except it will try to load a component if it was not present in the session
or the player is offline. The query will not run if at least one component is not present in the session and it has no
provider, or if a component could not be loaded. 
Any loaded components will be saved again after the function has run.
An error is returned when there was an error loading or saving a component.
```go
didRun, err := manager.QueryID(func(q1 peex.Query[*SampleComponent]) {
    /* do stuff */
})
```
