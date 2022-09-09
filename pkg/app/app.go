package app

import (
	"log"

	"github.com/mikestefanello/hooks"
	"github.com/samber/do"
)

var HookBoot = hooks.NewHook[*do.Injector]("boot")

func init() {
	// See what's happening within hooks in the logs
	hooks.SetLogger(func(format string, args ...any) {
		log.Printf(format+"\n", args...)
	})
}

// Boot boots up the entire application by requesting all dependencies be registered
func Boot() *do.Injector {
	HookBoot.Dispatch(do.DefaultInjector)

	log.Printf("registered %d dependencies: %v",
		len(do.DefaultInjector.ListProvidedServices()),
		do.DefaultInjector.ListProvidedServices(),
	)

	return do.DefaultInjector
}
