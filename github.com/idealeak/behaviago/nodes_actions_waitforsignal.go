package behaviago

/**
  Always return Running until the predicates of WaitforSignal node become true,
  or executing child node and return execution result.
*/
const (
	WaitforSignalNodeName = "WaitforSignal"
)

func init() {
	RegisteNodeCreator(WaitforSignalNodeName, BehaviorNodeCreatorWrapper(func() BehaviorNode {
		wfs := newWaitforSignalNode()
		return wfs
	}))
}

func newWaitforSignalNode() *WaitforSignal {
	wfs := &WaitforSignal{}
	wfs.SetClassName(WaitforSignalNodeName)
	wfs.SetSelf(wfs)
	return wfs
}

// ============================================================================
type WaitforSignal struct {
	BehaviorNodeBase
}

func (w *WaitforSignal) Load(version int, agentType string, properties []property_t) {
	w.BehaviorNodeBase.Load(version, agentType, properties)
}

func (w *WaitforSignal) CreateTask() BehaviorTask {
	BTGLog.Trace("(w *WaitforSignal) CreateTask()")
	return &WaitforSignalTask{}
}

func (w *WaitforSignal) IsValid(a *Agent, task BehaviorTask) bool {
	if _, ok := task.GetNode().(*WaitforSignal); !ok {
		return false
	}
	return true
}

func (w *WaitforSignal) CheckIfSignaled(a *Agent) bool {
	return w.EvaluteCustomCondition(a)
}

// ============================================================================
type WaitforSignalTask struct {
	SingeChildTask
	bTriggered bool
}

func (wft *WaitforSignalTask) OnEnter(a *Agent) bool {
	wft.bTriggered = false
	return wft.BehaviorTaskBase.OnEnter(a)
}

func (wft *WaitforSignalTask) OnExit(a *Agent, status EBTStatus) {
	wft.BehaviorTaskBase.OnExit(a, status)
}

func (wtf *WaitforSignalTask) Update(a *Agent, childStatus EBTStatus) EBTStatus {
	if childStatus != BT_RUNNING {
		return childStatus
	}
	if !wtf.bTriggered {
		if w, ok := wtf.GetNode().(*WaitforSignal); ok {
			wtf.bTriggered = w.CheckIfSignaled(a)
		}
	}
	if wtf.bTriggered {
		if wtf.root == nil {
			return BT_SUCCESS
		}
		return wtf.Update(a, childStatus)
	}
	return BT_RUNNING
}
