package behaviago

/**
  No matter what child return. DecoratorAlwaysFailure always return Failure. it can only has one child node.
*/
const (
	DecoratorAlwaysFailureNodeName = "DecoratorAlwaysFailure"
)

func init() {
	RegisteNodeCreator(DecoratorAlwaysFailureNodeName, BehaviorNodeCreatorWrapper(func() BehaviorNode {
		n := newDecoratorAlwaysFailureNode()
		return n
	}))
}

func newDecoratorAlwaysFailureNode() *DecoratorAlwaysFailure {
	n := &DecoratorAlwaysFailure{}
	n.SetClassName(DecoratorAlwaysFailureNodeName)
	n.SetSelf(n)
	return n
}

// ============================================================================
type DecoratorAlwaysFailure struct {
	DecoratorNode
}

func (d *DecoratorAlwaysFailure) Load(version int, agentType string, properties []property_t) {
	d.DecoratorNode.Load(version, agentType, properties)
}

func (d *DecoratorAlwaysFailure) CreateTask() BehaviorTask {
	BTGLog.Trace("(d *DecoratorAlwaysFailure) CreateTask()")
	return &DecoratorAlwaysFailureTask{}
}

func (n *DecoratorAlwaysFailure) IsValid(a *Agent, task BehaviorTask) bool {
	if _, ok := task.GetNode().(*DecoratorAlwaysFailure); !ok {
		return false
	}
	return true
}

// ============================================================================
type DecoratorAlwaysFailureTask struct {
	DecoratorTask
}

func (dt *DecoratorAlwaysFailureTask) OnEnter(a *Agent) bool {
	return dt.DecoratorTask.OnEnter(a)
}

func (dt *DecoratorAlwaysFailureTask) OnExit(a *Agent, status EBTStatus) {
	dt.DecoratorTask.OnExit(a, status)
}

func (dt *DecoratorAlwaysFailureTask) Decorate(status EBTStatus) EBTStatus {
	return BT_FAILURE
}
