# Hooks (examples) - A modular monolithic approach

## Overview

Aside from just providing usage examples for the [hooks](https://github.com/mikestefanello/hooks) library, this is an exploration in modular monolithic architectural patterns in Go by leveraging both [hooks](https://github.com/mikestefanello/hooks) and [do](https://github.com/samber/do) _(for dependency injection)_.  It's recommended you review and understand these libraries prior to reviewing this repository.

@todo

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

@todo

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

## Optional independent binaries

@todo

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

@todo