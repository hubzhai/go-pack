package behaviago

/**
  this node has three children: 'condition' branch, 'if' branch, 'else' branch. first, it executes
  conditon, until it returns success or failure. if it returns success, it then executes 'if' branch,
  else if it returns failure, it then executes 'else' branch.
*/
const (
	IfElseNodeName = "IfElse"
)

func init() {
	RegisteNodeCreator(IfElseNodeName, BehaviorNodeCreatorWrapper(func() BehaviorNode {
		s := newIfElseNode()
		return s
	}))
}

func newIfElseNode() *IfElse {
	s := &IfElse{}
	s.SetClassName(IfElseNodeName)
	s.SetSelf(s)
	return s
}

// ============================================================================
type IfElse struct {
	BehaviorNodeBase
}

func (ie *IfElse) Load(version int, agentType string, properties []property_t) {
	ie.BehaviorNodeBase.Load(version, agentType, properties)
}

func (ie *IfElse) CreateTask() BehaviorTask {
	BTGLog.Trace("(ie *IfElse) CreateTask()")
	return NewIfElseTask()
}

func (ie *IfElse) IsValid(a *Agent, task BehaviorTask) bool {
	if _, ok := task.GetNode().(*IfElse); !ok {
		return false
	}
	return true
}

// ============================================================================
type IfElseTask struct {
	CompositeTask
}

func NewIfElseTask() *IfElseTask {
	return &IfElseTask{CompositeTask: CompositeTask{activeChildIndex: -1}}
}

func (iet *IfElseTask) OnEnter(a *Agent) bool {
	iet.activeChildIndex = -1
	return len(iet.childs) == 3
}

func (iet *IfElseTask) OnExit(a *Agent, status EBTStatus) {
	iet.BehaviorTaskBase.OnExit(a, status)
}

func (iet *IfElseTask) Update(a *Agent, childStatus EBTStatus) EBTStatus {
	if len(iet.childs) != 3 {
		return childStatus
	}
	if iet.activeChildIndex == -1 {
		condition := iet.childs[0]
		if condition != nil {
			s := condition.Exec(a, childStatus)
			if s == BT_SUCCESS {
				iet.activeChildIndex = 1
			} else if s == BT_FAILURE {
				iet.activeChildIndex = 2
			}
		}
	}
	if iet.activeChildIndex != -1 {
		c := iet.childs[iet.activeChildIndex]
		if c != nil {
			return c.Exec(a, childStatus)
		}
	}
	return BT_RUNNING
}
