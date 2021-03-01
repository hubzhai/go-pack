package behaviago

import "strings"

///Assignment
/**
  Assign a right value to left par or agent property. a right value can be a par or agent property.
*/
const (
	AssignmentNodeName = "Assignment"
)

func init() {
	RegisteNodeCreator(AssignmentNodeName, BehaviorNodeCreatorWrapper(func() BehaviorNode {
		a := newAssignmentNode()
		return a
	}))
}

func newAssignmentNode() *Assignment {
	a := &Assignment{}
	a.SetClassName(AssignmentNodeName)
	a.SetSelf(a)
	return a
}

// ============================================================================
type Assignment struct {
	BehaviorNodeBase
	opl  *Property
	opr  *Property
	oprm *Method
}

func (at *Assignment) Load(version int, agentType string, properties []property_t) {
	at.BehaviorNodeBase.Load(version, agentType, properties)
	for i := 0; i < len(properties); i++ {
		switch properties[i].name {
		case "Opl":
			at.opl = ParseProperty(properties[i].value)
		case "Opr":
			if !strings.ContainsAny(properties[i].value, "(") {
				at.opr = ParseProperty(properties[i].value)
			} else {
				at.oprm = LoadMethod(properties[i].value)
			}
		}
	}
}

func (at *Assignment) CreateTask() BehaviorTask {
	BTGLog.Trace("(at *Assignment) CreateTask()")
	return &AssignmentTask{}
}

func (at *Assignment) IsValid(a *Agent, task BehaviorTask) bool {
	if _, ok := task.GetNode().(*Assignment); !ok {
		return false
	}
	return true
}

// ============================================================================
type AssignmentTask struct {
	LeafTask
}

func (at *AssignmentTask) OnEnter(a *Agent) bool {
	return at.LeafTask.OnEnter(a)
}

func (at *AssignmentTask) OnExit(a *Agent, status EBTStatus) {
	at.LeafTask.OnExit(a, status)
}

func (at *AssignmentTask) Update(a *Agent, childStatus EBTStatus) EBTStatus {
	if node, ok := at.GetNode().(*Assignment); ok {
		if EvaluteAssignment(a, node.opl, node.opr, node.oprm) {
			return BT_SUCCESS
		} else {
			return node.UpdateImpl(a, childStatus)
		}
	}
	return BT_SUCCESS
}

// ============================================================================
func EvaluteAssignment(a *Agent, opl *Property, opr *Property, oprm *Method) bool {
	if opl != nil {
		if oprm != nil {
			al := opl.GetParentAgent(a)
			if al != nil {
				return opl.SetFromByMethod(a, oprm, al)
			}
			return true
		} else if opr != nil {
			al := opl.GetParentAgent(a)
			ar := opr.GetParentAgent(a)
			if al != nil && ar != nil {
				return opl.SetFromByProperty(ar, opr, al)
			}
			return true
		}
	}
	return false
}
