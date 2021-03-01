package behaviago

/**
  DecoratorLoopUntil can be set two conditions, loop count 'C' and a return value 'R'.
  if current update count less equal than 'C' and child return value not equal to 'R',
  it returns Running. Or returns child value.
*/
const (
	DecoratorLoopUntilNodeName = "DecoratorLoopUntil"
)

func init() {
	RegisteNodeCreator(DecoratorLoopUntilNodeName, BehaviorNodeCreatorWrapper(func() BehaviorNode {
		n := newDecoratorLoopUntilNode()
		return n
	}))
}

func newDecoratorLoopUntilNode() *DecoratorLoopUntil {
	n := &DecoratorLoopUntil{}
	n.SetClassName(DecoratorLoopUntilNodeName)
	n.SetSelf(n)
	return n
}

// ============================================================================
type DecoratorLoopUntil struct {
	DecoratorCount
	bUntil bool
}

func (d *DecoratorLoopUntil) Load(version int, agentType string, properties []property_t) {
	d.DecoratorCount.Load(version, agentType, properties)
	for i := 0; i < len(properties); i++ {
		if properties[i].name == "Until" {
			if properties[i].value == "true" {
				d.bUntil = true
			}
		}
	}
}

func (d *DecoratorLoopUntil) CreateTask() BehaviorTask {
	BTGLog.Trace("(d *DecoratorLoopUntil) CreateTask()")
	return &DecoratorLoopUntilTask{}
}

func (n *DecoratorLoopUntil) IsValid(a *Agent, task BehaviorTask) bool {
	if _, ok := task.GetNode().(*DecoratorLoopUntil); !ok {
		return false
	}
	return true
}

// ============================================================================
///Returns BT_RUNNING until the child returns BT_SUCCESS. if the child returns BT_FAILURE, it still returns BT_RUNNING
/**
  however, if bUntil is false, the checking condition is inverted.
  i.e. it Returns BT_RUNNING until the child returns BT_FAILURE. if the child returns BT_SUCCESS, it still returns BT_RUNNING
*/
type DecoratorLoopUntilTask struct {
	DecoratorCountTask
}

func (dt *DecoratorLoopUntilTask) Decorate(status EBTStatus) EBTStatus {
	if dt.count > 0 {
		dt.count--
	}
	if dt.count == 0 {
		return BT_SUCCESS
	}

	if node, ok := dt.GetNode().(*DecoratorLoopUntil); ok {
		if node.bUntil {
			if status == BT_SUCCESS {
				return BT_SUCCESS
			}
		} else {
			if status == BT_FAILURE {
				return BT_FAILURE
			}
		}
	}
	return BT_RUNNING
}
