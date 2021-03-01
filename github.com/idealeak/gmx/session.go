package gmx

import (
	"github.com/astaxie/beego/cache"
	"sync"
	"time"
)

const (
	SESSION_AUTH_IP = "_$IP$_"
)

var SessionMgr cache.Cache

func NewSession(id string) *Session {
	s := &Session{
		Id:         id,
		CreateTime: time.Now(),
	}
	return s
}

func GetSession(id string) *Session {
	s := SessionMgr.Get(id)
	if sess, ok := s.(*Session); ok {
		return sess
	}
	return nil
}

func SetSession(s *Session) {
	if s != nil {
		SessionMgr.Put(s.Id, s, int64(Cfg.SessionLifeTime))
	}
}

type Session struct {
	sync.Map
	Id         string
	CreateTime time.Time
}

func (s *Session) Get(key interface{}) (interface{}, bool) {
	return s.Load(key)
}

func (s *Session) Set(key, val interface{}) {
	s.Store(key, val)
}

func init() {
	SessionMgr, _ = cache.NewCache("memory", `{"interval":60}`)
}
