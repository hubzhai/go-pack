package behaviago

///Execute behaviors in a random order
/**
  SequenceStochastic tick each of their children in a random order. If a child returns Failure,
  so does the Sequence. If it returns Success, the Sequence will move on to the next child in line
  and return Running.If a child returns Running, so does the Sequence and that same child will be
  ticked again next time the Sequence is ticked.Once the Sequence reaches the end of its child list,
  it returns Success and resets its child index – meaning the first child in the line will be ticked
  on the next tick of the Sequence.
*/
const (
	SequenceStochasticNodeName = "SequenceStochastic"
)

func init() {
	RegisteNodeCreator(SequenceStochasticNodeName, BehaviorNodeCreatorWrapper(func() BehaviorNode {
		s := newSequenceStochasticNode()
		return s
	}))
}

func newSequenceStochasticNode() *SequenceStochastic {
	s := &SequenceStochastic{}
	s.SetClassName(SequenceStochasticNodeName)
	s.SetSelf(s)
	return s
}

// ============================================================================
type SequenceStochastic struct {
	CompositeStochastic
}

func (s *SequenceStochastic) Load(version int, agentType string, properties []property_t) {
	s.BehaviorNodeBase.Load(version, agentType, properties)
}

func (s *SequenceStochastic) CreateTask() BehaviorTask {
	BTGLog.Trace("(s *SequenceStochastic) CreateTask()")
	return &SequenceStochasticTask{}
}

func (s *SequenceStochastic) IsValid(a *Agent, task BehaviorTask) bool {
	if _, ok := task.GetNode().(*SequenceStochastic); !ok {
		return false
	}
	return true
}

func (s *SequenceStochastic) CheckIfInterrupted(a *Agent) bool {
	return s.EvaluteCustomCondition(a)
}

// ============================================================================
type SequenceStochasticTask struct {
	CompositeStochasticTask
}

func (st *SequenceStochasticTask) OnEnter(a *Agent) bool {
	return st.CompositeStochasticTask.OnEnter(a)
}

func (st *SequenceStochasticTask) OnExit(a *Agent, status EBTStatus) {
	st.CompositeStochasticTask.OnExit(a, status)
}

func (st *SequenceStochasticTask) Update(a *Agent, childStatus EBTStatus) EBTStatus {
	if st.activeChildIndex < 0 || st.activeChildIndex >= len(st.childs) {
		return BT_FAILURE
	}
	s := childStatus
	bFirst := true
	if node, ok := st.GetNode().(*SequenceStochastic); ok {
		for {
			if !bFirst || s == BT_RUNNING {
				if node.CheckIfInterrupted(a) {
					return BT_FAILURE
				}
				childIndex := st.set[st.activeChildIndex]
				c := st.childs[childIndex]
				if c != nil {
					s = c.Exec(a, s)
				}
			}
			bFirst = false
			// If the child fails, or keeps running, do the same.
			if s != BT_SUCCESS {
				return s
			}
			// Hit the end of the array, job done!
			st.activeChildIndex++
			if st.activeChildIndex >= len(st.childs) {
				return BT_SUCCESS
			}
		}
	}

	return BT_FAILURE
}
