package app

import (
	"github.com/mikestefanello/hooks"
	"github.com/samber/do"
)

var HookBoot = hooks.NewHook[*do.Injector]("boot")

func Boot() *do.Injector {
	HookBoot.Dispatch(do.DefaultInjector)
	return do.DefaultInjector
}
