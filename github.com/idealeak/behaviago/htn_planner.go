package behaviago

type Planner struct {
	a            *Agent
	AutoReplan   bool
	rootTaskNode *Task
	rootTask     PlannerTask
}

func (p *Planner) GetAgent() *Agent {
	return p.a
}

func (p *Planner) Init(a *Agent, root *Task) {
	p.a = a
	p.rootTaskNode = root
}

func (p *Planner) Uninit() {
	p.onDisable()
}

func (p *Planner) onDisable() {
	if p.rootTask != nil {
		if p.rootTask.GetStatue() == BT_RUNNING {
			p.rootTask.Abort(p.a)
		}
		p.rootTask = nil
	}
}

func (p *Planner) Update() EBTStatus {
	if p.a == nil {
		return BT_INVALID
	}
	p.doAutoPlanning()
	if p.rootTask == nil {
		return BT_FAILURE
	}
	// Need a local reference in case the p.rootTask is cleared by an event handler
	rootTask := p.rootTask
	return rootTask.Exec(p.a, BT_RUNNING)
}

/// <summary>
/// Generate a new task for the <paramref name="agent"/> based on the current world state as
/// described by the <paramref name="agentState"/>.
/// </summary>
/// <param name="agent">The agent for which the task is being generated. This object instance must be
/// of the same type as the type for which the Task was developed</param>
/// <param name="agentState">The current world state required by the planner</param>
/// <returns></returns>
func (p *Planner) generatePlan() PlannerTask {
	// If the planner is currently executing a task marked NotInterruptable, do not generate
	// any n ew plans.
	if !p.canInterruptCurrentPlan() {
		return nil
	}
	newPlan := p.BuildPlan(p.rootTaskNode)
	if newPlan == nil {
		return nil
	}
	if !newPlan.IsHigherPriority(p.rootTask) {
		return nil
	}

	return newPlan
}

func (p *Planner) canInterruptCurrentPlan() bool {
	if p.rootTask == nil {
		return true
	}
	status := p.rootTask.GetStatue()
	if status != BT_RUNNING {
		return true
	}
	if !p.rootTask.CannNotInterruptable() {
		return true
	}
	return status == BT_FAILURE || status == BT_SUCCESS
}

func (p *Planner) doAutoPlanning() {
	if !p.AutoReplan {
		return
	}
	noPlan := p.rootTask == nil || p.rootTask.GetStatue() != BT_RUNNING
	if noPlan {
		newPlan := p.generatePlan()
		if newPlan != nil {
			if p.rootTask != nil {
				if p.rootTask.GetStatue() == BT_RUNNING {
					p.rootTask.Abort(p.a)
				}
			}
			p.rootTask = newPlan
		}
	}
}

func (p *Planner) BuildPlan(root *Task) PlannerTask {
	return nil
}

func (p *Planner) DecomposeNode(node BehaviorNode, depth int) PlannerTask {
	// Ensure that the planner does not get stuck in an infinite loop
	if depth > 256 {
		BTGLog.Warnf("Exceeded task nesting depth. Does the graph contain an invalid cycle?")
		return nil
	}

	p.LogPlanNodeBegin(p.a, node)

	var taskAdd PlannerTask
	isPreconditionOk := node.CheckPreconditions(p.a, false)
	if isPreconditionOk {
		bOk := true
		taskAdd = CreatePlannerTask(node, p.a)
		if seqTask, ok := taskAdd.(*PlannerTaskComplex); ok {
			bOk = p.DecomposeComplex(node, seqTask, depth)
		}
		if bOk {
			node.ApplyEffects(p.a, EBT_SUCCESS)
		} else {
			taskAdd = nil
		}
	} else {
		p.LogPlanNodePreconditionFailed(p.a, node)
	}

	if taskAdd != nil {
		p.LogPlanNodeEnd(p.a, node, "success")
	} else {
		p.LogPlanNodeEnd(p.a, node, "failure")
	}

	return taskAdd
}

func (p *Planner) DecomposeComplex(node BehaviorNode, seqTask *PlannerTaskComplex, depth int) bool {
	return node.Decompose(node, seqTask, depth, p)
}

func (p *Planner) DecomposeTask(task *Task, depth int) PlannerTask {
	childs := task.GetChilds()
	methodsCount := len(childs)
	if methodsCount == 0 {
		return nil
	}
	for i := 0; i < methodsCount; i++ {

	}
	return nil
}

func (p *Planner) LogPlanBegin(a *Agent, root *Task) {

}

func (p *Planner) LogPlanEnd(a *Agent, root *Task) {

}

func (p *Planner) LogPlanNodeBegin(a *Agent, node BehaviorNode) {

}

func (p *Planner) LogPlanNodePreconditionFailed(a *Agent, node BehaviorNode) {

}

func (p *Planner) LogPlanNodeEnd(a *Agent, node BehaviorNode, result string) {

}

func (p *Planner) LogPlanReferenceTreeEnter(a *Agent, tree ReferencedBehavior) {

}

func (p *Planner) LogPlanReferenceTreeExit(a *Agent, tree ReferencedBehavior) {

}

func (p *Planner) LogPlanMethodBegin(a *Agent, m BehaviorNode) {

}

func (p *Planner) LogPlanMethodEnd(a *Agent, m BehaviorNode, result string) {

}

func (p *Planner) LogPlanForEachBegin(a *Agent, foreach *DecoratorIterator, index, count int) {

}

func (p *Planner) LogPlanForEachEnd(a *Agent, foreach *DecoratorIterator, index, count int, result string) {

}
