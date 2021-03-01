package main

import (
	"games.agamestudio.com/jxjyhj/common"
	"github.com/idealeak/goserver.v3/core"
	"github.com/idealeak/goserver.v3/core/module"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().Unix())
	core.RegisterConfigEncryptor(common.ConfigFE)
	defer core.ClosePackages()
	core.LoadPackages("config.json")

	waiter := module.Start()
	waiter.Wait("main()")
}
