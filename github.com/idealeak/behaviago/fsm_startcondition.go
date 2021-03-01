package behaviago

import (
	"strconv"
)

// ============================================================================
type StartCondition struct {
	Precondition
	targetId  int
	effectors []*EffectorConfig
}

func (sc *StartCondition) GetTargetStateId() int {
	return sc.targetId
}

func (sc *StartCondition) SetTargetStateId(val int) {
	sc.targetId = val
}

func (sc *StartCondition) Load(version int, agentType string, properties []property_t) {
	sc.Precondition.Load(version, agentType, properties)
	if sc.loadAttachment {
		ec := &EffectorConfig{ActionConfig: &ActionConfig{}}
		ec.Load(properties)
		sc.effectors = append(sc.effectors, ec)
	}

	for i := 0; i < len(properties); i++ {
		if properties[i].name == "TargetFSMNodeId" {
			sc.targetId, _ = strconv.Atoi(properties[i].value)
		}
	}
}

func (sc *StartCondition) IsValid(a *Agent, task BehaviorTask) bool {
	if _, ok := task.GetNode().(*StartCondition); !ok {
		return false
	}
	return true
}

func (sc *StartCondition) ApplyEffects(a *Agent, phase EBTPhase) {
	for _, ec := range sc.effectors {
		if ec != nil {
			ec.Execute(a)
		}
	}
}

// ============================================================================
const (
	TransitionNodeName = "Transition"
)

func init() {
	RegisteNodeCreator(TransitionNodeName, BehaviorNodeCreatorWrapper(func() BehaviorNode {
		n := newTransitionNode()
		return n
	}))
}

func newTransitionNode() *Transition {
	n := &Transition{
		StartCondition: StartCondition{
			Precondition: Precondition{
				AttachAction: AttachAction{
					ac: &ActionConfig{},
				},
				pc: &PreconditionConfig{},
			},
		},
	}
	n.pc.ActionConfig = n.ac
	n.SetClassName(TransitionNodeName)
	n.SetSelf(n)
	return n
}

type Transition struct {
	StartCondition
}

func (t *Transition) IsValid(a *Agent, task BehaviorTask) bool {
	if _, ok := task.GetNode().(*Transition); !ok {
		return false
	}
	return true
}
