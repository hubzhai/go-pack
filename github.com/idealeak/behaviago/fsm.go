package behaviago

import "strconv"

const (
	FSMNodeName = "FSM"
)

func init() {
	RegisteNodeCreator(FSMNodeName, BehaviorNodeCreatorWrapper(func() BehaviorNode {
		f := newFSMNode()
		return f
	}))
}

func newFSMNode() *FSM {
	n := &FSM{}
	n.SetClassName(FSMNodeName)
	n.SetSelf(n)
	return n
}

// ============================================================================
type FSM struct {
	BehaviorNodeBase
	initialId int
}

func (fsm *FSM) SetInitialId(intialId int) {
	fsm.initialId = intialId
}

func (fsm *FSM) GetInitalId() int {
	return fsm.initialId
}

func (fsm *FSM) Load(version int, agentType string, properties []property_t) {
	fsm.BehaviorNodeBase.Load(version, agentType, properties)
	for i := 0; i < len(properties); i++ {
		if properties[i].name == "initialid" {
			fsm.initialId, _ = strconv.Atoi(properties[i].value)
		}
	}
}

func (fsm *FSM) CreateTask() BehaviorTask {
	BTGLog.Trace("(fsm *FSM) CreateTask()")
	return &FSMTask{}
}

func (fsm *FSM) IsValid(a *Agent, task BehaviorTask) bool {
	if _, ok := task.GetNode().(*FSM); !ok {
		return false
	}
	return true
}

// ============================================================================
type FSMTask struct {
	CompositeTask
}

func (ft *FSMTask) OnEnter(a *Agent) bool {
	if fsm, ok := ft.GetNode().(*FSM); ok {
		ft.activeChildIndex = 0
		ft.currentNodeId = fsm.GetInitalId()
		return true
	}
	return ft.CompositeTask.OnEnter(a)
}

func (ft *FSMTask) OnExit(a *Agent, status EBTStatus) {
	ft.currentNodeId = -1
	ft.CompositeTask.OnExit(a, status)
}

func (ft *FSMTask) UpdateFSM(a *Agent, childStatus EBTStatus) EBTStatus {
	state := childStatus
	bLoop := true
	for bLoop {
		currentState := ft.GetChildById(ft.currentNodeId)
		if currentState != nil {
			currentState.Exec(a, BT_RUNNING)
			if state, ok := currentState.(*StateTask); ok {
				if state.IsEndState() {
					return BT_SUCCESS
				}
			}
			nextStateId := currentState.GetNextStateId()
			if nextStateId < 0 { // don't know why, the nextStateID might be -2147483648, so change the condition
				//if not transitioned, don't go on next state, to exit
				bLoop = false
			} else {
				//if transitioned, go on next state
				ft.currentNodeId = nextStateId
			}
		}
	}
	return state
}

func (ft *FSMTask) Update(a *Agent, childStatus EBTStatus) EBTStatus {
	if ft.activeChildIndex < len(ft.childs) {
		return ft.UpdateFSM(a, childStatus)
	}
	return BT_RUNNING
}
