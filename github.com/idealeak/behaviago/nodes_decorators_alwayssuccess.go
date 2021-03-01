package behaviago

/**
  No matter what child return. DecoratorAlwaysSuccess always return Success. it can only has one child node.
*/
const (
	DecoratorAlwaysSuccessNodeName = "DecoratorAlwaysSuccess"
)

func init() {
	RegisteNodeCreator(DecoratorAlwaysSuccessNodeName, BehaviorNodeCreatorWrapper(func() BehaviorNode {
		n := newDecoratorAlwaysSuccessNode()
		return n
	}))
}

func newDecoratorAlwaysSuccessNode() *DecoratorAlwaysSuccess {
	n := &DecoratorAlwaysSuccess{}
	n.SetClassName(DecoratorAlwaysSuccessNodeName)
	n.SetSelf(n)
	return n
}

// ============================================================================
type DecoratorAlwaysSuccess struct {
	DecoratorNode
}

func (d *DecoratorAlwaysSuccess) Load(version int, agentType string, properties []property_t) {
	d.DecoratorNode.Load(version, agentType, properties)
}

func (d *DecoratorAlwaysSuccess) CreateTask() BehaviorTask {
	BTGLog.Trace("(d *DecoratorAlwaysSuccess) CreateTask()")
	return &DecoratorAlwaysSuccessTask{}
}

func (n *DecoratorAlwaysSuccess) IsValid(a *Agent, task BehaviorTask) bool {
	if _, ok := task.GetNode().(*DecoratorAlwaysSuccess); !ok {
		return false
	}
	return true
}

// ============================================================================
type DecoratorAlwaysSuccessTask struct {
	DecoratorTask
}

func (dt *DecoratorAlwaysSuccessTask) OnEnter(a *Agent) bool {
	return dt.DecoratorTask.OnEnter(a)
}

func (dt *DecoratorAlwaysSuccessTask) OnExit(a *Agent, status EBTStatus) {
	dt.DecoratorTask.OnExit(a, status)
}

func (dt *DecoratorAlwaysSuccessTask) Decorate(status EBTStatus) EBTStatus {
	return BT_SUCCESS
}
