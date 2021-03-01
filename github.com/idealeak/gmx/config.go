package gmx

import "net"

var Cfg Config

type Config struct {
	Addr            string   //gmx api地址
	Username        string   //gmx 交互用的用户名
	Password        string   //gmx 交互用的密码
	SessionLifeTime int      //gmx session生命周期,单位:秒
	Timeout         int      //gmx api超时时间,单位:秒
	WhiteAddr       []string //gmx api交互白名单地址
	Debug           bool     //gmx debug模式,log开关
	EnableOp        bool     //gmx api交互开关
}

func (cfg *Config) IsValid() bool {
	if cfg.Username == "" || cfg.Password == "" {
		return false
	}
	if cfg.SessionLifeTime <= 0 {
		return false
	}
	if cfg.Timeout <= 0 {
		return false
	}
	if _, _, err := net.SplitHostPort(cfg.Addr); err != nil {
		return false
	}
	return true
}
