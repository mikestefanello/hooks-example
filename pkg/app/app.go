package app

import (
	"log"

	"github.com/mikestefanello/hooks"
	"github.com/samber/do"
)

var HookBoot = hooks.NewHook[*do.Injector]("boot")

func Boot() *do.Injector {
	HookBoot.Dispatch(do.DefaultInjector)

	log.Printf("registered dependencies: %v", do.DefaultInjector.ListProvidedServices())

	return do.DefaultInjector
}
