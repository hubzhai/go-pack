package behaviago

/**
  Behavives similarly to SelectorTask, i.e. executing chidren until the first successful one.
  however, in the following ticks, it constantly monitors the higher priority nodes.if any
  one's precondtion node returns success, it picks it and execute it, and before executing,
  it first cleans up the original executing one. all its children are WithPreconditionTask
  or its derivatives.
*/
const (
	SelectorLoopNodeName = "SelectorLoop"
)

func init() {
	RegisteNodeCreator(SelectorLoopNodeName, BehaviorNodeCreatorWrapper(func() BehaviorNode {
		s := newSelectorLoopNode()
		return s
	}))
}

func newSelectorLoopNode() *SelectorLoop {
	s := &SelectorLoop{}
	s.SetClassName(SelectorLoopNodeName)
	s.SetSelf(s)
	return s
}

// ============================================================================
type SelectorLoop struct {
	BehaviorNodeBase
}

func (s *SelectorLoop) Load(version int, agentType string, properties []property_t) {
	s.BehaviorNodeBase.Load(version, agentType, properties)
}

func (s *SelectorLoop) CreateTask() BehaviorTask {
	return NewSelectorLoopTask()
}

func (s *SelectorLoop) IsValid(a *Agent, task BehaviorTask) bool {
	if _, ok := task.GetNode().(*SelectorLoop); !ok {
		return false
	}
	return true
}

func (s *SelectorLoop) IsManagingChildrenAsSubTrees() bool {
	return true
}

// ============================================================================
type SelectorLoopTask struct {
	CompositeTask
}

func NewSelectorLoopTask() *SelectorLoopTask {
	return &SelectorLoopTask{CompositeTask: CompositeTask{activeChildIndex: -1}}
}

func (slt *SelectorLoopTask) AddChild(c BehaviorTask) {
	BTGLog.Tracef("SelectorLoopTask.AddChild(%v) enter", c.GetClassNameString())
	if _, ok := c.(*WithPreconditionTask); ok {
		slt.CompositeTask.AddChild(c)
	} else {
		BTGLog.Warn("SelectorLoopTask.AddChild child type must be WithPreconditionTask!!!")
	}
}

func (slt *SelectorLoopTask) OnEnter(a *Agent) bool {
	slt.activeChildIndex = -1
	return slt.BehaviorTaskBase.OnEnter(a)
}

func (slt *SelectorLoopTask) OnExit(a *Agent, status EBTStatus) {
	slt.BehaviorTaskBase.OnExit(a, status)
}

func (slt *SelectorLoopTask) Update(a *Agent, childStatus EBTStatus) EBTStatus {
	idx := -1
	if childStatus != BT_RUNNING {
		if childStatus == BT_SUCCESS {
			return BT_SUCCESS
		} else if childStatus == BT_FAILURE {
			idx = slt.activeChildIndex
		}
	}

	//checking the preconditions and take the first action tree
	index := -1
	for i := idx + 1; i < len(slt.childs); i++ {
		if c, ok := slt.childs[i].(*WithPreconditionTask); ok {
			pre := c.PreconditionNode()
			if pre != nil {
				if pre.Exec(a, childStatus) == BT_SUCCESS {
					index = i
					break
				}
			}
		}
	}

	//clean up the current ticking action tree
	if index != -1 {
		if slt.activeChildIndex != -1 && slt.activeChildIndex != index {
			if currSubTree, ok := slt.childs[slt.activeChildIndex].(*WithPreconditionTask); ok {
				act := currSubTree.ActionNode()
				if act != nil {
					currSubTree.Abort(a)
				}
			}
		}

		for i := index; i < len(slt.childs); i++ {
			if childSubTree, ok := slt.childs[i].(*WithPreconditionTask); ok {
				if i > index {
					pre := childSubTree.PreconditionNode()
					if pre != nil {
						if pre.Exec(a, childStatus) != BT_SUCCESS {
							continue
						}
					}
				}

				action := childSubTree.ActionNode()
				if action != nil {
					s := action.Exec(a, childStatus)
					if s == BT_RUNNING {
						slt.activeChildIndex = i
						childSubTree.SetStatue(BT_RUNNING)
					} else {
						childSubTree.SetStatue(s)
						if s == BT_FAILURE {
							continue
						}
					}
					return s
				}
			}
		}
	}
	return BT_FAILURE
}
