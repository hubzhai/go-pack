package gmx

import (
	"net/http"
	"sync"
)

// 线程安全沙箱,确保对注册对象的访问发生在安全沙箱里
type ThreadSafeSandBox func(*OpEvent) error

var threadSafeSandBox ThreadSafeSandBox
var gsrv *http.Server
var lock sync.Mutex

func RegisteThreadSafeSandbox(sb ThreadSafeSandBox) {
	lock.Lock()
	threadSafeSandBox = sb
	lock.Unlock()
}

func Startup(cfg Config) error {
	lock.Lock()
	if gsrv != nil {
		gsrv.Close()
		gsrv = nil
	}
	Cfg = cfg
	gsrv = &http.Server{Addr: Cfg.Addr}
	lock.Unlock()
	return gsrv.ListenAndServe()
}

func Close() error {
	lock.Lock()
	defer lock.Unlock()
	if gsrv != nil {
		srv := gsrv
		gsrv = nil
		return srv.Close()
	}
	return nil
}

func IsRunning() bool {
	lock.Lock()
	defer lock.Unlock()
	if gsrv == nil {
		return false
	}
	return true
}

func EnableOp() {
	lock.Lock()
	Cfg.EnableOp = true
	lock.Unlock()
}

func DisableOp() {
	lock.Lock()
	Cfg.EnableOp = false
	lock.Unlock()
}
