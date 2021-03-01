package behaviago

import "strings"

/**
  Wait for the specified frames. always return Running until exceeds count.
*/
const (
	WaitFramesNodeName = "WaitFrames"
)

func init() {
	RegisteNodeCreator(WaitFramesNodeName, BehaviorNodeCreatorWrapper(func() BehaviorNode {
		wf := &WaitFrames{}
		return wf
	}))
}

func newWaitFramesNode() *WaitFrames {
	wf := &WaitFrames{}
	wf.SetClassName(WaitFramesNodeName)
	wf.SetSelf(wf)
	return wf
}

// ============================================================================
type WaitFrames struct {
	BehaviorNodeBase
	framesVar *Property
	framesMet *Method
}

func (w *WaitFrames) Load(version int, agentType string, properties []property_t) {
	w.BehaviorNodeBase.Load(version, agentType, properties)
	for i := 0; i < len(properties); i++ {
		if properties[i].name == "Frames" {
			if !strings.ContainsAny(properties[i].value, "(") {
				w.framesVar = ParseProperty(properties[i].value)
			} else {
				w.framesMet = LoadMethod(properties[i].value)
			}
		}
	}
}

func (w *WaitFrames) CreateTask() BehaviorTask {
	BTGLog.Trace("(w *WaitFrames) CreateTask()")
	return &WaitFramesTask{}
}

func (w *WaitFrames) IsValid(a *Agent, task BehaviorTask) bool {
	if _, ok := task.GetNode().(*WaitFrames); !ok {
		return false
	}
	return true
}

func (w *WaitFrames) GetFrames(a *Agent) int64 {
	if w.framesVar != nil {
		return w.framesVar.GetInt(a)
	} else if w.framesMet != nil {
		return w.framesMet.InvokeAndGetInt(a)
	}
	return 0
}

// ============================================================================
type WaitFramesTask struct {
	LeafTask
	start  int64
	frames int64
}

func (wft *WaitFramesTask) GetFrames(a *Agent) int64 {
	if w, ok := wft.GetNode().(*WaitFrames); ok {
		return w.GetFrames(a)
	}
	return 0
}
func (wft *WaitFramesTask) OnEnter(a *Agent) bool {
	wft.start = BTGWorkspace.GetFrameSinceStartup()
	wft.frames = wft.GetFrames(a)
	return wft.BehaviorTaskBase.OnEnter(a)
}

func (wft *WaitFramesTask) OnExit(a *Agent, status EBTStatus) {
	wft.BehaviorTaskBase.OnExit(a, status)
}

func (wtf *WaitFramesTask) Update(a *Agent, childStatus EBTStatus) EBTStatus {
	if BTGWorkspace.GetFrameSinceStartup()-wtf.start >= wtf.frames {
		return BT_SUCCESS
	}
	return BT_RUNNING
}
