package behaviago

import (
	"strings"
	"time"
)

/**
  Wait for the specified milliseconds. always return Running until time over.
*/
const (
	WaitNodeName = "Wait"
)

func init() {
	RegisteNodeCreator(WaitNodeName, BehaviorNodeCreatorWrapper(func() BehaviorNode {
		w := newWaitNode()
		return w
	}))
}

func newWaitNode() *Wait {
	w := &Wait{}
	w.SetClassName(WaitNodeName)
	w.SetSelf(w)
	return w
}

// ============================================================================
type Wait struct {
	BehaviorNodeBase
	timeVar *Property
	timeMet *Method
}

func (w *Wait) Load(version int, agentType string, properties []property_t) {
	w.BehaviorNodeBase.Load(version, agentType, properties)
	for i := 0; i < len(properties); i++ {
		if properties[i].name == "Time" {
			if !strings.ContainsAny(properties[i].value, "(") {
				w.timeVar = ParseProperty(properties[i].value)
			} else {
				w.timeMet = LoadMethod(properties[i].value)
			}
		}
	}
}

func (w *Wait) CreateTask() BehaviorTask {
	BTGLog.Trace("(w *Wait) CreateTask()")
	return &WaitTask{}
}

func (w *Wait) IsValid(a *Agent, task BehaviorTask) bool {
	if _, ok := task.GetNode().(*Wait); !ok {
		return false
	}
	return true
}

func (w *Wait) GetTime(a *Agent) float64 {
	BTGLog.Trace("(w *Wait) GetTime() enter")
	if w.timeVar != nil {
		return w.timeVar.GetFloat(a)
	} else if w.timeMet != nil {
		return w.timeMet.InvokeAndGetFloat(a)
	}
	return .0
}

// ============================================================================
type WaitTask struct {
	LeafTask
	start int64
	time  int64
}

func (wt *WaitTask) GetTime(a *Agent) float64 {
	BTGLog.Trace("(wt *WaitTask) GetTime() enter")
	if w, ok := wt.GetNode().(*Wait); ok {
		return w.GetTime(a)
	}
	return .0
}

func (wt *WaitTask) OnEnter(a *Agent) bool {
	wt.start = BTGWorkspace.GetTimeSinceStartup()
	wt.time = int64(wt.GetTime(a)) * int64(time.Millisecond)
	return wt.BehaviorTaskBase.OnEnter(a)
}

func (wt *WaitTask) OnExit(a *Agent, status EBTStatus) {
	wt.BehaviorTaskBase.OnExit(a, status)
}

func (wt *WaitTask) Update(a *Agent, childStatus EBTStatus) EBTStatus {
	BTGLog.Trace("(wt *WaitTask) Update() enter")
	if BTGWorkspace.GetTimeSinceStartup()-wt.start >= wt.time {
		return BT_SUCCESS
	}
	return BT_RUNNING
}
