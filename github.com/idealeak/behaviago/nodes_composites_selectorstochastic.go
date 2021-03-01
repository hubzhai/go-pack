package behaviago

/**
  the Selector runs the children from the first sequentially until the child which returns success.
  for SelectorStochastic, the children are not sequentially selected, instead it is selected stochasticly.

  for example: the children might be [0, 1, 2, 3, 4]
  Selector always select the child by the order of 0, 1, 2, 3, 4
  while SelectorStochastic, sometime, it is [4, 2, 0, 1, 3], sometime, it is [2, 3, 0, 4, 1], etc.
*/
const (
	SelectorStochasticNodeName = "SelectorStochastic"
)

func init() {
	RegisteNodeCreator(SelectorStochasticNodeName, BehaviorNodeCreatorWrapper(func() BehaviorNode {
		s := newSelectorStochasticNode()
		return s
	}))
}

func newSelectorStochasticNode() *SelectorStochastic {
	s := &SelectorStochastic{}
	s.SetClassName(SelectorStochasticNodeName)
	s.SetSelf(s)
	return s
}

// ============================================================================
type SelectorStochastic struct {
	CompositeStochastic
}

func (s *SelectorStochastic) Load(version int, agentType string, properties []property_t) {
	s.BehaviorNodeBase.Load(version, agentType, properties)
}

func (s *SelectorStochastic) CreateTask() BehaviorTask {
	BTGLog.Trace("(s *SelectorStochastic) CreateTask()")
	return &SelectorProbabilityTask{}
}

func (s *SelectorStochastic) IsValid(a *Agent, task BehaviorTask) bool {
	if _, ok := task.GetNode().(*SelectorStochastic); !ok {
		return false
	}
	return true
}

// ============================================================================
type SelectorStochasticTask struct {
	CompositeStochasticTask
}

func (sst *SelectorStochasticTask) OnEnter(a *Agent) bool {
	return sst.CompositeStochasticTask.OnEnter(a)
}

func (sst *SelectorStochasticTask) OnExit(a *Agent, status EBTStatus) {
	sst.CompositeStochasticTask.OnExit(a, status)
}

func (sst *SelectorStochasticTask) Update(a *Agent, childStatus EBTStatus) EBTStatus {
	bFirst := true
	if sst.activeChildIndex != -1 && len(sst.set) != 0 {
		// Keep going until a child behavior says its running.
		for {
			s := childStatus
			if !bFirst || s == BT_RUNNING {
				childIdx := sst.set[sst.activeChildIndex]
				c := sst.childs[childIdx]
				if c != nil {
					s = c.Exec(a, s)
				}
				bFirst = false
				// If the child succeeds, or keeps running, do the same.
				if s != BT_FAILURE {
					return s
				}
				// Hit the end of the array, job done!
				sst.activeChildIndex++
				if sst.activeChildIndex >= len(sst.childs) {
					return BT_FAILURE
				}
			}
		}
	}
	return BT_FAILURE
}
