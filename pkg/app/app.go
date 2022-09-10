package app

import (
	"log"

	"github.com/mikestefanello/hooks"
	"github.com/samber/do"
)

// HookBoot allows modules and service the ability to initialize and register dependencies
var HookBoot = hooks.NewHook[*do.Injector]("boot")

func init() {
	// See what's happening within hooks in the logs
	hooks.SetLogger(func(format string, args ...any) {
		log.Printf(format+"\n", args...)
	})
}

// Boot boots up the entire application by requesting all dependencies be registered
func Boot() *do.Injector {
	// Allow modules and service the ability to initialize and register dependencies
	HookBoot.Dispatch(do.DefaultInjector)

	// Log the dependencies provided via the hook
	d := do.DefaultInjector.ListProvidedServices()
	log.Printf("registered %d dependencies: %v", len(d), d)

	return do.DefaultInjector
}
