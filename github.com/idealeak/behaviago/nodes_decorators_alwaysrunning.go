package behaviago

/**
  No matter what child return. DecoratorAlwaysRunning always return Running. it can only has one child node.
*/
const (
	DecoratorAlwaysRunningNodeName = "DecoratorAlwaysRunning"
)

func init() {
	RegisteNodeCreator(DecoratorAlwaysRunningNodeName, BehaviorNodeCreatorWrapper(func() BehaviorNode {
		n := newDecoratorAlwaysRunningNode()
		return n
	}))
}

func newDecoratorAlwaysRunningNode() *DecoratorAlwaysRunning {
	n := &DecoratorAlwaysRunning{}
	n.SetClassName(DecoratorAlwaysRunningNodeName)
	n.SetSelf(n)
	return n
}

// ============================================================================
type DecoratorAlwaysRunning struct {
	DecoratorNode
}

func (d *DecoratorAlwaysRunning) Load(version int, agentType string, properties []property_t) {
	d.DecoratorNode.Load(version, agentType, properties)
}

func (d *DecoratorAlwaysRunning) CreateTask() BehaviorTask {
	BTGLog.Trace("(d *DecoratorAlwaysRunning) CreateTask()")
	return &DecoratorAlwaysRunningTask{}
}

func (n *DecoratorAlwaysRunning) IsValid(a *Agent, task BehaviorTask) bool {
	if _, ok := task.GetNode().(*DecoratorAlwaysRunning); !ok {
		return false
	}
	return true
}

// ============================================================================
type DecoratorAlwaysRunningTask struct {
	DecoratorTask
}

func (dt *DecoratorAlwaysRunningTask) OnEnter(a *Agent) bool {
	return dt.DecoratorTask.OnEnter(a)
}

func (dt *DecoratorAlwaysRunningTask) OnExit(a *Agent, status EBTStatus) {
	dt.DecoratorTask.OnExit(a, status)
}

func (dt *DecoratorAlwaysRunningTask) Decorate(status EBTStatus) EBTStatus {
	return BT_RUNNING
}
