package behaviago

const (
	StateNodeName = "State"
)

func init() {
	RegisteNodeCreator(StateNodeName, BehaviorNodeCreatorWrapper(func() BehaviorNode {
		n := newStateNode()
		return n
	}))
}

func newStateNode() *State {
	n := &State{}
	n.SetClassName(StateNodeName)
	n.SetSelf(n)
	return n
}

// ============================================================================
type State struct {
	BehaviorNodeBase
	isEndState  bool
	method      *Method
	transitions []*Transition
}

func (s *State) IsEndState() bool {
	return s.isEndState
}

func (s *State) Load(version int, agentType string, properties []property_t) {
	s.BehaviorNodeBase.Load(version, agentType, properties)
	for i := 0; i < len(properties); i++ {
		switch properties[i].name {
		case "Method":
			s.method = LoadMethod(properties[i].value)
		case "IsEndState":
			s.isEndState = properties[i].value == "true"
		}
	}
}

func (s *State) Attch(pAttachment BehaviorNode, bIsPrecondition, bIsEffector, bIsTransition bool) {
	if bIsTransition {
		if transition, ok := pAttachment.(*Transition); ok {
			s.transitions = append(s.transitions, transition)
		}
	}
	s.BehaviorNodeBase.Attach(pAttachment, bIsPrecondition, bIsEffector, bIsTransition)
}

func (s *State) CreateTask() BehaviorTask {
	BTGLog.Trace("(s *State) CreateTask()")
	return &StateTask{}
}

func (s *State) IsValid(a *Agent, task BehaviorTask) bool {
	if _, ok := task.GetNode().(*State); !ok {
		return false
	}
	return true
}

func (s *State) UpdateImpl(a *Agent, childStatus EBTStatus) EBTStatus {
	return BT_RUNNING
}

func (s *State) Execute(a *Agent, result EBTStatus) EBTStatus {
	ret := BT_RUNNING
	if s.method != nil {
		s.method.Invoke(a)
	} else {
		ret = s.self.UpdateImpl(a, BT_RUNNING)
	}
	return ret
}

//nextStateId holds the next state id if it returns running when a certain transition is satisfied
//otherwise, it returns success or failure if it ends
func (s *State) UpdateNext(a *Agent) (result EBTStatus, nextStateId int) {
	nextStateId = -1
	//when no method is specified(m_method == 0),
	//'update_impl' is used to return the configured result status for both xml/bson and c#
	result = s.self.Execute(a, BT_INVALID)
	if s.isEndState {
		result = BT_SUCCESS
	} else {
		var ok bool
		nextStateId, ok = s.UpdateTransitions(a, s, s.transitions, result)
		if ok {
			//it will transition to another state, set result as success so as it exits
			result = BT_SUCCESS
		}
	}
	return
}

func (s *State) UpdateTransitions(a *Agent, node BehaviorNode, transitions []*Transition, result EBTStatus) (nextStateId int, ok bool) {
	ok = false
	if len(transitions) != 0 {
		for _, t := range transitions {
			if t.Evaluate(a, result) {
				nextStateId = t.GetTargetStateId()
				//transition actions
				t.ApplyEffects(a, EBT_BOTH)
				ok = true
				break
			}
		}
	}
	return
}

// ============================================================================
type StateTask struct {
	LeafTask
	nextStateId int
}

func (st *StateTask) GetNextStateId() int {
	return st.nextStateId
}

func (st *StateTask) IsEndState() bool {
	if s, ok := st.GetNode().(*State); ok {
		return s.IsEndState()
	}
	return true
}

func (st *StateTask) OnEnter(a *Agent) bool {
	st.nextStateId = -1
	return true
}

func (st *StateTask) OnExit(a *Agent, status EBTStatus) {
	st.LeafTask.OnExit(a, status)
}

func (st *StateTask) Update(a *Agent, childStatus EBTStatus) EBTStatus {
	if node, ok := st.GetNode().(*State); ok {
		var result EBTStatus
		result, st.nextStateId = node.UpdateNext(a)
		return result
	}
	return BT_SUCCESS
}
