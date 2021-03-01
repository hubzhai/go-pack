package main

import (
	_ "github.com/idealeak/goserver.v3/mmo"

	"github.com/idealeak/goserver.v3/core"
	"github.com/idealeak/goserver.v3/core/module"
)

func main() {
	defer core.ClosePackages()
	core.LoadPackages("config.json")

	waiter := module.Start()
	waiter.Wait()
}
