package behaviago

/**
  UntilFailureUntil node always return Success until it reaches a specified number of count.
  when reach time exceed the count specified return Failure. If the specified number of count
  is -1, then always return Success.
*/

const (
	DecoratorSuccessUntilNodeName = "DecoratorSuccessUntil"
)

func init() {
	RegisteNodeCreator(DecoratorSuccessUntilNodeName, BehaviorNodeCreatorWrapper(func() BehaviorNode {
		n := newDecoratorSuccessUntilNode()
		return n
	}))
}

func newDecoratorSuccessUntilNode() *DecoratorSuccessUntil {
	n := &DecoratorSuccessUntil{}
	n.SetClassName(DecoratorSuccessUntilNodeName)
	n.SetSelf(n)
	return n
}

// ============================================================================
type DecoratorSuccessUntil struct {
	DecoratorCount
}

func (d *DecoratorSuccessUntil) CreateTask() BehaviorTask {
	BTGLog.Trace("(d *DecoratorSuccessUntil) CreateTask()")
	return &DecoratorSuccessUntilTask{}
}

func (n *DecoratorSuccessUntil) IsValid(a *Agent, task BehaviorTask) bool {
	if _, ok := task.GetNode().(*DecoratorSuccessUntil); !ok {
		return false
	}
	return true
}

// ============================================================================
///Returns BT_SUCCESS for the specified number of iterations, then returns BT_FAILURE after that
type DecoratorSuccessUntilTask struct {
	DecoratorCountTask
}

func (dt *DecoratorSuccessUntilTask) OnEnter(a *Agent) bool {
	//don't reset the count if it is restarted
	if dt.count == 0 {
		count := dt.GetCount(a)
		if count == 0 {
			return false
		}
		dt.count = count
	}
	return true
}

func (dt *DecoratorSuccessUntilTask) Decorate(status EBTStatus) EBTStatus {
	if dt.count > 0 {
		dt.count--
		if dt.count == 0 {
			return BT_FAILURE
		}
		return BT_SUCCESS
	}
	if dt.count == -1 {
		return BT_SUCCESS
	}
	return BT_FAILURE
}
