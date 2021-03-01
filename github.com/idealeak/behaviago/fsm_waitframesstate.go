package behaviago

import (
	"strings"
)

const (
	WaitFramesStateNodeName = "WaitFramesState"
)

func init() {
	RegisteNodeCreator(WaitFramesStateNodeName, BehaviorNodeCreatorWrapper(func() BehaviorNode {
		n := newWaitFramesStateNode()
		return n
	}))
}

func newWaitFramesStateNode() *WaitFramesState {
	n := &WaitFramesState{}
	n.SetClassName(WaitFramesStateNodeName)
	n.SetSelf(n)
	return n
}

// ============================================================================
type WaitFramesState struct {
	State
	framesVar *Property
	framesMet *Method
}

func (ws *WaitFramesState) Load(version int, agentType string, properties []property_t) {
	ws.State.Load(version, agentType, properties)
	for i := 0; i < len(properties); i++ {
		if properties[i].name == "Frames" {
			if strings.ContainsAny(properties[i].value, "(") {
				ws.framesMet = LoadMethod(properties[i].value)
			} else {
				ws.framesVar = ParseProperty(properties[i].value)
			}
		}
	}
}

func (ws *WaitFramesState) GetFrames(a *Agent) int64 {
	if ws.framesVar != nil {
		return ws.framesVar.GetInt(a)
	} else if ws.framesMet != nil {
		return ws.framesMet.InvokeAndGetInt(a)
	}
	return 0
}

func (ws *WaitFramesState) IsValid(a *Agent, task BehaviorTask) bool {
	if _, ok := task.GetNode().(*WaitFramesState); !ok {
		return false
	}
	return true
}

func (ws *WaitFramesState) CreateTask() BehaviorTask {
	BTGLog.Trace("(ws *WaitFramesState) CreateTask()")
	return &WaitFramesStateTask{}
}

// ============================================================================
type WaitFramesStateTask struct {
	StateTask
	start  int64
	frames int64
}

func (wft *WaitFramesStateTask) GetFrames(a *Agent) int64 {
	if node, ok := wft.GetNode().(*WaitFramesState); ok {
		return node.GetFrames(a)
	}
	return 0
}

func (wft *WaitFramesStateTask) OnEnter(a *Agent) bool {
	wft.nextStateId = -1
	wft.start = BTGWorkspace.GetFrameSinceStartup()
	wft.frames = wft.GetFrames(a)
	return wft.StateTask.OnEnter(a)
}

func (wft *WaitFramesStateTask) OnExit(a *Agent, status EBTStatus) {
	wft.StateTask.OnExit(a, status)
}

func (wtf *WaitFramesStateTask) Update(a *Agent, childStatus EBTStatus) EBTStatus {
	if BTGWorkspace.GetFrameSinceStartup()-wtf.start >= wtf.frames {
		if node, ok := wtf.GetNode().(*WaitFramesState); ok {
			_, wtf.nextStateId = node.UpdateNext(a)
		}
		return BT_SUCCESS
	}
	return BT_RUNNING
}
