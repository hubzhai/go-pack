package behaviago

/**
  Not Node inverts the return value of child. But keeping the Running value unchanged.
*/
const (
	DecoratorNotNodeName = "DecoratorNot"
)

func init() {
	RegisteNodeCreator(DecoratorNotNodeName, BehaviorNodeCreatorWrapper(func() BehaviorNode {
		n := newDecoratorNotNode()
		return n
	}))
}

func newDecoratorNotNode() *DecoratorNot {
	n := &DecoratorNot{}
	n.SetClassName(DecoratorNotNodeName)
	n.SetSelf(n)
	return n
}

// ============================================================================
type DecoratorNot struct {
	DecoratorNode
}

func (d *DecoratorNot) CreateTask() BehaviorTask {
	BTGLog.Trace("(d *DecoratorNot) CreateTask()")
	return &DecoratorNotTask{}
}

func (n *DecoratorNot) IsValid(a *Agent, task BehaviorTask) bool {
	if _, ok := task.GetNode().(*DecoratorNot); !ok {
		return false
	}
	return true
}

func (n *DecoratorNot) Evaluate(a *Agent, result EBTStatus) bool {
	childs := n.GetChilds()
	if len(childs) != 0 {
		return !childs[0].Evaluate(a, result)
	}
	return false
}

// ============================================================================
type DecoratorNotTask struct {
	DecoratorTask
}

func (dt *DecoratorNotTask) Decorate(status EBTStatus) EBTStatus {
	if status == BT_FAILURE {
		return BT_SUCCESS
	}
	if status == BT_SUCCESS {
		return BT_FAILURE
	}
	return status
}
