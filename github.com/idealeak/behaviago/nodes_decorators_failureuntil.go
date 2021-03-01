package behaviago

/**
  DecoratorFailureUntil node always return Failure until it reaches a specified number of count.
  when reach time exceed the count specified return Success. If the specified number of count
  is -1, then always return failed.
*/
const (
	DecoratorFailureUntilNodeName = "DecoratorFailureUntil"
)

func init() {
	RegisteNodeCreator(DecoratorFailureUntilNodeName, BehaviorNodeCreatorWrapper(func() BehaviorNode {
		n := newDecoratorFailureUntilNode()
		return n
	}))
}

func newDecoratorFailureUntilNode() *DecoratorFailureUntil {
	n := &DecoratorFailureUntil{}
	n.SetClassName(DecoratorFailureUntilNodeName)
	n.SetSelf(n)
	return n
}

// ============================================================================
type DecoratorFailureUntil struct {
	DecoratorCount
}

func (d *DecoratorFailureUntil) CreateTask() BehaviorTask {
	BTGLog.Trace("(d *DecoratorFailureUntil) CreateTask()")
	return &DecoratorFailureUntilTask{}
}

func (n *DecoratorFailureUntil) IsValid(a *Agent, task BehaviorTask) bool {
	if _, ok := task.GetNode().(*DecoratorFailureUntil); !ok {
		return false
	}
	return true
}

// ============================================================================

type DecoratorFailureUntilTask struct {
	DecoratorCountTask
}

func (dt *DecoratorFailureUntilTask) OnReset(a *Agent) {
	dt.count = 0
}

func (dt *DecoratorFailureUntilTask) OnEnter(a *Agent) bool {
	//don't reset the m_n if it is restarted
	if dt.count == 0 {
		dt.count = dt.GetCount(a)
		if dt.count == 0 {
			return false
		}
	}
	return true
}

func (dt *DecoratorFailureUntilTask) Decorate(status EBTStatus) EBTStatus {
	if dt.count > 0 {
		dt.count--
		if dt.count == 0 {
			return BT_SUCCESS
		}
		return BT_FAILURE
	}
	if dt.count == -1 {
		return BT_FAILURE
	}
	return BT_SUCCESS
}
