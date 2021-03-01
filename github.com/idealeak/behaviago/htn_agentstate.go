package behaviago

type AgentState struct {
	Variables
	parent *AgentState
	stack  []*AgentState
	forced bool
	pushed int
}

func NewAgentState(parent *AgentState) *AgentState {
	as := &AgentState{parent: parent}
	return as
}

func (as *AgentState) Depth() int {
	d := 1
	size := len(as.stack)
	if size != 0 {
		for i := size - 1; i >= 0; i-- {
			if as.stack[i] != nil {
				d += 1 + as.stack[i].pushed
			}
		}
	}
	return d
}

func (as *AgentState) Top() int {
	return len(as.stack) - 1
}

func (as *AgentState) Push(bForcePush bool) *AgentState {
	if !bForcePush {
		//if the top has nothing new added, to use it again

	}

	newly := NewAgentState(as)
	if newly != nil {
		newly.forced = bForcePush
		as.stack = append(as.stack, newly)
	}
	return newly
}

func (as *AgentState) Pop() {
	if as.parent == nil {
		return
	}
	if as.pushed > 0 {
		as.pushed--
	}
	size := len(as.stack)
	if size != 0 {
		top := as.stack[size-1]
		if top != nil {
			top.Pop()
			return
		}
	}

	as.parent.PopTop()
	as.parent = nil
}

func (as *AgentState) PopTop() {
	size := len(as.stack)
	if size != 0 {
		//remove the last one
		as.stack = as.stack[:size-1]
	}
}

func (as *AgentState) Clear(bFull bool) {

}
