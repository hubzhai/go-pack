package behaviago

/**
  Output message specified when it updates.
*/
const (
	DecoratorLogNodeName = "DecoratorLog"
)

func init() {
	RegisteNodeCreator(DecoratorLogNodeName, BehaviorNodeCreatorWrapper(func() BehaviorNode {
		n := newDecoratorLogNode()
		return n
	}))
}

func newDecoratorLogNode() *DecoratorLog {
	n := &DecoratorLog{}
	n.SetClassName(DecoratorLogNodeName)
	n.SetSelf(n)
	return n
}

// ============================================================================

type DecoratorLog struct {
	DecoratorNode
	message string
}

func (d *DecoratorLog) Load(version int, agentType string, properties []property_t) {
	d.DecoratorNode.Load(version, agentType, properties)
	for i := 0; i < len(properties); i++ {
		if properties[i].name == "Log" {
			d.message = properties[i].value
		}
	}
}

func (d *DecoratorLog) CreateTask() BehaviorTask {
	BTGLog.Trace("(d *DecoratorLog) CreateTask()")
	return &DecoratorLogTask{}
}

func (n *DecoratorLog) IsValid(a *Agent, task BehaviorTask) bool {
	if _, ok := task.GetNode().(*DecoratorLog); !ok {
		return false
	}
	return true
}

// ============================================================================
type DecoratorLogTask struct {
	DecoratorTask
}

func (dt *DecoratorLogTask) Decorate(status EBTStatus) EBTStatus {
	if node, ok := dt.GetNode().(*DecoratorLog); ok {
		BTGLog.Infof("DecoratorLogTask:%v", node.message)
	}
	return status
}
