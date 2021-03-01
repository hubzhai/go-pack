package behaviago

const (
	DecoratorRepeatNodeName = "DecoratorRepeat"
)

func init() {
	RegisteNodeCreator(DecoratorRepeatNodeName, BehaviorNodeCreatorWrapper(func() BehaviorNode {
		n := newDecoratorRepeatNode()
		return n
	}))
}

func newDecoratorRepeatNode() *DecoratorRepeat {
	n := &DecoratorRepeat{}
	n.SetClassName(DecoratorRepeatNodeName)
	n.SetSelf(n)
	return n
}

// ============================================================================
type DecoratorRepeat struct {
	DecoratorCount
}

func (d *DecoratorRepeat) CreateTask() BehaviorTask {
	BTGLog.Trace("(d *DecoratorRepeat) CreateTask()")
	return &DecoratorRepeatTask{}
}

func (n *DecoratorRepeat) IsValid(a *Agent, task BehaviorTask) bool {
	if _, ok := task.GetNode().(*DecoratorRepeat); !ok {
		return false
	}
	return true
}

// ============================================================================
type DecoratorRepeatTask struct {
	DecoratorCountTask
}

func (dt *DecoratorRepeatTask) Decorate(status EBTStatus) EBTStatus {
	BTGLog.Warn("DecoratorRepeatTask unsupport Decorate!!!")
	return BT_INVALID
}

func (dt *DecoratorRepeatTask) Update(a *Agent, childStatus EBTStatus) EBTStatus {
	if node, ok := dt.GetNode().(*DecoratorRepeat); ok {
		if dt.root != nil && dt.count > 0 {
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
		}
		return BT_SUCCESS
	}
	return dt.DecoratorCountTask.Update(a, childStatus)
}
