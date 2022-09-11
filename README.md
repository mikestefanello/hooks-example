# Hooks (examples) - A modular monolithic approach

## Overview

Aside from just providing usage examples for the [hooks](https://github.com/mikestefanello/hooks) library, this is an exploration in modular monolithic architectural patterns in Go by leveraging both [hooks](https://github.com/mikestefanello/hooks) and [do](https://github.com/samber/do) _(for dependency injection)_.  It's recommended you review and understand these libraries prior to reviewing this repository. [Do](https://github.com/samber/do) is not required to achieve the pattern illustrated in this application, but I find it to be a very helpful and elegant approach.

I'm by no means advocating (at this time) for this specific approach but rather using this as an experiment and place to iterate with these ideas. I have had a lot of success with modular monoliths with languages and frameworks prior to learning Go and I haven't come across any similar patterns within the Go ecosystem. While microservices have become more prominent, a modular monolith can not only be a better choice in certain circumstances, but if done well, can make transitioning to microservices easier.

The overall goals of this approach are:
1) Create self-contained _modules_ that represent segments of business logic.
2) Avoid any patterns that reach across the codebase (ie, the entrypoint being used to initialize all dependencies, a router that initializes all handlers and routes, etc).
3) Modules should be able to be added and removed without having to touch the _core_ codebase at all.

## Repo structure

Below describes the repo structure and is just a proposed idea for effective, clear organization, but there's no requirement to follow this.

```
hooks-example/
├─ modules/         # Modules that each represent a unit of independent business logic 
│  ├─ analytics/
│  ├─ todo/
├─ pkg/             # General-purpose, non-dependency packages which can be used across the application 
│  ├─ app/
├─ services/        # Services which are auto-registered as dependencies
│  ├─ cache/
│  ├─ config/
│  ├─ web/
├─ main.go
```

### Modules

See the `func init()` within the primary, self-named _.go_ file of each module to understand how the module auto-registers itself with the application.

- `modules/todo`: Provides a very simple todo-list implemenation with a `Todo` model, a service to interact with todos as a registered dependency, an HTTP handler as a registered dependency, some JSON REST endpoints, and hooks to allow other modules to alter todos prior to saving and react when they are saved.
- `modules/analytics`: Provides bare-bones analytics for the application including the number of web requests received and the amount of entities created. Included is an `Analytics` model, a service to interact with analytics as a registered dependency, an HTTP handler as a registered dependency, middleware to track requests, a GET endpoint to return analytics, a hook to broadcast updates to the analytics, a listener for todo creation in order to track entities.

## Hooks

### Dispatchers

- `pkg/app`
  - `HookBoot`: Indicates that the application is booting and allow dependencies to be registered across the entire application via `*do.Injector`.
- `services/web`
  - `HookBuildRouter`: Dispatched when the web router is being built which allows listeners to register their own web routes and middleware.
- `modules/todo`
  - `HookTodoPreInsert`: Dispatched prior to inserting a new _todo_ which allows listeners to make any required modifications.
  - `HookTodoInsert`: Dispatched after a new _todo_ is inserted.
- `modules/analytics`
  - `HookAnalyticsUpdate`: Dispatched when the analytics data is updated.

### Listeners

- `HookBoot`
  - `services/cache`: Registers a cache backend as a dependency.
  - `services/config`: Registers configuration as a dependency.
  - `services/web`: Registers a web server as a dependency.
  - `modules/analytics`: Registers analytics service and HTTP handler as dependencies.
  - `modules/todo`: Registers todo service and HTTP handler as dependencies.
- `HookBuildRouter`
  - `modules/analytics`: Registers web route and tracker middleware for analytics.
  - `modules/todo`: Registers web routes for todos.
- `HookTodoInsert`
  - `modules/analytics`: Increments analytics entity count when todos are created.

## Boot process and registration

Below is an attempt to illustrate how the entire application self-registers starting from a single hook that is invoked.

### Code

```go
func main() {
  i := app.Boot()
  
  server := do.MustInvoke[web.Web](i)
  _ = server.Start()
}
```

### Walkthrough

```
main.go/              app.Boot()
├─ pkg/app.go:        [Dispatch] HookBoot 
├─ services/cache.go  ├─  Register dependency: *cache.Cache
├─ services/config.go ├─  Register dependency: *config.Config
├─ services/web.go    ├─  Register dependency: *web.Web
├─ modules/analytics: ├─  Register dependency: *analytics.Service
├─ modules/analytics: ├─  Register dependency: *analytics.Handler
├─ modules/todo:      ├─  Register dependency: *todo.Service
├─ modules/todo:      ├─  Register dependency: *todo.Handler

main.go/              server := do.MustInvoke[web.Web](i)
├─ services/web.go:   ├─  Initialize *web.Web
├                     ├───  Initialize *config.Config
├                     ├───  [Dispatch] HookRouterBuild
├─ modules/analytics:      ├─  Register web routes and middleware
├                          ├───  Initialize *analytics.Handler
├                          ├─────  Initialize *analytics.Service
├                          ├───────  Initialize *cache.Cache  
├─ modules/todo:           ├─  Register web routes
├                          ├───  Initialize *todo.Handler
├                          ├─────  Initialize *todo.Service
├                          ├───────  Initialize *cache.Cache  
```

