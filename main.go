package main

import (
	"github.com/samber/do"

	"github.com/mikestefanello/hooks-example/pkg/app"

	// Services
	_ "github.com/mikestefanello/hooks-example/services/cache"
	_ "github.com/mikestefanello/hooks-example/services/config"
	"github.com/mikestefanello/hooks-example/services/web"
	// Modules
	_ "github.com/mikestefanello/hooks-example/modules/analytics"
	_ "github.com/mikestefanello/hooks-example/modules/todo"
)

func main() {
	i := app.Boot()

	server := do.MustInvoke[web.Web](i)
	_ = server.Start()
}
