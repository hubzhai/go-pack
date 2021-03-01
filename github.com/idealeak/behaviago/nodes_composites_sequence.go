package behaviago

///Execute behaviors from first to last
/**
  Sequences tick each of their children one at a time from top to bottom. If a child returns Failure,
  so does the Sequence. If it returns Success, the Sequence will move on to the next child in line
  and return Running.If a child returns Running, so does the Sequence and that same child will be
  ticked again next time the Sequence is ticked.Once the Sequence reaches the end of its child list,
  it returns Success and resets its child index meaning the first child in the line will be ticked
  on the next tick of the Sequence.
*/
const (
	SequenceNodeName = "Sequence"
)

func init() {
	RegisteNodeCreator(SequenceNodeName, BehaviorNodeCreatorWrapper(func() BehaviorNode {
		s := newSequenceNode()
		return s
	}))
}

func newSequenceNode() *Sequence {
	s := &Sequence{}
	s.SetClassName(SequenceNodeName)
	s.SetSelf(s)
	return s
}

// ============================================================================
type Sequence struct {
	BehaviorNodeBase
}

func (s *Sequence) Load(version int, agentType string, properties []property_t) {
	s.BehaviorNodeBase.Load(version, agentType, properties)
}

func (s *Sequence) CreateTask() BehaviorTask {
	BTGLog.Trace("(s *Sequence) CreateTask()")
	return NewSequenceTask()
}

func (s *Sequence) IsValid(a *Agent, task BehaviorTask) bool {
	if _, ok := task.GetNode().(*Sequence); !ok {
		return false
	}
	return true
}

func (s *Sequence) Evaluate(a *Agent, result EBTStatus) bool {
	BTGLog.Trace("(s *Sequence) Evaluate() enter")
	for _, c := range s.childs {
		if !c.Evaluate(a, result) {
			return false
		}
	}
	return true
}

func (s *Sequence) CheckIfInterrupted(a *Agent) bool {
	BTGLog.Trace("(s *Sequence) CheckIfInterrupted() enter")
	return s.EvaluteCustomCondition(a)
}

func (s *Sequence) Decompose(node BehaviorNode, seqTask *PlannerTaskComplex, depth int, planner *Planner) bool {
	BTGLog.Trace("(s *Sequence) Decompose() enter")
	childs := node.GetChilds()
	for i := 0; i < len(childs); i++ {
		c := childs[i]
		if c == nil {
			return false
		}

		childTask := planner.DecomposeNode(c, depth)
		if childTask == nil {
			return false
		}
		seqTask.AddChild(childTask)
	}
	return true
}

func (s *Sequence) SequenceUpdate(a *Agent, childStatus EBTStatus, activeChildIndex int, childs []BehaviorTask) (int, EBTStatus) {
	BTGLog.Trace("(s *Sequence) SequenceUpdate() enter")
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
		if status != BT_SUCCESS {
			return activeChildIndex, status
		}

		// Hit the end of the array, job done!
		activeChildIndex++
		if activeChildIndex >= childSize {
			return activeChildIndex, BT_SUCCESS
		}
		status = BT_RUNNING
	}
	return activeChildIndex, status
}

// ============================================================================
type SequenceTask struct {
	CompositeTask
}

func NewSequenceTask() *SequenceTask {
	return &SequenceTask{CompositeTask: CompositeTask{activeChildIndex: -1}}
}

func (st *SequenceTask) OnEnter(a *Agent) bool {
	BTGLog.Trace("(st *SequenceTask) OnEnter() enter")
	st.activeChildIndex = 0
	return true
}

func (st *SequenceTask) OnExit(a *Agent, status EBTStatus) {
	BTGLog.Trace("(st *SequenceTask) OnExit() enter")
	st.CompositeTask.OnExit(a, status)
}

func (st *SequenceTask) Update(a *Agent, childStatus EBTStatus) EBTStatus {
	BTGLog.Trace("(st *SequenceTask) Update() enter")
	if st.activeChildIndex >= 0 && st.activeChildIndex < len(st.childs) {
		if s, ok := st.GetNode().(*Sequence); ok {
			index, status := s.SequenceUpdate(a, childStatus, st.activeChildIndex, st.childs)
			st.activeChildIndex = index
			return status
		}
	}
	return BT_FAILURE
}
