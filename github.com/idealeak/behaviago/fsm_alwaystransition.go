package behaviago

// ============================================================================
type ETransitionPhase int

const (
	ETP_Always ETransitionPhase = iota
	ETP_Success
	ETP_Failure
	ETP_Exit
)

// ============================================================================
type AlwaysTransition struct {
	Transition
	transitionPhase ETransitionPhase
}

func (at *AlwaysTransition) Load(version int, agentType string, properties []property_t) {
	at.Transition.Load(version, agentType, properties)
	for i := 0; i < len(properties); i++ {
		if properties[i].name == "TransitionPhase" {
			switch properties[i].value {
			case "ETP_Exit":
				at.transitionPhase = ETP_Exit
			case "ETP_Success":
				at.transitionPhase = ETP_Success
			case "ETP_Failure":
				at.transitionPhase = ETP_Failure
			case "ETP_Always":
				at.transitionPhase = ETP_Always
			}
		}
	}
}

func (at *AlwaysTransition) IsValid(a *Agent, task BehaviorTask) bool {
	if _, ok := task.GetNode().(*AlwaysTransition); !ok {
		return false
	}
	return true
}

func (at *AlwaysTransition) Evaluate(a *Agent, status EBTStatus) bool {
	if at.transitionPhase == ETP_Always {
		return true
	} else if status == BT_SUCCESS && (at.transitionPhase == ETP_Success || at.transitionPhase == ETP_Exit) {
		return true
	} else if status == BT_FAILURE && (at.transitionPhase == ETP_Failure || at.transitionPhase == ETP_Exit) {
		return true
	}
	return false
}
