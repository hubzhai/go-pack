package behaviago

// ============================================================================
type EPreCondPhase int

const (
	E_PRECOND_ENTER EPreCondPhase = iota
	E_PRECOND_UPDATE
	E_PRECOND_BOTH
)

// ============================================================================
type PreconditionInterface interface {
	GetPhase() EPreCondPhase
	SetPhase(EPreCondPhase)
	IsAnd() bool
	SetIsAnd(bool)
}

// ============================================================================
type PreconditionConfig struct {
	*ActionConfig
	phase EPreCondPhase
	bAnd  bool
}

func (pc *PreconditionConfig) Load(properties []property_t) {
	pc.ActionConfig.Load(properties)
	for i := 0; i < len(properties); i++ {
		switch properties[i].name {
		case "BinaryOperator":
			switch properties[i].value {
			case "Or":
				pc.bAnd = false
			case "And":
				pc.bAnd = true
			}
		case "Phase":
			switch properties[i].value {
			case "Enter":
				pc.phase = E_PRECOND_ENTER
			case "Update":
				pc.phase = E_PRECOND_UPDATE
			case "Both":
				pc.phase = E_PRECOND_BOTH
			}
		}
	}
}

// ============================================================================
const (
	PreconditionNodeName = "Precondition"
)

func init() {
	RegisteNodeCreator(PreconditionNodeName, BehaviorNodeCreatorWrapper(func() BehaviorNode {
		n := newPreconditionNode()
		return n
	}))
}

func newPreconditionNode() *Precondition {
	n := &Precondition{
		AttachAction: AttachAction{
			ac: &ActionConfig{},
		},
		pc: &PreconditionConfig{},
	}
	n.pc.ActionConfig = n.ac
	n.SetClassName(PreconditionNodeName)
	n.SetSelf(n)
	return n
}

type Precondition struct {
	AttachAction
	pc *PreconditionConfig
}

func (p *Precondition) Load(version int, agentType string, properties []property_t) {
	p.AttachAction.Load(version, agentType, properties)
	p.pc.Load(properties)
}

func (p *Precondition) GetPhase() EPreCondPhase {
	return p.pc.phase
}
func (p *Precondition) SetPhase(value EPreCondPhase) {
	p.pc.phase = value
}
func (p *Precondition) IsAnd() bool {
	return p.pc.bAnd
}
func (p *Precondition) SetIsAnd(isAnd bool) {
	p.pc.bAnd = isAnd
}
func (p *Precondition) IsValid(a *Agent, task BehaviorTask) bool {
	if _, ok := task.GetNode().(*Precondition); !ok {
		return false
	}
	return true
}
