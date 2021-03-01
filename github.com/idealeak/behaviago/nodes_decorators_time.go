package behaviago

import (
	"strings"
	"time"
)

/**
  It returns Running result until it reaches the time limit specified, no matter which
  value its child return. Or return the child's value.
*/

const (
	DecoratorTimeNodeName = "DecoratorTime"
)

func init() {
	RegisteNodeCreator(DecoratorTimeNodeName, BehaviorNodeCreatorWrapper(func() BehaviorNode {
		n := newDecoratorTimeNode()
		return n
	}))
}

func newDecoratorTimeNode() *DecoratorTime {
	n := &DecoratorTime{}
	n.SetClassName(DecoratorTimeNodeName)
	n.SetSelf(n)
	return n
}

// ============================================================================
type DecoratorTime struct {
	DecoratorNode
	timeVar *Property
	timeMet *Method
}

func (d *DecoratorTime) Load(version int, agentType string, properties []property_t) {
	d.DecoratorNode.Load(version, agentType, properties)
	for i := 0; i < len(properties); i++ {
		if properties[i].name == "Time" {
			if !strings.ContainsAny(properties[i].value, "(") {
				d.timeVar = ParseProperty(properties[i].value)
			} else {
				d.timeMet = LoadMethod(properties[i].value)
			}
		}
	}
}

func (n *DecoratorTime) CreateTask() BehaviorTask {
	BTGLog.Trace("(d *DecoratorTime) CreateTask()")
	return &DecoratorTimeTask{}
}

func (n *DecoratorTime) IsValid(a *Agent, task BehaviorTask) bool {
	if _, ok := task.GetNode().(*DecoratorTime); !ok {
		return false
	}
	return true
}

func (n *DecoratorTime) GetTime(a *Agent) float64 {
	if n.timeVar != nil {
		return n.timeVar.GetFloat(a)
	} else if n.timeMet != nil {
		return n.timeMet.InvokeAndGetFloat(a)
	}
	return .0
}

// ============================================================================
type DecoratorTimeTask struct {
	DecoratorTask
	start int64
	time  int64
}

func (dt *DecoratorTimeTask) GetTime(a *Agent) float64 {
	if n, ok := dt.GetNode().(*DecoratorTime); ok {
		return n.GetTime(a)
	}
	return 0
}

func (dt *DecoratorTimeTask) OnEnter(a *Agent) bool {
	dt.start = BTGWorkspace.GetTimeSinceStartup()
	dt.time = int64(dt.GetTime(a)) * int64(time.Millisecond)
	return true
}

func (dt *DecoratorTimeTask) Decorate(status EBTStatus) EBTStatus {
	if BTGWorkspace.GetTimeSinceStartup()-dt.start >= dt.time {
		return BT_SUCCESS
	}
	return BT_RUNNING
}
