package behaviago

import "github.com/howeyc/fsnotify"

type FileSysMonitor struct {
	started    bool
	monitorDir string
	files      []string
	watcher    *fsnotify.Watcher
	innerQ     chan struct{}
}

func NewFileSysMonitor() *FileSysMonitor {
	fsm := &FileSysMonitor{
		started: false,
		innerQ:  make(chan struct{}),
	}
	return fsm
}

func (fsm *FileSysMonitor) SetMonitorDir(dir string) {
	fsm.monitorDir = dir
}

func (fsm *FileSysMonitor) Start() error {
	if fsm.started {
		return nil
	}
	fsm.started = true
	var err error
	fsm.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		BTGLog.Error("fsnotify.NewWatcher err:", err)
		return err
	}
	// Process events
	go func() {
		for {
			select {
			case ev, ok := <-fsm.watcher.Event:
				if ok && ev != nil {
					if ev.IsModify() {
						BTGLog.Info("fsnotify event:", ev)
						//todo:change
					}
				} else {
					fsm.watcher.RemoveWatch(fsm.monitorDir)
					return
				}
			case err := <-fsm.watcher.Error:
				BTGLog.Error("fsnotify error:", err)
				fsm.watcher.RemoveWatch(fsm.monitorDir)
				return
			case <-fsm.innerQ:
				fsm.watcher.RemoveWatch(fsm.monitorDir)
				return
			}
		}
	}()
	fsm.watcher.Watch(fsm.monitorDir)

	return nil
}

func (fsm *FileSysMonitor) Stop() {
	if !fsm.started {
		return
	}
	fsm.innerQ <- struct{}{}
	fsm.started = false
}

func (fsm *FileSysMonitor) Restart() error {
	fsm.Stop()
	return fsm.Start()
}
