package behaviago

import (
	"strings"
)

/**
  It returns Running result until it reaches the frame count specified, no matter which
  value its child return. Or return the child's value.
*/
const (
	DecoratorFramesNodeName = "DecoratorFrames"
)

func init() {
	RegisteNodeCreator(DecoratorFramesNodeName, BehaviorNodeCreatorWrapper(func() BehaviorNode {
		n := newDecoratorFramesNode()
		return n
	}))
}

func newDecoratorFramesNode() *DecoratorFrames {
	n := &DecoratorFrames{}
	n.SetClassName(DecoratorFramesNodeName)
	n.SetSelf(n)
	return n
}

// ============================================================================
type DecoratorFrames struct {
	DecoratorNode
	framesVar *Property
	framesMet *Method
}

func (d *DecoratorFrames) Load(version int, agentType string, properties []property_t) {
	d.DecoratorNode.Load(version, agentType, properties)
	for i := 0; i < len(properties); i++ {
		if properties[i].name == "Frames" {
			if strings.Contains(properties[i].value, "(") {
				d.framesMet = LoadMethod(properties[i].value)
			} else {
				d.framesVar = ParseProperty(properties[i].value)
			}
		}
	}
}

func (d *DecoratorFrames) CreateTask() BehaviorTask {
	BTGLog.Trace("(d *DecoratorFrames) CreateTask()")
	return &DecoratorFramesTask{}
}

func (d *DecoratorFrames) GetFrames(a *Agent) int64 {
	if d.framesVar != nil {
		return d.framesVar.GetInt(a)
	} else if d.framesMet != nil {
		return d.framesMet.InvokeAndGetInt(a)
	}
	return 0
}

// ============================================================================
type DecoratorFramesTask struct {
	DecoratorTask
	start  int64
	frames int64
}

func (dt *DecoratorFramesTask) GetFrames(a *Agent) int64 {
	if node, ok := dt.GetNode().(*DecoratorFrames); ok {
		return node.GetFrames(a)
	}
	return 0
}

func (dt *DecoratorFramesTask) OnEnter(a *Agent) bool {
	dt.start = BTGWorkspace.GetFrameSinceStartup()
	dt.frames = dt.GetFrames(a)
	return dt.DecoratorTask.OnEnter(a)
}

func (dt *DecoratorFramesTask) Decorate(status EBTStatus) EBTStatus {
	if BTGWorkspace.GetFrameSinceStartup()-dt.start+1 > dt.frames {
		return BT_SUCCESS
	}
	return BT_FAILURE
}
