package behaviago

///Pick a child to Execute
/**
  Choose a child to execute based on the probability have set. then return the child execute result.
*/
const (
	SelectorProbabilityNodeName = "SelectorProbability"
)

func init() {
	RegisteNodeCreator(SelectorProbabilityNodeName, BehaviorNodeCreatorWrapper(func() BehaviorNode {
		s := newSelectorProbabilityNode()
		return s
	}))
}

func newSelectorProbabilityNode() *SelectorProbability {
	s := &SelectorProbability{}
	s.SetClassName(SelectorProbabilityNodeName)
	s.SetSelf(s)
	return s
}

// ============================================================================
type SelectorProbability struct {
	BehaviorNodeBase
	method *Method
}

func (s *SelectorProbability) Load(version int, agentType string, properties []property_t) {
	s.BehaviorNodeBase.Load(version, agentType, properties)
	for i := 0; i < len(properties); i++ {
		if properties[i].name == "RandomGenerator" {
			s.method = LoadMethod(properties[i].value)
		}
	}
}

func (s *SelectorProbability) CreateTask() BehaviorTask {
	BTGLog.Trace("(s *SelectorProbability) CreateTask()")
	return NewSelectorProbabilityTask()
}

func (s *SelectorProbability) IsValid(a *Agent, task BehaviorTask) bool {
	if _, ok := task.GetNode().(*SelectorProbability); !ok {
		return false
	}
	return s.BehaviorNodeBase.IsValid(a, task)
}

func (s *SelectorProbability) AddChild(c BehaviorNode) {
	BTGLog.Tracef("SelectorProbability.AddChild(%v) enter", c.GetClassName())
	if _, ok := c.(*DecoratorWeight); ok {
		s.BehaviorNodeBase.AddChild(c)
	} else {
		BTGLog.Warn("SelectorProbability.AddChild only DecoratorWeightTask can be children")
	}
}

///Executes behaviors randomly, based on a given set of weights.
/** The weights are not percentages, but rather simple ratios.
  For example, if there were two children with a weight of one, each would have a 50% chance of being executed.
  If another child with a weight of eight were added, the previous children would have a 10% chance of being executed, and the new child would have an 80% chance of being executed.
  This weight system is intended to facilitate the fine-tuning of behaviors.
*/
// ============================================================================
type SelectorProbabilityTask struct {
	CompositeTask
	weights  []int64
	totalSum int64
}

func NewSelectorProbabilityTask() *SelectorProbabilityTask {
	return &SelectorProbabilityTask{CompositeTask: CompositeTask{activeChildIndex: -1}}
}

func (spt *SelectorProbabilityTask) AddChild(c BehaviorTask) {
	BTGLog.Tracef("SelectorProbabilityTask.AddChild(%v) enter", c.GetClassNameString())
	if _, ok := c.(*DecoratorWeightTask); ok {
		spt.CompositeTask.AddChild(c)
	} else {
		BTGLog.Warn("SelectorProbabilityTask.AddChild child type must be DecoratorWeightTask!!!")
	}
}

func (spt *SelectorProbabilityTask) OnEnter(a *Agent) bool {
	spt.activeChildIndex = -1
	spt.totalSum = 0
	if spt.weights != nil {
		spt.weights = spt.weights[0:0]
	}
	for _, c := range spt.childs {
		if task, ok := c.(*DecoratorWeightTask); ok {
			weight := task.GetWeight(a)
			spt.weights = append(spt.weights, weight)
			spt.totalSum += weight
		}
	}
	BTGLog.Tracef("SelectorProbabilityTask.OnEnter() enter totalSum=%v", spt.totalSum)
	return len(spt.weights) == len(spt.childs)
}

func (spt *SelectorProbabilityTask) OnExit(a *Agent, status EBTStatus) {
	spt.BehaviorTaskBase.OnExit(a, status)
	spt.activeChildIndex = -1
}

func (spt *SelectorProbabilityTask) Update(a *Agent, childStatus EBTStatus) EBTStatus {
	BTGLog.Trace("SelectorProbabilityTask.Update() enter")
	if childStatus != BT_RUNNING {
		return childStatus
	}
	//check if we've already chosen a node to run
	if spt.activeChildIndex != -1 {
		c := spt.childs[spt.activeChildIndex]
		if c != nil {
			return c.Exec(a, childStatus)
		}
	}
	if node, ok := spt.GetNode().(*SelectorProbability); ok {
		chosen := int64(GetRandomValue(node.method, a) * float64(spt.totalSum))
		BTGLog.Tracef("SelectorProbabilityTask.Update() enter roll randvalue=%v", chosen)
		var sum int64
		for i := 0; i < len(spt.childs); i++ {
			sum += spt.weights[i]
			if spt.weights[i] > 0 && sum >= chosen {
				c := spt.childs[i]
				if c != nil {
					s := c.Exec(a, childStatus)
					if s == BT_RUNNING {
						spt.activeChildIndex = i
					} else {
						spt.activeChildIndex = -1
					}
					return s
				}
			}
		}
	}

	return BT_FAILURE
}
