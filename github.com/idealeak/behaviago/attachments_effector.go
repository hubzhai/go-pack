package behaviago

// ============================================================================
type EffectorInterface interface {
	GetPhace() EBTPhase
	SetPhace(EBTPhase)
	Evaluate(a *Agent, result EBTStatus) bool
}

// ============================================================================
type EffectorConfig struct {
	*ActionConfig
	phase EBTPhase
}

func (ec *EffectorConfig) Load(properties []property_t) {
	ec.ActionConfig.Load(properties)
	for i := 0; i < len(properties); i++ {
		if properties[i].name == "Phase" {
			switch properties[i].value {
			case "Success":
				ec.phase = EBT_SUCCESS
			case "Failure":
				ec.phase = EBT_FAILURE
			case "Both":
				ec.phase = EBT_BOTH
			}
		}
	}
}

// ============================================================================
const (
	EffectorNodeName = "Effector"
)

func init() {
	RegisteNodeCreator(EffectorNodeName, BehaviorNodeCreatorWrapper(func() BehaviorNode {
		n := newEffectorNode()
		return n
	}))
}

func newEffectorNode() *Effector {
	n := &Effector{
		AttachAction: AttachAction{
			ac: &ActionConfig{},
		},
		ec: &EffectorConfig{},
	}
	n.ec.ActionConfig = n.ac
	n.SetClassName(EffectorNodeName)
	n.SetSelf(n)
	return n
}

type Effector struct {
	AttachAction
	ec *EffectorConfig
}

func (e *Effector) GetPhace() EBTPhase {
	return e.ec.phase
}

func (e *Effector) SetPhace(value EBTPhase) {
	e.ec.phase = value
}

func (e *Effector) IsValid(a *Agent, task BehaviorTask) bool {
	if _, ok := task.GetNode().(*Effector); !ok {
		return false
	}
	return true
}
