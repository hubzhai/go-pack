package behaviago

const (
	WaitTransitionNodeName = "WaitTransition"
)

func init() {
	RegisteNodeCreator(WaitTransitionNodeName, BehaviorNodeCreatorWrapper(func() BehaviorNode {
		n := newWaitTransitionNode()
		return n
	}))
}

func newWaitTransitionNode() *WaitTransition {
	n := &WaitTransition{
		Transition: Transition{
			StartCondition: StartCondition{
				Precondition: Precondition{
					AttachAction: AttachAction{
						ac: &ActionConfig{},
					},
					pc: &PreconditionConfig{},
				},
			},
		},
	}
	n.pc.ActionConfig = n.ac
	n.SetClassName(WaitTransitionNodeName)
	n.SetSelf(n)
	return n
}

type WaitTransition struct {
	Transition
}

func (wt *WaitTransition) IsValid(a *Agent, task BehaviorTask) bool {
	if _, ok := task.GetNode().(*WaitTransition); !ok {
		return false
	}
	return true
}

func (wt *WaitTransition) Evaluate(a *Agent, status EBTStatus) bool {
	return true
}
