package behaviago

/**
  Boolean arithmetical operation &&
*/
const (
	AndNodeName = "And"
)

func init() {
	RegisteNodeCreator(AndNodeName, BehaviorNodeCreatorWrapper(func() BehaviorNode {
		a := &And{}
		return a
	}))
}

func newAndNode() *And {
	a := &And{}
	a.SetClassName(AndNodeName)
	a.SetSelf(a)
	return a
}

// ============================================================================
type And struct {
	BehaviorNodeBase
}

func (a *And) Load(version int, agentType string, properties []property_t) {
	a.BehaviorNodeBase.Load(version, agentType, properties)
}

func (a *And) CreateTask() BehaviorTask {
	BTGLog.Trace("(a *And) CreateTask()")
	return &AndTask{}
}

func (a *And) IsValid(aa *Agent, task BehaviorTask) bool {
	if _, ok := task.GetNode().(*And); !ok {
		return false
	}
	return true
}

func (a *And) Evaluate(aa *Agent, result EBTStatus) bool {
	for _, n := range a.childs {
		if !n.Evaluate(aa, result) {
			return false
		}
	}
	return true
}

// ============================================================================
type AndTask struct {
	SequenceTask
}

func (at *AndTask) OnEnter(a *Agent) bool {
	return at.BehaviorTaskBase.OnEnter(a)
}

func (at *AndTask) OnExit(a *Agent, status EBTStatus) {
	at.BehaviorTaskBase.OnExit(a, status)
}

func (at *AndTask) Update(a *Agent, childStatus EBTStatus) EBTStatus {
	for _, c := range at.childs {
		if c.Exec(a, childStatus) == BT_FAILURE {
			return BT_FAILURE
		}
	}
	return BT_SUCCESS
}
