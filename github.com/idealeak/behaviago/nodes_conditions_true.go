package behaviago

/**
  True is a leaf node that always return Success.
*/
const (
	TrueNodeName = "True"
)

func init() {
	RegisteNodeCreator(TrueNodeName, BehaviorNodeCreatorWrapper(func() BehaviorNode {
		t := newTrueNode()
		return t
	}))
}

func newTrueNode() *True {
	t := &True{}
	t.SetClassName(TrueNodeName)
	t.SetSelf(t)
	return t
}

// ============================================================================
type True struct {
	BehaviorNodeBase
}

func (t *True) CreateTask() BehaviorTask {
	BTGLog.Trace("(t *True) CreateTask()")
	return &TrueTask{}
}

func (t *True) IsValid(aa *Agent, task BehaviorTask) bool {
	if _, ok := task.GetNode().(*True); !ok {
		return false
	}
	return true
}

// ============================================================================
type TrueTask struct {
	LeafTask
}

func (tt *TrueTask) Update(a *Agent, childStatus EBTStatus) EBTStatus {
	return BT_SUCCESS
}
