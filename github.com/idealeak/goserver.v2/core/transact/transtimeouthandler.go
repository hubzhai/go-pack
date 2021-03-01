package transact

import (
	"github.com/idealeak/goserver.v2/core/timer"
)

type transactTimerAction struct {
}

func (t transactTimerAction) OnTimer(h timer.TimerHandle, ud interface{}) bool {
	if trans, ok := ud.(*TransNode); ok {
		trans.timeout()
		return true
	}
	return false
}