## Imports

It's important to note that if you want a module or service to self-register, it must be imported. This is why you see this in `main.go`:

```go
// Services
_ "github.com/mikestefanello/hooks-example/services/cache"
_ "github.com/mikestefanello/hooks-example/services/config"
"github.com/mikestefanello/hooks-example/services/web"
// Modules
_ "github.com/mikestefanello/hooks-example/modules/analytics"
_ "github.com/mikestefanello/hooks-example/modules/todo"
```

This is needed to ensure that `init()` executes in each package which is what they are using to listen to hooks.

### Logs

To help illustrate the app boot process:

```go
2022/09/09 15:50:22 hook created: boot
2022/09/09 15:50:22 registered listener with hook: boot
2022/09/09 15:50:22 registered listener with hook: boot
2022/09/09 15:50:22 hook created: router.build
2022/09/09 15:50:22 registered listener with hook: boot
2022/09/09 15:50:22 hook created: todo.pre_insert
2022/09/09 15:50:22 hook created: todo.insert
2022/09/09 15:50:22 registered listener with hook: boot
2022/09/09 15:50:22 registered listener with hook: router.build
2022/09/09 15:50:22 hook created: analytics.update
2022/09/09 15:50:22 registered listener with hook: boot
2022/09/09 15:50:22 registered listener with hook: router.build
2022/09/09 15:50:22 registered listener with hook: todo.insert
2022/09/09 15:50:22 dispatching hook boot to 5 listeners (async: false)
2022/09/09 15:50:22 dispatch to hook boot complete
2022/09/09 15:50:22 registered 7 dependencies: [*analytics.Handler *cache.Cache *config.Config *web.Web *todo.Service *todo.Handler *analytics.Service]
2022/09/09 15:50:22 dispatching hook router.build to 2 listeners (async: false)
2022/09/09 15:50:22 dispatch to hook router.build complete
2022/09/09 15:50:22 registered 5 routes: [GET_/ GET_/todo GET_/todo/:todo POST_/todo GET_/analytics]
```

### Module registration

Below is the code used by the `analytics` module to register itself:

```go
func init() {
    // Provide dependencies during app boot process
    app.HookBoot.Listen(func(e hooks.Event[*do.Injector]) {
        do.Provide(e.Msg, NewAnalyticsService)
        do.Provide(e.Msg, NewAnalyticsHandler)
    })

    // Provide web routes
    web.HookBuildRouter.Listen(func(e hooks.Event[*echo.Echo]) {
        h := do.MustInvoke[Handler](do.DefaultInjector)
        e.Msg.GET("/analytics", h.Get)
        e.Msg.Use(h.WebRequestMiddleware)
    })

    // React to new todos being inserted
    todo.HookTodoInsert.Listen(func(e hooks.Event[todo.Todo]) {
        h := do.MustInvoke[Service](do.DefaultInjector)
        if err := h.IncrementEntities(); err != nil {
            log.Error(err)
        }
    })
}
```

## Optional independent binaries

It is possible to create separate entrypoints that only register one or some of your modules, allowing for a monolithic codebase that could be used to create separate applications/services.

For example, in `main.go`, simply remove the import `_ "github.com/mikestefanello/hooks-example/modules/analytics"` and the application will run without the `analytics` modules (and everything within it).

## Run the application

`go run main.go`

### Endpoints

_NOTE:_ Data created is stored in memory and will be lost when the application restarts.

- `GET /`: Hello world
- `GET /todo`: Get all _todos_
- `GET /todo/:todo`: Get a _todo_ by ID
- `POST /todo`: Create a todo
- `GET /analytics`: Get analytics

## Downsides

Nothing is without downsides and this approach certainly has them. It lacks overall explicitness by hiding details within hook listeners and by injecting all dependencies inside a single container. It could make understanding and debugging the codebase harder than one following a very straight-forward approach, especially since you lose some power of your IDE. This certainly goes a bit against the overall philosophy of Go itself. It's also hard to tell how well this would scale with a large codebase and even with multiple development teams.

But there are pros, in my opinion. I'll leave it to the reader to make their own judgements and I encourage you to share them here.