package behaviago

/**
  Nothing to do, just return success.
*/
const (
	NoopNodeName = "Noop"
)

func init() {
	RegisteNodeCreator(NoopNodeName, BehaviorNodeCreatorWrapper(func() BehaviorNode {
		n := newNoop()
		return n
	}))
}

func newNoop() *Noop {
	n := &Noop{}
	n.SetClassName(NoopNodeName)
	n.SetSelf(n)
	return n
}

// ============================================================================
type Noop struct {
	BehaviorNodeBase
}

func (n *Noop) Load(version int, agentType string, properties []property_t) {
	n.BehaviorNodeBase.Load(version, agentType, properties)
}

func (n *Noop) CreateTask() BehaviorTask {
	BTGLog.Trace("(n *Noop) CreateTask()")
	return &NoopTask{}
}

func (n *Noop) IsValid(a *Agent, task BehaviorTask) bool {
	if _, ok := task.GetNode().(*Noop); !ok {
		return false
	}
	return true
}

// ============================================================================
type NoopTask struct {
	LeafTask
}

func (nt *NoopTask) OnEnter(a *Agent) bool {
	return nt.BehaviorTaskBase.OnEnter(a)
}

func (nt *NoopTask) OnExit(a *Agent, status EBTStatus) {
	nt.BehaviorTaskBase.OnExit(a, status)
}

func (nt *NoopTask) Update(a *Agent, childStatus EBTStatus) EBTStatus {
	return nt.Update(a, childStatus)
}
