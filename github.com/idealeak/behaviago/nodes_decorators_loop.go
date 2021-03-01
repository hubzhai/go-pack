package behaviago

/**
  DecoratorLoop can be set a integer Count value. It increases inner count value when it updates.
  It always return Running until inner count less equal than integer Count value. Or returns the child
  value. It always return Running when the count limit equal to -1.
*/
const (
	DecoratorLoopNodeName = "DecoratorLoop"
)

func init() {
	RegisteNodeCreator(DecoratorLoopNodeName, BehaviorNodeCreatorWrapper(func() BehaviorNode {
		n := newDecoratorLoop()
		return n
	}))
}

func newDecoratorLoop() *DecoratorLoop {
	n := &DecoratorLoop{}
	n.SetClassName(DecoratorLoopNodeName)
	n.SetSelf(n)
	return n
}

// ============================================================================
type DecoratorLoop struct {
	DecoratorCount
	bDoneWithinFrame bool
}

func (d *DecoratorLoop) Load(version int, agentType string, properties []property_t) {
	d.DecoratorCount.Load(version, agentType, properties)
	for i := 0; i < len(properties); i++ {
		if properties[i].name == "DoneWithinFrame" {
			if properties[i].value == "true" {
				d.bDoneWithinFrame = true
			}
		}
	}
}

func (d *DecoratorLoop) Decompose(node BehaviorNode, seqTask *PlannerTaskComplex, depth int, planner *Planner) bool {
	BTGLog.Trace("(d *DecoratorLoop) Decompose enter")
	if dl, ok := node.(*DecoratorLoop); ok {
		childs := dl.GetChilds()
		if len(childs) > 0 {
			c := childs[0]
			childTask := planner.DecomposeNode(c, depth)
			if childTask == nil {
				return false
			}
			seqTask.AddChild(childTask)
			return true
		}
	}

	return false
}

func (d *DecoratorLoop) CreateTask() BehaviorTask {
	BTGLog.Trace("(d *DecoratorLoop) CreateTask()")
	return &DecoratorLoopTask{}
}

func (n *DecoratorLoop) IsValid(a *Agent, task BehaviorTask) bool {
	if _, ok := task.GetNode().(*DecoratorLoop); !ok {
		return false
	}
	return true
}

// ============================================================================
///Returns BT_FAILURE for the specified number of iterations, then returns BT_SUCCESS after that
type DecoratorLoopTask struct {
	DecoratorCountTask
}

func (dt *DecoratorLoopTask) Decorate(status EBTStatus) EBTStatus {
	BTGLog.Tracef("(d *DecoratorLoopTask) Decorate enter(count=%v)", dt.count)
	if dt.count > 0 {
		dt.count--
		if dt.count == 0 {
			return BT_SUCCESS
		}
		return BT_RUNNING
	}
	if dt.count == -1 {
		return BT_RUNNING
	}
	return BT_SUCCESS
}

func (dt *DecoratorLoopTask) Update(a *Agent, childStatus EBTStatus) EBTStatus {
	BTGLog.Trace("(d *DecoratorLoopTask) Update enter")
	if node, ok := dt.GetNode().(*DecoratorLoop); ok {
		if node.bDoneWithinFrame && dt.root != nil && dt.count > 0 {
			status := BT_INVALID
			for i := int64(0); i < dt.count; i++ {
				status = dt.root.Exec(a, childStatus)
				if node.IsDecorateWhenChildEnds() {
					for status == BT_RUNNING {
						status = dt.DecoratorCountTask.Update(a, childStatus)
					}
				}
				if status == BT_FAILURE {
					return BT_FAILURE
				}
			}
			return BT_SUCCESS
		}
	}
	return dt.DecoratorCountTask.Update(a, childStatus)
}
