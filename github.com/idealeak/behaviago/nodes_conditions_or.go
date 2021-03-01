package behaviago

/**
  Boolean arithmetical operation &&
*/
const (
	OrNodeName = "Or"
)

func init() {
	RegisteNodeCreator(OrNodeName, BehaviorNodeCreatorWrapper(func() BehaviorNode {
		o := newOrNode()
		return o
	}))
}

func newOrNode() *Or {
	o := &Or{}
	o.SetClassName(OrNodeName)
	o.SetSelf(o)
	return o
}

// ============================================================================
type Or struct {
	BehaviorNodeBase
}

func (o *Or) Load(version int, agentType string, properties []property_t) {
	o.BehaviorNodeBase.Load(version, agentType, properties)
}

func (o *Or) CreateTask() BehaviorTask {
	BTGLog.Trace("(o *Or) CreateTask()")
	return &OrTask{}
}

func (o *Or) IsValid(aa *Agent, task BehaviorTask) bool {
	if _, ok := task.GetNode().(*Or); !ok {
		return false
	}
	return true
}

func (o *Or) Evaluate(aa *Agent, result EBTStatus) bool {
	for _, n := range o.childs {
		if n.Evaluate(aa, result) {
			return true
		}
	}
	return false
}

// ============================================================================
type OrTask struct {
	SelectorTask
}

func (ot *OrTask) OnEnter(a *Agent) bool {
	return ot.BehaviorTaskBase.OnEnter(a)
}

func (ot *OrTask) OnExit(a *Agent, status EBTStatus) {
	ot.BehaviorTaskBase.OnExit(a, status)
}

func (ot *OrTask) Update(a *Agent, childStatus EBTStatus) EBTStatus {
	for _, c := range ot.childs {
		if c.Exec(a, childStatus) == BT_SUCCESS {
			return BT_SUCCESS
		}
	}
	return BT_FAILURE
}
