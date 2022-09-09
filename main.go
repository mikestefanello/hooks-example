package main

import (
	"log"

	"github.com/mikestefanello/hooks"
	"github.com/samber/do"

	"github.com/mikestefanello/hooks-example/pkg/app"

	// All services
	_ "github.com/mikestefanello/hooks-example/services/cache"
	_ "github.com/mikestefanello/hooks-example/services/config"
	"github.com/mikestefanello/hooks-example/services/web"
	// All modules
	_ "github.com/mikestefanello/hooks-example/modules/analytics"
	_ "github.com/mikestefanello/hooks-example/modules/todo"
)

func init() {
	// See what's happening within hooks in the logs
	hooks.SetLogger(func(format string, args ...any) {
		log.Printf(format+"\n", args...)
	})
}

func main() {
	i := app.Boot()

	server := do.MustInvoke[web.Web](i)
	server.Start()
}
