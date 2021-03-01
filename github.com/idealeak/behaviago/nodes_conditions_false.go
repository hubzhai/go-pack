package behaviago

/**
  False is a leaf node that always return Failure.
*/
const (
	FalseNodeName = "False"
)

func init() {
	RegisteNodeCreator(FalseNodeName, BehaviorNodeCreatorWrapper(func() BehaviorNode {
		f := newFalseNode()
		return f
	}))
}

func newFalseNode() *False {
	f := &False{}
	f.SetClassName(FalseNodeName)
	f.SetSelf(f)
	return f
}

// ============================================================================
type False struct {
	BehaviorNodeBase
}

func (f *False) CreateTask() BehaviorTask {
	BTGLog.Trace("(f *False) CreateTask()")
	return &FalseTask{}
}

func (f *False) IsValid(aa *Agent, task BehaviorTask) bool {
	if _, ok := task.GetNode().(*False); !ok {
		return false
	}
	return f.BehaviorNodeBase.IsValid(aa, task)
}

// ============================================================================
type FalseTask struct {
	LeafTask
}

func (ft *FalseTask) Update(a *Agent, childStatus EBTStatus) EBTStatus {
	return BT_FAILURE
}
