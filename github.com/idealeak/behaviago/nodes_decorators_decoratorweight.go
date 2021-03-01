package behaviago

/**
  Specified the weight value of SelectorProbability child node.
*/
const (
	DecoratorWeightNodeName = "DecoratorWeight"
)

func init() {
	RegisteNodeCreator(DecoratorWeightNodeName, BehaviorNodeCreatorWrapper(func() BehaviorNode {
		n := newDecoratorWeightNode()
		return n
	}))
}

func newDecoratorWeightNode() *DecoratorWeight {
	n := &DecoratorWeight{}
	n.SetClassName(DecoratorWeightNodeName)
	n.SetSelf(n)
	return n
}

// ============================================================================
type DecoratorWeight struct {
	DecoratorNode
	varWeight *Property
}

func (dw *DecoratorWeight) IsManagingChildrenAsSubTrees() bool {
	return false
}

func (dw *DecoratorWeight) Load(version int, agentType string, properties []property_t) {
	dw.DecoratorNode.Load(version, agentType, properties)
	for i := 0; i < len(properties); i++ {
		if properties[i].name == "Weight" {
			dw.varWeight = ParseProperty(properties[i].value)
		}
	}
}

func (dw *DecoratorWeight) CreateTask() BehaviorTask {
	BTGLog.Trace("(dw *DecoratorWeight) CreateTask()")
	return &DecoratorWeightTask{}
}

func (dw *DecoratorWeight) IsValid(a *Agent, task BehaviorTask) bool {
	if _, ok := task.GetNode().(*DecoratorWeight); !ok {
		return false
	}
	return true
}

func (dw *DecoratorWeight) GetWeight(a *Agent) int64 {
	if dw.varWeight != nil {
		return dw.varWeight.GetInt(a)
	}
	return 0
}

// ============================================================================
type DecoratorWeightTask struct {
	DecoratorTask
}

func (dwt *DecoratorWeightTask) GetWeight(a *Agent) int64 {
	if dw, ok := dwt.GetNode().(*DecoratorWeight); ok {
		if dw != nil {
			return dw.GetWeight(a)
		}
	}
	return 0
}
