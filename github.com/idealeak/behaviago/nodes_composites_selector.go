package behaviago

///Execute behaviors from first to last
/**
  Selectors tick each of their children one at a time from top to bottom. If a child returns
  Success, then so does the Selector. If it returns Failure, the Selector will move on to the
  next child in line and return Running.If a child returns Running, so does the Selector and
  that same child will be ticked again next time the Selector is ticked. Once the Selector
  reaches the end of its child list, it returns Failure and resets its child index – meaning
  the first child in the line will be ticked on the next tick of the Selector.
*/
const (
	SelectorNodeName = "Selector"
)

func init() {
	RegisteNodeCreator(SelectorNodeName, BehaviorNodeCreatorWrapper(func() BehaviorNode {
		s := newSelectorNode()
		return s
	}))
}

func newSelectorNode() *Selector {
	s := &Selector{}
	s.SetClassName(SelectorNodeName)
	s.SetSelf(s)
	return s
}

// ============================================================================
type Selector struct {
	BehaviorNodeBase
}

func (s *Selector) Load(version int, agentType string, properties []property_t) {
	s.BehaviorNodeBase.Load(version, agentType, properties)
}

func (s *Selector) CreateTask() BehaviorTask {
	BTGLog.Trace("(s *Selector) CreateTask()")
	return NewSelectorTask()
}

func (s *Selector) IsValid(a *Agent, task BehaviorTask) bool {
	if _, ok := task.GetNode().(*Selector); !ok {
		return false
	}
	return true
}

func (s *Selector) Evaluate(a *Agent, result EBTStatus) bool {
	for _, c := range s.childs {
		if c.Evaluate(a, result) {
			return true
		}
	}
	return false
}

func (s *Selector) CheckIfInterrupted(a *Agent) bool {
	return s.EvaluteCustomCondition(a)
}

func (s *Selector) Decompose(node BehaviorNode, seqTask *PlannerTaskComplex, depth int, planner *Planner) bool {
	childs := node.GetChilds()
	for i := 0; i < len(childs); i++ {
		c := childs[i]
		if c != nil {
			childTask := planner.DecomposeNode(c, depth)
			if childTask != nil {
				seqTask.AddChild(childTask)
				return true
			}
		}
	}
	return false
}

func (s *Selector) SelectorUpdate(a *Agent, childStatus EBTStatus, activeChildIndex int, childs []BehaviorTask) (int, EBTStatus) {
	status := childStatus
	childSize := len(childs)
	for {
		if activeChildIndex >= childSize || activeChildIndex < 0 {
			return activeChildIndex, BT_FAILURE
		}

		if status == BT_RUNNING {
			if s.CheckIfInterrupted(a) {
				return activeChildIndex, BT_FAILURE
			}
			c := childs[activeChildIndex]
			if c != nil {
				status = c.Exec(a, status)
			}
		}

		// If the child fails, or keeps running, do the same.
		if status != BT_FAILURE {
			return activeChildIndex, status
		}

		// Hit the end of the array, job done!
		activeChildIndex++
		if activeChildIndex >= childSize {
			return activeChildIndex, BT_FAILURE
		}
		status = BT_RUNNING
	}
	return activeChildIndex, status
}

// ============================================================================
type SelectorTask struct {
	CompositeTask
}

func NewSelectorTask() *SelectorTask {
	return &SelectorTask{CompositeTask: CompositeTask{activeChildIndex: -1}}
}
func (st *SelectorTask) OnEnter(a *Agent) bool {
	st.activeChildIndex = 0
	return st.BehaviorTaskBase.OnEnter(a)
}

func (st *SelectorTask) OnExit(a *Agent, status EBTStatus) {
	st.BehaviorTaskBase.OnExit(a, status)
}

func (st *SelectorTask) Update(a *Agent, childStatus EBTStatus) EBTStatus {
	if st.activeChildIndex >= 0 && st.activeChildIndex < len(st.childs) {
		if s, ok := st.GetNode().(*Selector); ok {
			index, status := s.SelectorUpdate(a, childStatus, st.activeChildIndex, st.childs)
			st.activeChildIndex = index
			return status
		}
	}
	return BT_FAILURE
}
