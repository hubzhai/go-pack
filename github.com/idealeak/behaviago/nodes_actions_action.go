package behaviago

///An action is a member function of agent
/**
  Action node is the bridge between behavior tree and agent member function.
  a member function can be assigned to an action node. function will be
  invoked when Action node ticked. function attached can be up to eight parameters most.
*/
const (
	ActionNodeName = "Action"
)

func init() {
	RegisteNodeCreator(ActionNodeName, BehaviorNodeCreatorWrapper(func() BehaviorNode {
		actNode := newActionNode()
		return actNode
	}))
}

func newActionNode() *Action {
	actNode := &Action{}
	actNode.SetClassName(ActionNodeName)
	actNode.SetSelf(actNode)
	return actNode
}

// ============================================================================
type Action struct {
	BehaviorNodeBase
	resultOption  EBTStatus
	method        *Method
	resultFunctor *Method
}

func (act *Action) Load(version int, agentType string, properties []property_t) {
	act.BehaviorNodeBase.Load(version, agentType, properties)

	for i := 0; i < len(properties); i++ {
		switch properties[i].name {
		case "Method":
			act.method = LoadMethod(properties[i].value)
		case "ResultOption":
			switch properties[i].value {
			case "BT_INVALID":
				act.resultOption = BT_INVALID
			case "BT_FAILURE":
				act.resultOption = BT_FAILURE
			case "BT_RUNNING":
				act.resultOption = BT_RUNNING
			default:
				act.resultOption = BT_SUCCESS
			}
		case "ResultFunctor":
			act.resultFunctor = LoadMethod(properties[i].value)
		}
	}
}

func (act *Action) Execute(a *Agent, childStatus EBTStatus) EBTStatus {
	BTGLog.Tracef("(act *Action) Execute enter")
	result := BT_SUCCESS
	if act.method != nil {
		rets := a.InvokeMethod(act.method)
		if act.resultOption != BT_INVALID {
			result = act.resultOption
		} else if act.resultFunctor != nil {
			rets = a.InvokeMethod(act.resultFunctor)
			result = CheckResult(rets)
		} else {
			result = CheckResult(rets)
		}
	} else {
		result = act.UpdateImpl(a, childStatus)
	}
	return result
}

func (act *Action) CreateTask() BehaviorTask {
	BTGLog.Trace("(act *Action) CreateTask()")
	return &ActionTask{}
}

func (act *Action) IsValid(a *Agent, task BehaviorTask) bool {
	if _, ok := task.GetNode().(*Action); !ok {
		return false
	}
	return true
}

// ============================================================================
type ActionTask struct {
	LeafTask
}

func (at *ActionTask) OnEnter(a *Agent) bool {
	BTGLog.Tracef("(at *ActionTask) OnEnter enter")
	return true
}

func (at *ActionTask) OnExit(a *Agent, s EBTStatus) {
	BTGLog.Tracef("(at *ActionTask) OnExit enter")
	at.LeafTask.OnExit(a, s)
}

func (at *ActionTask) Update(a *Agent, s EBTStatus) EBTStatus {
	BTGLog.Tracef("(at *ActionTask) Update enter")
	n := at.GetNode()
	if n != nil {
		return n.Execute(a, s)
	}
	return BT_FAILURE
}
