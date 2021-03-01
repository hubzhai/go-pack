package behaviago

import (
	"strings"
	"time"
)

const (
	WaitStateNodeName = "WaitState"
)

func init() {
	RegisteNodeCreator(WaitStateNodeName, BehaviorNodeCreatorWrapper(func() BehaviorNode {
		n := newWaitStateNode()
		return n
	}))
}

func newWaitStateNode() *WaitState {
	n := &WaitState{}
	n.SetClassName(WaitStateNodeName)
	n.SetSelf(n)
	return n
}

// ============================================================================
type WaitState struct {
	State
	timeVar *Property
	timeMet *Method
}

func (ws *WaitState) Load(version int, agentType string, properties []property_t) {
	ws.State.Load(version, agentType, properties)
	for i := 0; i < len(properties); i++ {
		if properties[i].name == "Time" {
			if strings.ContainsAny(properties[i].value, "(") {
				ws.timeMet = LoadMethod(properties[i].value)
			} else {
				ws.timeVar = ParseProperty(properties[i].value)
			}
		}
	}
}

func (ws *WaitState) IsValid(a *Agent, task BehaviorTask) bool {
	if _, ok := task.GetNode().(*WaitState); !ok {
		return false
	}
	return true
}

func (ws *WaitState) GetTime(a *Agent) float64 {
	if ws.timeVar != nil {
		return ws.timeVar.GetFloat(a)
	} else if ws.timeMet != nil {
		return ws.timeMet.InvokeAndGetFloat(a)
	}
	return .0
}

func (ws *WaitState) CreateTask() BehaviorTask {
	BTGLog.Trace("(ws *WaitState) CreateTask()")
	return &WaitStateTask{}
}

// ============================================================================
type WaitStateTask struct {
	StateTask
	start int64
	time  int64
}

func (wt *WaitStateTask) GetTime(a *Agent) float64 {
	if w, ok := wt.GetNode().(*WaitState); ok {
		return w.GetTime(a)
	}
	return 0
}

func (wt *WaitStateTask) OnEnter(a *Agent) bool {
	wt.nextStateId = -1
	wt.start = BTGWorkspace.GetTimeSinceStartup()
	wt.time = int64(wt.GetTime(a)) * int64(time.Millisecond)
	return wt.StateTask.OnEnter(a)
}

func (wt *WaitStateTask) OnExit(a *Agent, status EBTStatus) {
	wt.StateTask.OnExit(a, status)
}

func (wt *WaitStateTask) Update(a *Agent, childStatus EBTStatus) EBTStatus {
	if BTGWorkspace.GetTimeSinceStartup()-wt.start >= wt.time {
		return BT_SUCCESS
	}
	return BT_RUNNING
}
