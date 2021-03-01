package behaviago

/**
  WithPrecondition is the precondition of SelectorLoop child. must be used in conjunction with SelectorLoop.
  WithPrecondition can return SUCCESS or FAILURE. child would execute when it returns SUCCESS, or not.
*/
const (
	WithPreconditionNodeName = "WithPrecondition"
)

func init() {
	RegisteNodeCreator(WithPreconditionNodeName, BehaviorNodeCreatorWrapper(func() BehaviorNode {
		w := &WithPrecondition{}
		return w
	}))
}

func newWithPreconditionNode() *WithPrecondition {
	w := &WithPrecondition{}
	w.SetClassName(WithPreconditionNodeName)
	w.SetSelf(w)
	return w
}

// ============================================================================
type WithPrecondition struct {
	BehaviorNodeBase
}

func (w *WithPrecondition) Load(version int, agentType string, properties []property_t) {
	w.BehaviorNodeBase.Load(version, agentType, properties)
}

func (w *WithPrecondition) CreateTask() BehaviorTask {
	BTGLog.Trace("(w *WithPrecondition) CreateTask()")
	return &WithPreconditionTask{}
}

func (w *WithPrecondition) IsValid(a *Agent, task BehaviorTask) bool {
	if _, ok := task.GetNode().(*WithPrecondition); !ok {
		return false
	}
	return true
}

// ============================================================================
type WithPreconditionTask struct {
	SequenceTask
}

func (wt *WithPreconditionTask) PreconditionNode() BehaviorTask {
	if len(wt.childs) == 2 {
		return wt.childs[0]
	}
	return nil
}

func (wt *WithPreconditionTask) ActionNode() BehaviorTask {
	if len(wt.childs) == 2 {
		return wt.childs[1]
	}
	return nil
}

func (wt *WithPreconditionTask) OnEnter(a *Agent) bool {
	return wt.BehaviorTaskBase.OnEnter(a)
}

func (wt *WithPreconditionTask) OnExit(a *Agent, status EBTStatus) {
	wt.BehaviorTaskBase.OnExit(a, status)
}

func (wt *WithPreconditionTask) Update(a *Agent, childStatus EBTStatus) EBTStatus {
	return BT_RUNNING
}
