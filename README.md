# Peex

A multi-handler & player session system for Dragonfly, partly inspired by ECS and Dragonfly's command system.

Peex aims to keep a modular approach, without boilerplate code while keeping enough speed and simplicity.
I have personally tried multiple approaches for multiple handlers per player in the past,
from manually calling other handlers in the main handler to more sophisticated approaches.
Ultimately I think this approach is my favourite one so far.

### Experimental/concept branch for component loading/saving.

The aim of this branch is to allow for seamless use of data regardless of where the dats is or if the player is even 
online. Components can have providers that will be called whenever a component is added/removed.
There is a new `(*peex.Manager).QueryID()` function which can execute a query on both online and offline players.
Read the function's documentation for more info.

## How it works
This section will show the basics of how peex works.
The example used here will be a basic implementation of some sort of minigame system.

#### The manager & sessions
Firstly you will need to make a new `*peex.Manager`.
This will store all active sessions, and will allow you to assign a session to a player.

```go
manager := peex.New( /* ... handlers go here (more on that shortly). */ )

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

In our example you would add the component when a player joins a miningame and remove it when they leave t.
i
#### Handlers
Now that our player has components, we can write handlers to handle events for the player.
A handler is just a struct that implements some methods fom `player.Handler`.
Note that your handler does not actually need to implement `player.Handler`.
In fact, it is recommended to **not implement events you dont use** for performance reasons.

Struct fields can be used to add different queries to the handler.
The handler will only run if all the queried components are present in the session
and will also allow the handler to access these values.

Let's create a handler that will handle events when the player is in a minigame.
We will make a simple that subtracts score when the player dies.
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
Here didRun is a boolean that returns whether the the query was able to run or not.

You can also run queries on multiple players at once, using the manager.QueryAll() method.
This works the same as session.Query(), just for every player.
The method will return the amount of players the function actually ran for.