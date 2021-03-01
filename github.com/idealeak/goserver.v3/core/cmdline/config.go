package cmdline

import (
	"github.com/idealeak/goserver.v3/core"
)

var Config = Configuration{}

type Configuration struct {
	SupportCmdLine bool
}

func (c *Configuration) Name() string {
	return "cmdline"
}

func (c *Configuration) Init() error {
	return nil
}

func (c *Configuration) Close() error {
	return nil
}

func init() {
	core.RegistePackage(&Config)
}
